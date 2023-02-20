package template

import (
	"bytes"
	"io/ioutil"
	"path"
	"text/template"
	yaml "gopkg.in/yaml.v2"
)

type TemplateParameters struct {
	RepoDir string
	Service string
	Release string
}

type FsConventions struct {
	Naming    string
	FluxcdDir string `yaml:"fluxcd_directory"`
	AppsDir   string `yaml:"apps_directory"`
}

type Orchestration struct {
	FsConventions *FsConventions 
	FluxcdFile    string
	AppFiles      []string
}

func loadFsConventions(fPath string, params *TemplateParameters) (*FsConventions, error) {
	var conv FsConventions
	
	f, err := ioutil.ReadFile(path.Join(fPath, "filesystem-conventions.yml"))
	if err != nil {
		return nil, err
	}

	tmpl, tErr := template.New("FsConventions").Parse(string(f))
	if tErr != nil {
		return nil, tErr
	}

	var b bytes.Buffer
	exErr := tmpl.Execute(&b, params)
	if exErr != nil {
		return nil, exErr
	}

	yamlErr := yaml.Unmarshal(b.Bytes(), &conv)
	if yamlErr != nil {
		return nil, yamlErr
	}

	return &conv, nil
}

func loadFluxcdFile(fPath string, params *TemplateParameters) (string, error) {
	f, err := ioutil.ReadFile(path.Join(fPath, "fluxcd.yml"))
	if err != nil {
		return "", err
	}

	tmpl, tErr := template.New("FluxcdFile").Parse(string(f))
	if tErr != nil {
		return "", tErr
	}

	var b bytes.Buffer
	exErr := tmpl.Execute(&b, params)
	if exErr != nil {
		return "", exErr
	}


	return string(b.Bytes()), nil
}

func loadAppFiles(aPath string, params *TemplateParameters) ([]string, error) {
	//path.Join(path, "app")
	return []string{}, nil
}

func processTemplatePath(tPath string, params *TemplateParameters) (string, error) {
	tmpl, tErr := template.New("string").Parse(tPath)
	if tErr != nil {
		return "", tErr
	}

	var b bytes.Buffer
	exErr := tmpl.Execute(&b, params)
	if exErr != nil {
		return "", exErr
	}

	return string(b.Bytes()), nil
}

func LoadTemplate(tPath string, params *TemplateParameters) (*Orchestration, error) {
	var o Orchestration

	var tPathErr error
	tPath, tPathErr = processTemplatePath(tPath, params)
	if tPathErr != nil {
		return nil, tPathErr
	}

	var fsCoErr error
	o.FsConventions, fsCoErr = loadFsConventions(tPath, params)
	if fsCoErr != nil {
		return nil, fsCoErr
	}

	var flFiErr error
	o.FluxcdFile, flFiErr = loadFluxcdFile(tPath, params)
	if flFiErr != nil {
		return nil, flFiErr
	}

	var appFlsErr error
	o.AppFiles, appFlsErr = loadAppFiles(tPath, params)
	if appFlsErr != nil {
		return nil, appFlsErr
	}

	return &o, nil
}