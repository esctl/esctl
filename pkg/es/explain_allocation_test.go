package es

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/stretchr/testify/assert"
)

func Test_elasticSearchClient_ExplainAllocation(t *testing.T) {
	resp := `{
	"error": {
		"root_cause": [
			{
				"type": "illegal_argument_exception",
				"reason": "No shard was specified in the request which means the response should explain a randomly-chosen unassigned shard, but there are no unassigned shards in this cluster. To explain the allocation of an assigned shard you must specify the target shard in the request."
			}
		],
		"type": "illegal_argument_exception",
		"reason": "No shard was specified in the request which means the response should explain a randomly-chosen unassigned shard, but there are no unassigned shards in this cluster. To explain the allocation of an assigned shard you must specify the target shard in the request."
	},
	"status": 400
}`
	client := elasticSearchClient{}
	mockC, _ := elasticsearch.NewClient(elasticsearch.Config{
		Transport: &mockTransport{
			roundTripFunc: func(*http.Request) (*http.Response, error) {
				return &http.Response{
						Body:   ioutil.NopCloser(strings.NewReader(resp)),
						Header: http.Header{"X-Elastic-Product": []string{"Elasticsearch"}},
					},
					nil
			},
		}})

	client.c = mockC
	actualResp, err := client.ExplainAllocation()
	assert.Nil(t, err, "error not expected")
	assert.Equal(t, resp, actualResp, "Incorrect explain allocation response")
}

func Test_elasticSearchClient_ExplainAllocation_Error(t *testing.T) {
	e := errors.New("Something went wrong")
	client := elasticSearchClient{}
	mockC, _ := elasticsearch.NewClient(elasticsearch.Config{
		Transport: &mockTransport{
			roundTripFunc: func(*http.Request) (*http.Response, error) {
				return &http.Response{
						Body:   ioutil.NopCloser(strings.NewReader("")),
						Header: http.Header{"X-Elastic-Product": []string{"Elasticsearch"}},
					},
					e
			},
		}})
	client.c = mockC
	_, err := client.ExplainAllocation()
	assert.NotNil(t, err, "error was expected, are we not returning it ?")

}
