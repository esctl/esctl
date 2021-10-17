package config

import (
	"fmt"
	"log"
	"path"
	"strings"

	"github.com/esctl/esctl/pkg/fs"
	"gopkg.in/AlecAivazis/survey.v1"
	yaml "gopkg.in/yaml.v3"
)

type ClusterConfig struct {
	CurrentCluster string    `yaml:"CurrentCluster"`
	Clusters       []Cluster `yaml:"Clusters"`
	cfgFile        string
}

func (cc ClusterConfig) String() string {
	return fmt.Sprintf("CurrentCluster: %v\nClusters:\n%v", cc.CurrentCluster, cc.Clusters)
}

type Cluster struct {
	Name  string   `yaml:"Name"`
	Hosts []string `yaml:"Hosts"`
}

func (c Cluster) String() string {
	return fmt.Sprintf("Name: %v  Hosts: %v\n", c.Name, c.Hosts)
}

func (c *ClusterConfig) Load(cfgFile string, readFn fs.ReadFn) error {

	if cfgFile == "" {
		cfgFile = path.Join(fs.HomeDir(), ".esctl")
	}

	yamlFile, err := readFn(cfgFile)

	if err != nil {
		return fmt.Errorf("error reading config file %v", err)
	}

	err = yaml.Unmarshal(yamlFile, c)

	if err != nil {
		return fmt.Errorf("error loading config %v", err)
	}
	c.cfgFile = cfgFile
	return nil
}

// Write current config to file
func (c *ClusterConfig) Write(writeFn fs.WriteFn) error {
	yamlData, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("error while Marshaling %v", err)
	}

	err = writeFn(c.cfgFile, yamlData, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (c *ClusterConfig) AddCluster() {
	var qs = []*survey.Question{
		{
			Name:      "name",
			Prompt:    &survey.Input{Message: "Enter name of the cluster: "},
			Validate:  survey.Required,
			Transform: survey.Title,
		},
		{
			Name:     "hosts",
			Validate: survey.Required,
			Prompt: &survey.Input{
				Message: "Enter comma separated list of hosts: ",
			},
		},
	}
	answers := struct {
		Name  string
		Hosts string
	}{}

	err := survey.Ask(qs, &answers)
	if err != nil {
		log.Fatal(err.Error())
	}

	hosts := strings.Split(answers.Hosts, ",")
	cluster := Cluster{
		Name:  answers.Name,
		Hosts: hosts,
	}
	c.Clusters = append(c.Clusters, cluster)
}
