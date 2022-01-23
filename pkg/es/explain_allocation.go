package es

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/elastic/go-elasticsearch/v7/esapi"
)

func (e *elasticSearchClient) ExplainAllocation() (string, error) {

	req := esapi.ClusterAllocationExplainRequest{}
	res, err := req.Do(context.Background(), e.c)
	if err != nil {
		return "", fmt.Errorf("calling allocation explain request to elastic search failed, %w", err)
	}

	resData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("reading response data for allocation explain failed, %w", err)
	}

	var prettyJSON bytes.Buffer
	err = json.Indent(&prettyJSON, resData, "", "\t")
	if err != nil {
		return "", fmt.Errorf("indenting allocation explain failed, %w", err)
	}

	return prettyJSON.String(), nil
}
