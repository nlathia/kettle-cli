package templates

import (
	"encoding/json"
	"io/ioutil"
	"path"
)

const (
	configFileName = "operator.json"
)

type Template struct {
	Config struct {
		Runtime        string `json:"runtime"`
		CloudProvider  string `json:"cloud_provider"`
		DeploymentType string `json:"deployment_type"`
		FunctionName   string `json:"entry_function"`
	} `json:"config"`
	Template []struct {
		Prompt string `json:"prompt"`
		Type   string `json:"type"`
		Key    string `json:"key"`
	} `json:"template"`
}

func ReadConfig(templatePath string) (*Template, error) {
	configPath := path.Join(templatePath, configFileName)
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	template := &Template{}
	err = json.Unmarshal(data, template)
	if err != nil {
		return nil, err
	}
	return template, nil
}
