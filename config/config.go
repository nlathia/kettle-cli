package config

import (
	"io/ioutil"
	"path"
	"time"

	"gopkg.in/yaml.v2"
)

type TemplateConfig struct {
	CloudName     string    `yaml:"cloud"`
	DirectoryName string    `yaml:"name"`
	FunctionName  string    `yaml:"entrypoint"`
	Runtime       string    `yaml:"runtime"`
	Type          string    `yaml:"type"`
	Deployed      time.Time `yaml:"deployed_utc,omitempty"`
	PackageName   string    `yaml:"package,omitempty"`
}

func GetConfigFilePath(directoryPath string) string {
	return path.Join(directoryPath, DeploymentConfig)
}

func WriteConfig(config *TemplateConfig, directoryPath string) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	filePath := GetConfigFilePath(directoryPath)
	err = ioutil.WriteFile(filePath, []byte(data), 0644)
	if err != nil {
		return err
	}
	return nil
}

func ReadConfig(directoryPath string, values *TemplateConfig) error {
	filePath := GetConfigFilePath(directoryPath)
	contents, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(contents, values); err != nil {
		return err
	}
	return nil
}
