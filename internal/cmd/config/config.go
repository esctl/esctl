package config

import (
	"fmt"
	"log"

	"github.com/esctl/esctl/pkg/config"
	"github.com/spf13/cobra"
)

func NewCmd(cfg *config.ClusterConfig) *cobra.Command {

	configCmd := &cobra.Command{
		Use:   "config",
		Short: "",
	}
	configCmd.AddCommand(newConfigListCmd(cfg))
	configCmd.AddCommand(newConfigAddCmd(cfg))
	configCmd.AddCommand(newConfigSetActiveCmd(cfg))
	configCmd.AddCommand(newConfigDeleteCmd(cfg))
	configCmd.AddCommand(newConfigUpdateCmd(cfg))
	return configCmd
}

func newConfigListCmd(cfg *config.ClusterConfig) *cobra.Command {

	configListCmd := &cobra.Command{
		Use:   "list",
		Short: "",
		Run: func(cmd *cobra.Command, args []string) {
			config.Print(cfg)
		},
	}
	return configListCmd
}

func newConfigAddCmd(cfg *config.ClusterConfig) *cobra.Command {

	configAddCmd := &cobra.Command{
		Use:   "add",
		Short: "",
		Run: func(cmd *cobra.Command, args []string) {
			err := cfg.AddCluster()
			if err != nil {
				log.Fatalf("Error adding cluster config: %v", err)
			}
		},
	}
	return configAddCmd
}

func newConfigSetActiveCmd(cfg *config.ClusterConfig) *cobra.Command {

	configAddCmd := &cobra.Command{
		Use:   "set-active [cluster_name]",
		Short: "sa",
		Run: func(cmd *cobra.Command, args []string) {
			name := ""
			if len(args) > 0 {
				name = args[0]
			}
			if err := cfg.SetActive(name); err != nil {
				log.Fatalf("Error setting active cluster config: %v", err)
			}
		},
		Args: cobra.MaximumNArgs(1),
	}

	return configAddCmd
}

func newConfigDeleteCmd(cfg *config.ClusterConfig) *cobra.Command {

	configDeleteCmd := &cobra.Command{
		Use:   "delete [cluster_name]",
		Short: "",
		Run: func(cmd *cobra.Command, args []string) {
			name := ""
			if len(args) > 0 {
				name = args[0]
			}
			if err := cfg.DeleteCluster(name); err != nil {
				log.Fatalf("Error deleting cluster config: %v", err)
			}
		},
		Args: cobra.MaximumNArgs(1),
	}
	return configDeleteCmd
}

func newConfigUpdateCmd(cfg *config.ClusterConfig) *cobra.Command {

	configUpdateCmd := &cobra.Command{
		Use:   "update [cluster_name] [new_host(s)]",
		Short: "",
		Run: func(cmd *cobra.Command, args []string) {
			name, newHosts := "", ""
			if len(args) > 1 {
				name = args[0]
				newHosts = args[1]
			}
			if err := cfg.UpdateCluster(name, newHosts); err != nil {
				log.Fatalf("Error updating cluster config: %v", err)
			}
		},
		// Accept positional arguments only if both cluster name and new host name(s) are present. If not, fall back to survey mode.
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 0 && len(args) < 2 {
				return fmt.Errorf("invalid usage. Please provide both [cluster_name] and [new_host(s)] as arguments or use survey mode")
			}
			return nil
		},
	}
	return configUpdateCmd
}
