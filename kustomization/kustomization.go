package kustomization

import (
	_ "embed"
	"bytes"
	"io/ioutil"
	"text/template"
	"sort"
	yaml "gopkg.in/yaml.v2"
)

var (
	//go:embed kustomization.yaml
	kustomizationTemplate string
)

type Kustomization struct {
	ApiVersion string
	Kind       string
	Namespace  string
	Resources  []string
}

func (k *Kustomization) AddResource(resource string) {
	for _, r := range k.Resources {
		if r == resource {
			return
		}
	}

	k.Resources = append(k.Resources, resource)
	sort.Strings(k.Resources)
}

func (k *Kustomization) RemoveResource(resource string) {
	for idx, r := range k.Resources {
		if r == resource {
			k.Resources = append(k.Resources[:idx], k.Resources[idx + 1:]...)
			break
		}
	}
}

func (k *Kustomization) Render() (string, error) {
	tmpl, tErr := template.New("template").Parse(kustomizationTemplate)
	if tErr != nil {
		return "", tErr
	}

	var b bytes.Buffer
	exErr := tmpl.Execute(&b, *k)
	if exErr != nil {
		return "", exErr
	}

	return string(b.Bytes()), nil
}

func GetKustomization(fPath string) (*Kustomization, error) {
	var k Kustomization

	b, err := ioutil.ReadFile(fPath)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(b, &k)
	if err != nil {
		return nil, err
	}

	return &k, nil
}