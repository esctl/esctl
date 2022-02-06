package es

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/elastic/go-elasticsearch/v7/esapi"
)

type ListIndexResponse []string

func (l ListIndexResponse) Print() {
	fmt.Printf("%v\n", strings.Join(l, "\n"))
}

func (e *elasticSearchClient) ListIndex() (ListIndexResponse, error) {
	req := esapi.IndicesGetRequest{
		Index: []string{"_all"},
	}
	res, err := req.Do(context.Background(), e.c)
	if err != nil {
		return nil, fmt.Errorf("calling get indices failed, %w", err)
	}

	resData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response failed, %w", err)
	}

	// Extracts the index names from the raw response body
	c := make(map[string]json.RawMessage)
	err = json.Unmarshal(resData, &c)
	if err != nil {
		return nil, fmt.Errorf("extract index name(s): json unmarshall response failed, %w", err)
	}
	indices := []string{}
	for key := range c {
		indices = append(indices, key)
	}

	return indices, nil
}

func (e *elasticSearchClient) AllIndexSettings() (string, error) {
	return e.IndexSettings([]string{"_all"})
}

func (e *elasticSearchClient) IndexSettings(indices []string) (string, error) {
	req := esapi.IndicesGetRequest{
		Index: indices,
	}
	res, err := req.Do(context.Background(), e.c)
	if err != nil {
		return "", fmt.Errorf("calling get indices failed, %w", err)
	}

	resData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("reading response failed, %w", err)
	}

	var prettyJSON bytes.Buffer
	err = json.Indent(&prettyJSON, resData, "", "  ")
	if err != nil {
		return "", fmt.Errorf("indenting index settings failed, %w", err)
	}

	return prettyJSON.String(), nil
}
