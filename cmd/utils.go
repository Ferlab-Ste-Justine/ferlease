package cmd

import (
	"os"
	"path"
	"strings"

	"ferlab/ferlease/config"
	"ferlab/ferlease/template"
)

func expandPath(fpath string, homedir string) string {
	if strings.HasPrefix(fpath, "~/") {
		fpath = path.Join(homedir, fpath[2:])
	}

	return fpath
}

func processConfig(p *template.TemplateParameters, c *config.Config) error {
	homeDir, homeDirErr := os.UserHomeDir()
	if homeDirErr == nil {
		c.GitSshKey = expandPath(c.GitSshKey, homeDir)
		c.GitKnownKey = expandPath(c.GitKnownKey, homeDir)
		c.TemplateDirectory = expandPath(c.TemplateDirectory, homeDir)
	}

	dir, tplErr := template.ExecuteString(c.TemplateDirectory, p)
	if tplErr != nil {
		return tplErr
	}

	c.TemplateDirectory = dir

	return nil
}