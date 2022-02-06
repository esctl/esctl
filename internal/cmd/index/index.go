package index

import (
	"fmt"
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
	indexCmd.AddCommand(newIndexSettingsCmd(c))

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

func newIndexSettingsCmd(c *config.Cluster) *cobra.Command {
	listCmd := &cobra.Command{
		Use:   "settings",
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

			allF, err := cmd.Flags().GetBool("all")
			if err != nil {
				log.Fatalf("Error fetching value for flag: `all`, reason=[%v]", err)
			}
			if allF {
				r, err := client.AllIndexSettings()
				if err != nil {
					log.Fatalf("Error getting index settings, reason=[%v]", err)
				}
				fmt.Printf("%v\n", r)
				return
			}

			if len(args) < 1 {
				log.Fatalf("Invalid usage, pass the `--all` flag to fetch all indices settings, or space seprated list of index names to filter by name\n")
			}

			r, err := client.IndexSettings(args)
			if err != nil {
				log.Fatalf("Error getting index settings, reason=[%v]", err)
			}
			fmt.Printf("%v\n", r)
		},
	}

	listCmd.Flags().BoolP("all", "a", false, "Fetch settings for all indices")
	return listCmd
}
