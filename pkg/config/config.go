package config

import (
	"fmt"
	"log"
	"path"
	"strings"

	survey "github.com/AlecAivazis/survey/v2"
	"github.com/esctl/esctl/pkg/fs"
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
func (c *ClusterConfig) write() error {
	yamlData, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("yaml marshall: %w", err)
	}

	err = c.w(c.cfgFile, yamlData, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (c *ClusterConfig) AddCluster() error {
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

	err = c.write()
	if err != nil {
		return fmt.Errorf("persist config: %w", err)
	}
	return nil
}

func (c *ClusterConfig) SetActive(name string) error {
	if len(c.Clusters) == 0 {
		return fmt.Errorf("no clusters found")
	}

	if name != "" {
		_, err := c.find(name)
		if err != nil {
			return err
		}
		c.CurrentCluster = name
	} else {
		prompt := &survey.Select{
			Message: "Choose name of the cluster:",
			Options: c.names(),
		}
		err := survey.AskOne(prompt, &name, survey.WithValidator(survey.Required))
		if err != nil {
			return fmt.Errorf("survey select: %w", err)
		}
		c.CurrentCluster = name
	}

	err := c.write()
	if err != nil {
		return fmt.Errorf("persist config: %w", err)
	}

	fmt.Printf("Set %v as active cluster\n", name)
	return nil
}

func (c *ClusterConfig) DeleteCluster(name string) error {
	if len(c.Clusters) == 0 {
		return fmt.Errorf("no clusters found")
	}

	if name != "" {
		_, err := c.find(name)
		if err != nil {
			return err
		}
		return c.delete(name)
	}

	prompt := &survey.Select{
		Message: "Choose name of the cluster:",
		Options: c.names(),
	}
	err := survey.AskOne(prompt, &name, survey.WithValidator(survey.Required))
	if err != nil {
		return fmt.Errorf("survey select: %w", err)
	}

	return c.delete(name)
}

func (c *ClusterConfig) names() (names []string) {
	names = make([]string, 0, len(c.Clusters))
	for _, cl := range c.Clusters {
		names = append(names, cl.Name)
	}
	return
}

func (c *ClusterConfig) find(name string) (*Cluster, error) {
	for _, cluster := range c.Clusters {
		if name == cluster.Name {
			return &cluster, nil
		}
	}
	return nil, fmt.Errorf("cluster %v not found", name)
}

func (c *ClusterConfig) delete(name string) error {
	clusters := make([]Cluster, 0, len(c.Clusters))
	for _, cluster := range c.Clusters {
		if cluster.Name != name {
			clusters = append(clusters, cluster)
		}
	}

	c.Clusters = clusters
	// Reset current cluster if it's deleted from config
	if c.CurrentCluster == name {
		c.CurrentCluster = ""
	}

	err := c.write()
	if err != nil {
		return fmt.Errorf("persist config: %w", err)
	}

	fmt.Printf("Deleted %v from cluster config\n", name)
	return nil
}
