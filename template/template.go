package template

import (
	"bytes"
	"io/fs"
	"io/ioutil"
	"path"
	"path/filepath"
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
	AppFiles      map[string]string
}

func applyFile(fPath string, params *TemplateParameters) ([]byte, error) {
	f, err := ioutil.ReadFile(fPath)
	if err != nil {
		return []byte{}, err
	}

	tmpl, tErr := template.New("template").Parse(string(f))
	if tErr != nil {
		return []byte{}, tErr
	}

	var b bytes.Buffer
	exErr := tmpl.Execute(&b, params)
	if exErr != nil {
		return []byte{}, exErr
	}

	return b.Bytes(), nil
}

func loadFsConventions(fPath string, params *TemplateParameters) (*FsConventions, error) {
	var conv FsConventions
	
	res, err := applyFile(path.Join(fPath, "filesystem-conventions.yml"), params)
	if err != nil {
		return nil, err
	}

	yamlErr := yaml.Unmarshal(res, &conv)
	if yamlErr != nil {
		return nil, yamlErr
	}

	return &conv, nil
}

func loadFluxcdFile(fPath string, params *TemplateParameters) (string, error) {
	res, err := applyFile(path.Join(fPath, "fluxcd.yml"), params)
	if err != nil {
		return "", err
	}

	return string(res), nil
}

func loadAppFiles(aPath string, params *TemplateParameters) (map[string]string, error) {
	dir := path.Join(aPath, "app")

	app := map[string]string{}

	err := filepath.Walk(dir, func(fPath string, fInfo fs.FileInfo, fErr error) error {
		if fErr != nil {
			return fErr
		}

		if fInfo.IsDir() {
			return nil
		}

		res, appErr := applyFile(fPath, params)
		if appErr != nil {
			return appErr
		}	
		app[fPath] = string(res)

		return nil
	})

	return app, err
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