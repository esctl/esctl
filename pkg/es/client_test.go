package es

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/esctl/esctl/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestEsClient_New_ReturnsClient(t *testing.T) {
	cfg := config.Cluster{
		Name: "Hello",
		Hosts: []string{
			"http://hello1:9700",
		},
	}
	c, err := New(&cfg)
	assert.Nil(t, err, "Error should be nil")
	assert.NotNil(t, c, "Client should not be nil")
}

func TestEsClient_New_Returns_err_when_host_missing(t *testing.T) {
	cfg := config.Cluster{
		Name:  "Hello",
		Hosts: []string{},
	}
	c, err := New(&cfg)
	assert.NotNil(t, err, "Error should be nil")
	assert.Equal(t, errors.New("no hosts provided in config"), err)
	assert.Nil(t, c, "Client should not be nil")
}

func TestESSearchClient_GetHealth(t *testing.T) {
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
	client := ElasticSearchClient{}
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

func TestESSearchClient_GetHealth_Error(t *testing.T) {

	client := ElasticSearchClient{}
	mockC, _ := elasticsearch.NewClient(elasticsearch.Config{Transport: &mockTransport{
		roundTripFunc: func(*http.Request) (*http.Response, error) {
			return &http.Response{
					Body:   ioutil.NopCloser(strings.NewReader("{}")),
					Header: http.Header{"X-Elastic-Product": []string{"Elasticsearch"}},
				},
				errors.New("An error occurred")
		},
	}})

	client.c = mockC

	actualResp, err := client.GetHealth()
	assert.NotNil(t, err, "error expected")
	assert.Equal(t, errors.New("An error occurred"), err)
	assert.Equal(t, HealthResponse{}, actualResp)
}