package es

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/esctl/esctl/pkg/config"
)

func New(c *config.Cluster) (*ElasticSearchClient, error) {

	if len(c.Hosts) == 0 {
		return nil, errors.New("no hosts provided in config")
	}
	cfg := elasticsearch.Config{
		Addresses: c.Hosts,
	}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, err
	}
	return &ElasticSearchClient{
		c: es,
	}, nil
}

type ElasticSearchClient struct {
	c *elasticsearch.Client
}

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

type mockTransport struct {
	roundTripFunc func(*http.Request) (*http.Response, error)
}

func (t *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return t.roundTripFunc(req)
}
