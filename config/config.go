package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	Service           string
	Release           string
	Repo              string
	GitSshKey         string `yaml:"git_ssh_key"`
	GitKnownKey       string `yaml:"git_known_key"`
	TemplateDirectory string `yaml:"template_directory"`
}

func getConfigFilePath() string {
	path := os.Getenv("FERTURE_CONFIG_FILE")
	if path == "" {
	  return "config.yml"
	}
	return path
}

func GetConfig() (*Config, error) {
	var c Config

	b, err := ioutil.ReadFile(getConfigFilePath())
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error reading the configuration file: %s", err.Error()))
	}
	err = yaml.Unmarshal(b, &c)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error parsing the configuration file: %s", err.Error()))
	}

	return &c, nil
}