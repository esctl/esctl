package config

import (
	"testing"

	"github.com/esctl/esctl/pkg/fs"
	"github.com/stretchr/testify/assert"
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
	type fields struct {
		CurrentCluster string
		Clusters       []Cluster
		cfgFile        string
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ClusterConfig{
				CurrentCluster: tt.fields.CurrentCluster,
				Clusters:       tt.fields.Clusters,
				cfgFile:        tt.fields.cfgFile,
			}
			c.Write()
		})
	}
}

func TestClusterConfig_AddCluster(t *testing.T) {
	type fields struct {
		CurrentCluster string
		Clusters       []Cluster
		cfgFile        string
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ClusterConfig{
				CurrentCluster: tt.fields.CurrentCluster,
				Clusters:       tt.fields.Clusters,
				cfgFile:        tt.fields.cfgFile,
			}
			c.AddCluster()
		})
	}
}
