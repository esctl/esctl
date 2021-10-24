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
	w              fs.WriteFn
	r              fs.ReadFn
}

func New(w fs.WriteFn, r fs.ReadFn) *ClusterConfig {
	return &ClusterConfig{
		w: w,
		r: r,
	}
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

func (c *ClusterConfig) Load(cfgFile string) error {

	if cfgFile == "" {
		cfgFile = path.Join(fs.HomeDir(), ".esctl")
	}

	yamlFile, err := c.r(cfgFile)

	if err != nil {
		return fmt.Errorf("reading config file %w", err)
	}

	err = yaml.Unmarshal(yamlFile, c)

	if err != nil {
		return fmt.Errorf("loading config %w", err)
	}
	c.cfgFile = cfgFile
	return nil
}

// Write current config to file
func (c *ClusterConfig) Write() error {
	yamlData, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("while Marshaling %w", err)
	}

	err = c.w(c.cfgFile, yamlData, 0644)
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

func (c *ClusterConfig) SetActive(name string) error {
	var found bool
	for _, v := range c.Clusters {
		if name == v.Name {
			found = true
			break
		}
	}
	if found {
		c.CurrentCluster = name
		return c.Write()
	}

	return fmt.Errorf("cluster %v not found", name)
}
