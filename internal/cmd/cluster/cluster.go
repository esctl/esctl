package cluster

import (
	"log"

	"github.com/esctl/esctl/pkg/config"
	"github.com/esctl/esctl/pkg/es"
	"github.com/spf13/cobra"
)

func NewCmd(c *config.Cluster) *cobra.Command {

	clusterCmd := &cobra.Command{
		Use:   "cluster",
		Short: "",
	}
	clusterCmd.AddCommand(newClusterHealthCmd(c))

	return clusterCmd
}

func newClusterHealthCmd(c *config.Cluster) *cobra.Command {
	healthCmd := &cobra.Command{
		Use:   "health",
		Short: "",
		Run: func(cmd *cobra.Command, args []string) {

			if c == nil {
				log.Fatal(`Cluster config is nil, Please add and set an active cluster config
  esctl config add
  esctl config set-active 
`)
			}

			client, err := es.New(c)
			if err != nil {
				log.Fatalf("Error creating elastic search client, reason=[%v]", err)
			}

			r, err := client.GetHealth()
			if err != nil {
				log.Fatalf("Error getting cluster health, reason=[%v]", err)
			}

			r.Print()
		},
	}
	return healthCmd
}
