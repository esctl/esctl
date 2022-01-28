package es

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/esctl/esctl/pkg/printer"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/stretchr/testify/assert"
)

func TestESClient_GetHealth(t *testing.T) {
	resp := `
	{
		"active_primary_shards": 1,
		"active_shards": 2,
		"active_shards_percent_as_number": 100.0,
		"cluster_name": "es-docker-cluster",
		"delayed_unassigned_shards": 0,
		"initializing_shards": 0,
		"number_of_data_nodes": 3,
		"number_of_in_flight_fetch": 0,
		"number_of_nodes": 3,
		"number_of_pending_tasks": 0,
		"relocating_shards": 0,
		"status": "green",
		"task_max_waiting_in_queue_millis": 0,
		"timed_out": false,
		"unassigned_shards": 0
	}
	`
	client := elasticSearchClient{}
	mockC, _ := elasticsearch.NewClient(elasticsearch.Config{Transport: &mockTransport{
		roundTripFunc: func(*http.Request) (*http.Response, error) {
			return &http.Response{
					Body:   ioutil.NopCloser(strings.NewReader(resp)),
					Header: http.Header{"X-Elastic-Product": []string{"Elasticsearch"}},
				},
				nil
		},
	}})

	client.c = mockC

	actualResp, err := client.GetHealth()
	assert.Nil(t, err, "error not expected")
	assert.Equal(t, "green", actualResp.Status)
	assert.Equal(t, "es-docker-cluster", actualResp.ClusterName)
}

func TestESClient_GetHealth_Error(t *testing.T) {

	expectedError := errors.New("An error occurred")
	client := elasticSearchClient{}
	mockC, _ := elasticsearch.NewClient(elasticsearch.Config{Transport: &mockTransport{
		roundTripFunc: func(*http.Request) (*http.Response, error) {
			return &http.Response{
					Body:   ioutil.NopCloser(strings.NewReader("{}")),
					Header: http.Header{"X-Elastic-Product": []string{"Elasticsearch"}},
				},
				expectedError
		},
	}})

	client.c = mockC

	actualResp, err := client.GetHealth()
	assert.NotNil(t, err, "error expected")
	assert.Equal(t, fmt.Errorf("calling health request to elastic search failed, %w", expectedError), err)
	assert.Equal(t, HealthResponse{}, actualResp)
}

func TestHealResponse_Print(t *testing.T) {
	p := printer.MockPrinter{}
	expectedText := "Green in big text"
	hr := HealthResponse{
		ClusterName: "TestCluster",
		Status:      "green",
	}

	p.On("BigTextWithColor", hr.Status, hr.Status).Return(expectedText, nil)
	p.On("HighlightText", hr.ClusterName).Return(hr.ClusterName)

	err := hr.Print(&p)
	assert.Nil(t, err, "Not expecting an error here")
	p.AssertExpectations(t)

}

func TestHealResponse_Print_Error(t *testing.T) {
	p := printer.MockPrinter{}
	expectedText := "Gree in big text"
	hr := HealthResponse{
		ClusterName: "TestCluster",
		Status:      "green",
	}
	renderErr := errors.New("render error")
	expectedError := fmt.Errorf("printing health output failed, %w", renderErr)
	p.On("BigTextWithColor", hr.Status, hr.Status).Return(expectedText, renderErr)

	err := hr.Print(&p)
	assert.NotNil(t, err)
	assert.Equal(t, expectedError, err, "error values do not match")
	p.AssertExpectations(t)

}
