package es

import (
	"errors"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/esctl/esctl/pkg/config"
)

func New(c *config.Cluster) (*elasticSearchClient, error) {

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
	return &elasticSearchClient{
		c: es,
	}, nil
}

type elasticSearchClient struct {
	c *elasticsearch.Client
}
