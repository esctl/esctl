package config

import (
	stdlib_fs "io/fs"
	"testing"

	"github.com/esctl/esctl/pkg/fs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	yaml "gopkg.in/yaml.v3"
)

func TestClusterConfig_Load(t *testing.T) {
	type fields struct {
		CurrentCluster string
		Clusters       []Cluster
		cfgFile        string
		ReadFn         fs.ReadFn
	}

	type args struct {
		cfgFile string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Should load config from valid yaml file",
			args: args{
				cfgFile: "/some/config.yaml",
			},
			fields: fields{
				ReadFn: func(s string) ([]byte, error) {
					c := `
                          CurrentCluster: ""
                          Clusters:
                            - Name: Test
                              Hosts:
                                - "http://node1:1234"
                                - "http://node2:1234"`
					b := []byte(c)
					return b, nil
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ClusterConfig{
				CurrentCluster: tt.fields.CurrentCluster,
				Clusters:       tt.fields.Clusters,
				cfgFile:        tt.fields.cfgFile,
				r:              tt.fields.ReadFn,
			}
			err := c.Load(tt.args.cfgFile)
			assert.Nil(t, err, "error not expected")
			assert.Equal(t, "Test", c.Clusters[0].Name, "Cluster Names do not match")
			assert.Equal(t, "http://node1:1234", c.Clusters[0].Hosts[0], "Host do not match")
			assert.Equal(t, "http://node2:1234", c.Clusters[0].Hosts[1], "Host do not match")
			assert.Equal(t, "/some/config.yaml", c.cfgFile, "cfgFile do not match")
		})
	}
}

func TestClusterConfig_Write(t *testing.T) {

	localClusterConfig := ClusterConfig{
		cfgFile:        "/some/conf.yaml",
		CurrentCluster: "",
		Clusters: []Cluster{
			{
				Name:  "local",
				Hosts: []string{"http://node1:1234", "http://node2:1234"},
			},
		},
	}
	localClusterConfig.w = func(s string, b []byte, fm stdlib_fs.FileMode) error {
		expected, err := yaml.Marshal(localClusterConfig)
		assert.Nil(t, err, "Unexpected yaml marshal error")
		assert.Equal(t, expected, b, "config data mismatch")
		return nil
	}

	tests := []struct {
		name   string
		fields ClusterConfig
	}{
		{
			name:   "should write config to file",
			fields: localClusterConfig,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ClusterConfig{
				CurrentCluster: tt.fields.CurrentCluster,
				Clusters:       tt.fields.Clusters,
				cfgFile:        tt.fields.cfgFile,
				w:              tt.fields.w,
			}
			err := c.write()
			assert.Nil(t, err, "error not expected")
		})
	}
}

func TestClusterConfig_SetActive(t *testing.T) {
	c := &ClusterConfig{
		Clusters: []Cluster{
			{Name: "Server1", Hosts: []string{"http:node1:1234"}},
			{Name: "Server2", Hosts: []string{"http://node2:5678"}},
		},
	}

	c.w = func(s string, b []byte, fm stdlib_fs.FileMode) error {
		expected, err := yaml.Marshal(c)
		assert.Nil(t, err, "Unexpected yaml marshal error")
		assert.Equal(t, expected, b, "Config data mismatch")
		return nil
	}

	err := c.SetActive("Server1")

	require.NoError(t, err, "Unexpected set active cluster error")
	assert.Equal(t, "Server1", c.CurrentCluster, "Incorrect current cluster name")
}

func TestClusterConfig_SetActiveForInvalidClusterName(t *testing.T) {
	c := &ClusterConfig{
		Clusters: []Cluster{
			{Name: "Server1", Hosts: []string{"http:node1:1234"}},
			{Name: "Server2", Hosts: []string{"http://node2:5678"}},
		},
	}

	c.w = func(s string, b []byte, fm stdlib_fs.FileMode) error {
		expected, err := yaml.Marshal(c)
		assert.Nil(t, err, "Unexpected yaml marshal error")
		assert.Equal(t, expected, b, "Config data mismatch")
		return nil
	}

	err := c.SetActive("Server3")

	assert.EqualError(t, err, "cluster Server3 not found", "Incorrect set active cluster error")
	assert.Equal(t, "", c.CurrentCluster, "Incorrect current cluster name")
}

func TestClusterConfig_DeleteCluster(t *testing.T) {
	c := &ClusterConfig{
		Clusters: []Cluster{
			{Name: "Server1", Hosts: []string{"http:node1:1234"}},
			{Name: "Server2", Hosts: []string{"http://node2:5678"}},
		},
		CurrentCluster: "Server2",
	}

	c.w = func(s string, b []byte, fm stdlib_fs.FileMode) error {
		expected, err := yaml.Marshal(c)
		assert.Nil(t, err, "Unexpected yaml marshal error")
		assert.Equal(t, expected, b, "Config data mismatch")
		return nil
	}

	err := c.DeleteCluster("Server2")

	require.NoError(t, err, "Unexpected set active cluster error")
	assert.Len(t, c.Clusters, 1, "Incorrect count of clusters")
	assert.Equal(t, "Server1", c.Clusters[0].Name, "Incorrect cluster name")
	assert.Equal(t, "", c.CurrentCluster, "Incorrect current cluster name")
}

func TestClusterConfig_DeleteClusterForInvalidClusterName(t *testing.T) {
	c := &ClusterConfig{
		Clusters: []Cluster{
			{Name: "Server1", Hosts: []string{"http:node1:1234"}},
			{Name: "Server2", Hosts: []string{"http://node2:5678"}},
		},
	}

	c.w = func(s string, b []byte, fm stdlib_fs.FileMode) error {
		expected, err := yaml.Marshal(c)
		assert.Nil(t, err, "Unexpected yaml marshal error")
		assert.Equal(t, expected, b, "Config data mismatch")
		return nil
	}

	err := c.DeleteCluster("Server3")
	assert.EqualError(t, err, "cluster Server3 not found", "Incorrect delete cluster error")
	assert.Len(t, c.Clusters, 2, "Incorrect count of clusters")
}

func TestClusterConfig_GetCurrentCluster(t *testing.T) {
	s1 := Cluster{Name: "Server1", Hosts: []string{"http:node1:1234"}}
	s2 := Cluster{Name: "Server2", Hosts: []string{"http://node2:5678"}}

	c := &ClusterConfig{
		Clusters:       []Cluster{s1, s2},
		CurrentCluster: "Server1",
	}

	actualCurrentCluster := c.GetCurrentCluster()
	expectedCurrentCluster := &s1

	assert.Equal(t, expectedCurrentCluster, actualCurrentCluster, "Unexpected current cluster config")
}

func TestClusterConfig_GetCurrentCluster_Return_Nil_when_not_Set(t *testing.T) {
	s1 := Cluster{Name: "Server1", Hosts: []string{"http:node1:1234"}}
	s2 := Cluster{Name: "Server2", Hosts: []string{"http://node2:5678"}}

	c := &ClusterConfig{
		Clusters: []Cluster{s1, s2},
	}

	var expected *Cluster = nil
	actualCurrentCluster := c.GetCurrentCluster()

	assert.Equal(t, expected, actualCurrentCluster, "Unexpected current cluster config")
}

func TestClusterConfig_UpdateCluster(t *testing.T) {
	c := &ClusterConfig{
		Clusters: []Cluster{
			{Name: "Server1", Hosts: []string{"http:node1:1234"}},
			{Name: "Server2", Hosts: []string{"http://node2:5678"}},
		},
		CurrentCluster: "Server2",
	}

	c.w = func(s string, b []byte, fm stdlib_fs.FileMode) error {
		expected, err := yaml.Marshal(c)
		assert.Nil(t, err, "Unexpected yaml marshal error")
		assert.Equal(t, expected, b, "Config data mismatch")
		return nil
	}

	err := c.UpdateCluster("Server2", "http:node2:1234,http:node2:5678")

	require.NoError(t, err, "Unexpected update cluster hosts error")
	assert.Len(t, c.Clusters, 2, "Incorrect count of clusters")
	assert.Equal(t, "Server2", c.Clusters[1].Name, "Incorrect cluster name")
	assert.Equal(t, []string{"http:node2:1234", "http:node2:5678"}, c.Clusters[1].Hosts, "Incorrect cluster hosts")
}

func TestClusterConfig_UpdateCluster_InvalidClusterName(t *testing.T) {
	c := &ClusterConfig{
		Clusters: []Cluster{
			{Name: "Server1", Hosts: []string{"http:node1:1234"}},
			{Name: "Server2", Hosts: []string{"http://node2:5678"}},
		},
		CurrentCluster: "Server2",
	}

	c.w = func(s string, b []byte, fm stdlib_fs.FileMode) error {
		expected, err := yaml.Marshal(c)
		assert.Nil(t, err, "Unexpected yaml marshal error")
		assert.Equal(t, expected, b, "Config data mismatch")
		return nil
	}

	err := c.UpdateCluster("Server3", "http:node2:1234,http:node2:5678")

	require.Error(t, err, "Missing update cluster hosts error")
	assert.EqualError(t, err, "cluster Server3 not found", "Incorrect update cluster error message")
}
