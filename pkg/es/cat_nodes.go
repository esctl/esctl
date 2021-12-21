package es

import (
	"context"
	"fmt"
	"io/ioutil"

	"github.com/elastic/go-elasticsearch/v7/esapi"
)

func (e *elasticSearchClient) CatNodes() (string, error) {
	verbose := true
	req := esapi.CatNodesRequest{
		Pretty: true,
		Human:  true,
		V:      &verbose,
	}
	res, err := req.Do(context.Background(), e.c)
	if err != nil {
		return "", fmt.Errorf("calling cat nodes to elastic search failed, %w", err)
	}

	resData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("reading response data failed, %w", err)
	}

	return string(resData), nil
}
