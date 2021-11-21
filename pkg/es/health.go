package es

import (
	"context"
	"encoding/json"

	"github.com/elastic/go-elasticsearch/v7/esapi"
)

func (e *ElasticSearchClient) GetHealth() (HealthResponse, error) {
	healthResponse := HealthResponse{}
	req := esapi.ClusterHealthRequest{}
	res, err := req.Do(context.Background(), e.c)
	if err != nil {
		return healthResponse, err
	}

	err = json.NewDecoder(res.Body).Decode(&healthResponse)
	if err != nil {
		return healthResponse, err
	}

	return healthResponse, nil
}

type HealthResponse struct {
	ClusterName string `json:"cluster_name"`
	Status      string
}
