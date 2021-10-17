package config

import (
	stdlib_fs "io/fs"
	"testing"

	"github.com/esctl/esctl/pkg/fs"
	"github.com/stretchr/testify/assert"
	yaml "gopkg.in/yaml.v3"
)

func TestClusterConfig_Load(t *testing.T) {
	type fields struct {
		CurrentCluster string
		Clusters       []Cluster
		cfgFile        string
	}

	type args struct {
		cfgFile string
		ReadFn  fs.ReadFn
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
			}
			c.Load(tt.args.cfgFile, tt.args.ReadFn)
			assert.Equal(t, "Test", c.Clusters[0].Name, "Cluster Names do not match")
			assert.Equal(t, "http://node1:1234", c.Clusters[0].Hosts[0], "Host do not match")
			assert.Equal(t, "http://node2:1234", c.Clusters[0].Hosts[1], "Host do not match")
			assert.Equal(t, "/some/config.yaml", c.cfgFile, "cfgFile do not match")
		})
	}
}

func TestClusterConfig_Write(t *testing.T) {
	type args struct {
		write fs.WriteFn
	}

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
	tests := []struct {
		name   string
		fields ClusterConfig
		args   args
	}{
		{
			name:   "should write config to file",
			fields: localClusterConfig,
			args: args{
				write: func(s string, b []byte, fm stdlib_fs.FileMode) error {
					expected, err := yaml.Marshal(localClusterConfig)
					assert.Nil(t, err, "Not expecting an error here")
					assert.Equal(t, expected, b, "config data mismatch")
					return nil
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
			}
			c.Write(tt.args.write)
		})
	}
}
