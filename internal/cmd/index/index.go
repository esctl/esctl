package index

import (
	"log"

	"github.com/esctl/esctl/pkg/config"
	"github.com/esctl/esctl/pkg/es"
	"github.com/spf13/cobra"
)

func NewCmd(c *config.Cluster) *cobra.Command {
	indexCmd := &cobra.Command{
		Use:   "index",
		Short: "",
	}
	indexCmd.AddCommand(newIndexListCmd(c))

	return indexCmd
}

func newIndexListCmd(c *config.Cluster) *cobra.Command {
	listCmd := &cobra.Command{
		Use:   "list",
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

			r, err := client.ListIndex()
			if err != nil {
				log.Fatalf("Error getting indices, reason=[%v]", err)
			}

			r.Print()
		},
	}
	return listCmd
}
