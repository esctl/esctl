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
	"fmt"
	"log"
	"os"

	"github.com/esctl/esctl/internal/cmd/cluster"
	"github.com/esctl/esctl/internal/cmd/root"
	"github.com/esctl/esctl/pkg/config"
	"github.com/esctl/esctl/pkg/fs"
	"github.com/spf13/cobra"
)

func main() {
	setup()
}

var cfg = &config.ClusterConfig{}
var cfgFile string

func setup() {

	rootCmd := root.NewCmd()
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.esctl.yaml)")

	err := cfg.Load(cfgFile, fs.Read)
	if err != nil {
		log.Fatal(err)
	}

	initSubCommands(rootCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}

func initSubCommands(rootCmd *cobra.Command) {
	rootCmd.AddCommand(cluster.NewCmd(cfg))
}
