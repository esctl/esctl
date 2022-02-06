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

var (
	indexResp = `{
		"application-logs-001": {
			"settings": {
				"number_of_shards": "1",
				"provided_name": "application-logs-001",
				"creation_date": "1643516445292",
				"number_of_replicas": "1"
			}
		},
		"application-logs-002": {
			"settings": {
				"number_of_shards": "1",
				"provided_name": "application-logs-002",
				"creation_date": "1643516445292",
				"number_of_replicas": "1"
			}
		}
	}`
)

func Test_elasticSearchClient_ListIndex(t *testing.T) {
	expectedResp := ListIndexResponse{"application-logs-001", "application-logs-002"}
	mockC, _ := elasticsearch.NewClient(elasticsearch.Config{Transport: &mockTransport{
		roundTripFunc: func(*http.Request) (*http.Response, error) {
			return &http.Response{
					Body: ioutil.NopCloser(strings.NewReader(indexResp)),
					Header: http.Header{
						"content-type":      []string{"text/plain"},
						"X-Elastic-Product": []string{"Elasticsearch"},
					},
				},
				nil
		},
	}})
	client := elasticSearchClient{c: mockC}

	actualResp, err := client.ListIndex()
	assert.Nil(t, err, "error not expected")
	assert.ElementsMatch(t, expectedResp, actualResp, "Incorrect list index response")
}

func Test_elasticSearchClient_ListIndex_Error(t *testing.T) {
	mockC, _ := elasticsearch.NewClient(elasticsearch.Config{Transport: &mockTransport{
		roundTripFunc: func(*http.Request) (*http.Response, error) {
			return &http.Response{
					Body: ioutil.NopCloser(strings.NewReader(``)),
					Header: http.Header{
						"content-type":      []string{"text/plain"},
						"X-Elastic-Product": []string{"Elasticsearch"},
					},
				},
				errors.New("An error occurred")
		},
	}})
	client := elasticSearchClient{c: mockC}

	_, err := client.ListIndex()
	assert.Error(t, err, "Missing list index error")
	assert.Contains(t, err.Error(), "calling get indices failed, ", "Incorrect get index error message")
}

func Test_elasticSearchClient_ListIndex_JSON_Unmarshall_Error(t *testing.T) {
	resp := `invalid json`
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
	client := elasticSearchClient{c: mockC}

	_, err := client.ListIndex()
	assert.Error(t, err, "Missing list index error")
	assert.Contains(t, err.Error(), "extract index name(s): json unmarshall response failed, ", "Incorrect get index error message")
}

func Test_elasticSearchClient_AllIndexSettings(t *testing.T) {
	mockC, _ := elasticsearch.NewClient(elasticsearch.Config{Transport: &mockTransport{
		roundTripFunc: func(*http.Request) (*http.Response, error) {
			return &http.Response{
					Body: ioutil.NopCloser(strings.NewReader(indexResp)),
					Header: http.Header{
						"content-type":      []string{"text/plain"},
						"X-Elastic-Product": []string{"Elasticsearch"},
					},
				},
				nil
		},
	}})
	client := elasticSearchClient{c: mockC}

	actualResp, err := client.AllIndexSettings()
	assert.Nil(t, err, "error not expected")
	assert.JSONEq(t, indexResp, actualResp, "Incorrect list index response")
}

func Test_elasticSearchClient_IndexSettings(t *testing.T) {
	resp := `{
		"application-logs-001": {
			"settings": {
				"number_of_shards": "1",
				"provided_name": "application-logs-001",
				"creation_date": "1643516445292",
				"number_of_replicas": "1"
			}
		}
	}`
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
	client := elasticSearchClient{c: mockC}

	actualResp, err := client.IndexSettings([]string{"application-logs-001"})
	assert.Nil(t, err, "error not expected")
	assert.JSONEq(t, resp, actualResp, "Incorrect list index response")
}

func Test_elasticSearchClient_IndexSettings_Error(t *testing.T) {
	mockC, _ := elasticsearch.NewClient(elasticsearch.Config{Transport: &mockTransport{
		roundTripFunc: func(*http.Request) (*http.Response, error) {
			return &http.Response{
					Body: ioutil.NopCloser(strings.NewReader(indexResp)),
					Header: http.Header{
						"content-type":      []string{"text/plain"},
						"X-Elastic-Product": []string{"Elasticsearch"},
					},
				},
				errors.New("An error occurred")
		},
	}})
	client := elasticSearchClient{c: mockC}

	_, err := client.IndexSettings([]string{"application-logs-001"})
	assert.Error(t, err, "Missing list index error")
	assert.Contains(t, err.Error(), "calling get indices failed, ", "Incorrect get index error message")

}

func Test_elasticSearchClient_IndexSettings_JSON_Indent_Error(t *testing.T) {
	resp := `invalid json`
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
	client := elasticSearchClient{c: mockC}

	_, err := client.IndexSettings([]string{"application-logs-001"})
	assert.Error(t, err, "Missing list index error")
	assert.Contains(t, err.Error(), "indenting index settings failed, ", "Incorrect get index error message")
}
