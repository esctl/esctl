package cluster

import (
	"fmt"
	"log"

	"github.com/esctl/esctl/pkg/config"
	"github.com/spf13/cobra"
)

func NewCmd(cfg *config.ClusterConfig) *cobra.Command {

	clusterCmd := &cobra.Command{
		Use:   "cluster",
		Short: "",
	}
	clusterCmd.AddCommand(newClusterListCmd(cfg))
	clusterCmd.AddCommand(newClusterAddCmd(cfg))
	clusterCmd.AddCommand(newClusterSetActiveCmd(cfg))
	clusterCmd.AddCommand(newClusterDeleteCmd(cfg))
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
			err := cfg.Write()

			if err != nil {
				log.Fatalf("Error writing config file: %v", err)
			}
		},
	}
	return clusterAddCmd
}

func newClusterSetActiveCmd(cfg *config.ClusterConfig) *cobra.Command {

	clusterAddCmd := &cobra.Command{
		Use:   "set-active",
		Short: "sa",
		Run: func(cmd *cobra.Command, args []string) {
			name := ""
			if len(args) > 0 {
				name = args[0]
			}
			if err := cfg.SetActive(name); err != nil {
				log.Fatalf("error: setting active cluster: %v", err)
			}
		},
		Args: cobra.MaximumNArgs(1),
	}

	return clusterAddCmd
}

func newClusterDeleteCmd(cfg *config.ClusterConfig) *cobra.Command {

	clusterDeleteCmd := &cobra.Command{
		Use:   "delete",
		Short: "",
		Run: func(cmd *cobra.Command, args []string) {
			name := ""
			if len(args) > 0 {
				name = args[0]
			}
			if err := cfg.DeleteCluster(name); err != nil {
				log.Fatalf("Error deleting cluster: %v", err)
			}
		},
		Args: cobra.MaximumNArgs(1),
	}
	return clusterDeleteCmd
}
