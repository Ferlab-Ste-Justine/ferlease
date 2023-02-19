package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	Service           string
	Release           string
	Repo              string
	RepoBranch        string `yaml:"repo_branch"`
	GitSshKey         string `yaml:"git_ssh_key"`
	GitKnownKey       string `yaml:"git_known_key"`
	TemplateDirectory string `yaml:"template_directory"`
}

func GetConfig(path string) (*Config, error) {
	var c Config

	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error reading the configuration file: %s", err.Error()))
	}
	err = yaml.Unmarshal(b, &c)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error parsing the configuration file: %s", err.Error()))
	}

	return &c, nil
}