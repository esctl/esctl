/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"log"

	"github.com/esctl/esctl/internal/cmd/cluster"
	confCmd "github.com/esctl/esctl/internal/cmd/config"
	"github.com/esctl/esctl/internal/cmd/root"
	"github.com/esctl/esctl/pkg/config"
	"github.com/esctl/esctl/pkg/fs"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

func main() {
	setup()
}

var cfg = config.New(fs.Write, fs.Read)
var cfgFile string
var generateDocs bool

func setup() {

	rootCmd := root.NewCmd()
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.esctl.yaml)")

	rootCmd.PersistentFlags().BoolVar(&generateDocs, "generate-docs", false, "this option only work with source code")
	err := rootCmd.PersistentFlags().MarkHidden("generate-docs")
	if err != nil {
		log.Fatal(err)
	}

	err = cfg.Load(cfgFile)
	if err != nil {
		log.Fatalf("error loading config %v", err)
	}

	initSubCommands(rootCmd)
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
	if generateDocs {
		err = doc.GenMarkdownTree(rootCmd, "./docs")
		if err != nil {
			log.Fatal("Error generating markdown docs", err)
		}
	}
}

func initSubCommands(rootCmd *cobra.Command) {
	rootCmd.AddCommand(confCmd.NewCmd(cfg))
	rootCmd.AddCommand(cluster.NewCmd(cfg.GetCurrentCluster()))
}
