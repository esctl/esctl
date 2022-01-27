package es

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/esctl/esctl/pkg/printer"
)

func (e *elasticSearchClient) GetHealth() (HealthResponse, error) {
	healthResponse := HealthResponse{}
	req := esapi.ClusterHealthRequest{}
	res, err := req.Do(context.Background(), e.c)
	if err != nil {
		return healthResponse, fmt.Errorf("calling health request to elastic search failed, %w", err)
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

func (h HealthResponse) Print(p printer.Printer) error {
	s, err := p.BigLettersWithColor(h.Status, h.Status)
	if err != nil {
		return fmt.Errorf("printing health output failed, %w", err)
	}
	fmt.Println(s)
	fmt.Printf("Cluster Name :%s\n", p.HighlightText(h.ClusterName))
	return nil
}
