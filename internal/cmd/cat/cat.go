package cat

import (
	"fmt"
	"log"

	"github.com/esctl/esctl/pkg/config"
	"github.com/esctl/esctl/pkg/es"
	"github.com/spf13/cobra"
)

func NewCmd(c *config.Cluster) *cobra.Command {

	catCmd := &cobra.Command{
		Use:   "cat",
		Short: "",
	}
	catCmd.AddCommand(newClusterCatNodesCmd(c))

	return catCmd
}

func newClusterCatNodesCmd(c *config.Cluster) *cobra.Command {
	catNCmd := &cobra.Command{
		Use:   "nodes",
		Short: "Returns information about a cluster's nodes",
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
			r, err := client.CatNodes()
			if err != nil {
				log.Fatalf("Error getting cluster nodes status, reason=[%v]", err)
			}

			fmt.Printf("%v\n", r)
		},
	}

	return catNCmd
}
