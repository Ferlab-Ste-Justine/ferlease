package config

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"text/template"
	"time"
	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	Operation         string        `yaml:"-"`
	Service           string
	Release           string
	Environment       string
	Repo              string
	RepoDir           string        `yaml:"-"`
	Ref               string
	GitSshKey         string        `yaml:"git_ssh_key"`
	GitKnownKey       string        `yaml:"git_known_key"`
	TemplateDirectory string        `yaml:"template_directory"`
	CommitMessage     string        `yaml:"commit_message"`
	PushRetries       int64         `yaml:"push_retries"`
	PushRetryInterval time.Duration `yaml:"push_retry_interval"`
}

func renderStr(s string, c *Config) (string, error) {
	tmpl, tErr := template.New("string").Parse(s)
	if tErr != nil {
		return "", tErr
	}

	var b bytes.Buffer
	exErr := tmpl.Execute(&b, *c)
	if exErr != nil {
		return "", exErr
	}

	return string(b.Bytes()), nil
}


func expandPath(fpath string, homedir string) string {
	if strings.HasPrefix(fpath, "~/") {
		fpath = path.Join(homedir, fpath[2:])
	}

	return fpath
}

func GetConfig(path string, operation string) (*Config, error) {
	var c Config
	c.Operation = operation

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

	c.RepoDir = fmt.Sprintf("%s-%s", c.Service, c.Release)

	var str string

	str, err = renderStr(c.TemplateDirectory, &c)
	if err != nil {
		return nil, err
	}
	c.TemplateDirectory = str

	str, err = renderStr(c.CommitMessage, &c)
	if err != nil {
		return nil, err
	}
	c.CommitMessage = str

	return &c, nil
}