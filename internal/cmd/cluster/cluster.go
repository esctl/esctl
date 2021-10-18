package cluster

import (
	"fmt"
	"log"

	"github.com/esctl/esctl/pkg/config"
	"github.com/esctl/esctl/pkg/fs"
	"github.com/spf13/cobra"
)

func NewCmd(cfg *config.ClusterConfig) *cobra.Command {

	clusterCmd := &cobra.Command{
		Use:   "cluster",
		Short: "",
	}
	clusterCmd.AddCommand(newClusterListCmd(cfg))
	clusterCmd.AddCommand(newClusterAddCmd(cfg))
	return clusterCmd
}

func newClusterListCmd(cfg *config.ClusterConfig) *cobra.Command {

	clusterListCmd := &cobra.Command{
		Use:   "list",
		Short: "",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("%v", cfg)
		},
	}
	return clusterListCmd
}

func newClusterAddCmd(cfg *config.ClusterConfig) *cobra.Command {

	clusterAddCmd := &cobra.Command{
		Use:   "add",
		Short: "",
		Run: func(cmd *cobra.Command, args []string) {
			cfg.AddCluster()
			err := cfg.Write(fs.Write)

			if err != nil {
				log.Fatalf("Error writing config file %v", err)
			}
		},
	}
	return clusterAddCmd
}
