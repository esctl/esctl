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

func TestESSearchClient_CatNodes(t *testing.T) {
	resp := `ip         heap.percent ram.percent cpu load_1m load_5m load_15m node.role   master name
	172.20.0.2           59          67  76    1.73    0.41     0.14 cdfhilmrstw -      es03
	172.20.0.3           40          67  76    1.73    0.41     0.14 cdfhilmrstw -      es02
	172.20.0.4           36          67  76    1.73    0.41     0.14 cdfhilmrstw *      es01`
	client := elasticSearchClient{}
	mockC, _ := elasticsearch.NewClient(elasticsearch.Config{Transport: &mockTransport{
		roundTripFunc: func(*http.Request) (*http.Response, error) {
			return &http.Response{
					Body: ioutil.NopCloser(strings.NewReader(resp)),
					Header: http.Header{
						"content-type":      []string{"text/plain"},
						"X-Elastic-Product": []string{"Elasticsearch"},
					},
				},
				nil
		},
	}})

	client.c = mockC

	actualResp, err := client.CatNodes()
	assert.Nil(t, err, "error not expected")
	assert.Equal(t, resp, actualResp, "Incorrect cat nodes response")
}

func TestESSearchClient_CatNodes_Error(t *testing.T) {
	expectedError := errors.New("An error occurred")
	client := elasticSearchClient{}
	mockC, _ := elasticsearch.NewClient(elasticsearch.Config{Transport: &mockTransport{
		roundTripFunc: func(*http.Request) (*http.Response, error) {
			return &http.Response{
					Body: ioutil.NopCloser(strings.NewReader(``)),
					Header: http.Header{
						"content-type":      []string{"text/plain"},
						"X-Elastic-Product": []string{"Elasticsearch"},
					},
				},
				expectedError
		},
	}})

	client.c = mockC

	_, err := client.CatNodes()
	assert.Error(t, err, "Missing cat nodes error")
	assert.Contains(t, err.Error(), "calling cat nodes to elastic search failed, ", "Incorrect cat nodes error message")
}
