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
	ProjectName string `json:"name"`
	Config      struct {
		Runtime        string `json:"runtime"`
		CloudProvider  string `json:"cloud_provider"`
		DeploymentType string `json:"deployment_type"`
		FunctionName   string `json:"entry_function"`
	} `json:"config"`
	Template []struct {
		Prompt string `json:"prompt"`
		Type   string `json:"type"`
		Key    string `json:"key"`
		Value  string `json:"value"`
	} `json:"template,omitempty"`
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

func WriteConfig(projectPath string, config *Template) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	configPath := path.Join(projectPath, configFileName)
	return ioutil.WriteFile(configPath, data, 0644)
}

func HasConfigFile(directory string) (bool, error) {
	configFilePath := path.Join(directory, configFileName)
	exists, err := PathExists(configFilePath)
	if err != nil {
		return false, err
	}
	if !exists {
		return false, nil
	}
	return true, nil
}
