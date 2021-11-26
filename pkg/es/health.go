package es

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/pterm/pterm"
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

func (h HealthResponse) Print() {
	var st pterm.Color
	switch h.Status {
	case ClusterStatusGreen:
		st = pterm.FgGreen
	case ClusterStatusYellow:
		st = pterm.FgYellow
	case ClusterStatusRed:
		st = pterm.FgRed
	}

	err := pterm.DefaultBigText.WithLetters(pterm.NewLettersFromStringWithStyle(h.Status, pterm.NewStyle(st))).Render()
	if err != nil {
		log.Fatalf("Error printing health output %v", err)
	}
	pterm.DefaultBasicText.Println("Cluster Name", pterm.FgBlue.Sprint(h.ClusterName))
}
