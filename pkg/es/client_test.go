package es

import (
	"errors"
	"testing"

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
