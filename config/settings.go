package config

import (
	"io/ioutil"
	"path"

	"github.com/mitchellh/go-homedir"
	"gopkg.in/yaml.v2"
)

func getSettingsFilePath() (string, error) {
	// Find the home directory.
	home, err := homedir.Dir()
	if err != nil {
		return "", err
	}
	return path.Join(home, ".operator.yaml"), nil
}

func WriteSettings(cfg *Settings) error {
	// Get the path to the settings file
	filePath, err := getSettingsFilePath()
	if err != nil {
		return err
	}

	// Marshal & write the data
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filePath, []byte(data), 0644)
	if err != nil {
		return err
	}
	return nil
}

func ReadSettings() (*Settings, error) {
	// Get the path to the settings file
	filePath, err := getSettingsFilePath()
	if err != nil {
		return nil, err
	}

	// Read and unmarshal the file
	contents, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	values := Settings{}
	if err := yaml.Unmarshal(contents, &values); err != nil {
		return nil, err
	}
	return &values, nil
}
