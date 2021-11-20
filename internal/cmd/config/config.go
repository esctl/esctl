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
	return configCmd
}

func newConfigListCmd(cfg *config.ClusterConfig) *cobra.Command {

	configListCmd := &cobra.Command{
		Use:   "list",
		Short: "",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("%v", cfg)
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
		Use:   "set-active",
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
		Use:   "delete",
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
