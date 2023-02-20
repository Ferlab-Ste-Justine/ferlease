package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
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

func expandPath(fpath string, homedir string) string {
	if strings.HasPrefix(fpath, "~/") {
		fpath = path.Join(homedir, fpath[2:])
	}

	return fpath
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

	homeDir, homeDirErr := os.UserHomeDir()
	if homeDirErr == nil {
		c.GitSshKey = expandPath(c.GitSshKey, homeDir)
		c.GitKnownKey = expandPath(c.GitKnownKey, homeDir)
		c.TemplateDirectory = expandPath(c.TemplateDirectory, homeDir)
	}

	return &c, nil
}