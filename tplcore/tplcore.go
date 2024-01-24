package tplcore

import (
	"bytes"
	"io/fs"
	"io/ioutil"
	"path"
	"path/filepath"
	"strings"
	"text/template"
)

type TemplateParameters struct {
	Service      string
	Release      string
	Environment  string
	CustomParams map[string]string
}

func (params *TemplateParameters) LoadFile(fPath string) ([]byte, error) {
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

func (params *TemplateParameters) LoadDir(aPath string) (map[string]string, error) {
	dir := path.Join(aPath, "app")

	app := map[string]string{}

	err := filepath.Walk(dir, func(fPath string, fInfo fs.FileInfo, fErr error) error {
		if fErr != nil {
			return fErr
		}

		if fInfo.IsDir() {
			return nil
		}

		res, appErr := params.LoadFile(fPath)
		if appErr != nil {
			return appErr
		}	
		app[strings.TrimPrefix(fPath, dir)] = string(res)

		return nil
	})

	return app, err
}