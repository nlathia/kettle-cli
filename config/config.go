package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

func ReadConfig(templatePath string) (*Config, error) {
	configPath := path.Join(templatePath, configFileName)
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	template := &Config{}
	err = json.Unmarshal(data, template)
	if err != nil {
		return nil, err
	}
	return template, nil
}

func WriteConfig(projectPath string, config *Config) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	configPath := path.Join(projectPath, configFileName)
	return ioutil.WriteFile(configPath, data, 0644)
}

func HasConfigFile(directory string) (bool, error) {
	configFilePath := path.Join(directory, configFileName)
	exists, err := pathExists(configFilePath)
	if err != nil {
		return false, err
	}
	if !exists {
		return false, nil
	}
	return true, nil
}

func pathExists(path string) (bool, error) {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func GetKey(cfg *Config, key string) (string, error) {
	for _, template := range cfg.Template {
		if template.Key == key {
			return template.Value, nil
		}
	}
	return "", fmt.Errorf("this template has not defined a '%s'", key)
}
