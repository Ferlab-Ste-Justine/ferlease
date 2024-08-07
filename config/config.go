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

type AuthorConfig struct {
	Name  string
	Email string
}

type CommitSignatureConfig struct {
	Key        string
	Passphrase string
}

type GitAuthConfig struct {
	SshKey   string `yaml:"ssh_key"`
    KnownKey string `yaml:"known_key"`
}

type Orchestration struct {
	Type               string
	Repo               string
	Ref                string
	GitAuth            GitAuthConfig         `yaml:"git_auth"`
	CommitSignature    CommitSignatureConfig `yaml:"commit_signature"`
	CommitMessage      string                `yaml:"commit_message"`
	AcceptedSignatures string                `yaml:"accepted_signatures"`
	TemplateDirectory  string                `yaml:"template_directory"`
}

type Config struct {
	Operation          string                `yaml:"-"`
	RepoDir            string                `yaml:"-"`
	Service            string
	Release            string
	Environment        string
	CustomParams       map[string]string     `yaml:"custom_parameters"`
	Author             AuthorConfig
	CommitMessage      string                `yaml:"commit_message"`
	PushRetries        int64                 `yaml:"push_retries"`
	PushRetryInterval  time.Duration         `yaml:"push_retry_interval"`
	Orchestrations     []Orchestration
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
		for idx, _ := range c.Orchestrations {
			c.Orchestrations[idx].GitAuth.SshKey = expandPath(c.Orchestrations[idx].GitAuth.SshKey, homeDir)
			c.Orchestrations[idx].GitAuth.KnownKey = expandPath(c.Orchestrations[idx].GitAuth.KnownKey, homeDir)
			c.Orchestrations[idx].TemplateDirectory = expandPath(c.Orchestrations[idx].TemplateDirectory, homeDir)
		}
	}

	c.RepoDir = fmt.Sprintf("%s-%s", c.Service, c.Release)

	var str string

	for idx, _ := range c.Orchestrations {
		str, err = renderStr(c.Orchestrations[idx].TemplateDirectory, &c)
		if err != nil {
			return nil, err
		}
		c.Orchestrations[idx].TemplateDirectory = str
	}

	if c.CommitMessage != "" {
		str, err = renderStr(c.CommitMessage, &c)
		if err != nil {
			return nil, err
		}
		c.CommitMessage = str
	}

	for idx, _ := range c.Orchestrations {
		if c.Orchestrations[idx].CommitMessage == "" {
			if c.CommitMessage == "" {
				return nil, errors.New("Commit message for one of the orchestrations is empty. Either define a commit message for every orchestration or define a default commit message at the root of the configuration")
			}

			c.Orchestrations[idx].CommitMessage = c.CommitMessage
			continue
		}

		str, err = renderStr(c.Orchestrations[idx].CommitMessage, &c)
		if err != nil {
			return nil, err
		}
		c.Orchestrations[idx].CommitMessage = str
	}

	for _, orch := range c.Orchestrations {
		if orch.Type != "fluxcd" && orch.Type != "terraform" {
			return nil, errors.New(fmt.Sprintf("Orchestration of type '%s' is not recognized. Valid values are 'fluxcd' or 'terraform", orch.Type))
		}
	}

	return &c, nil
}