package settings

import (
	"io/ioutil"
	"os"
	"path"

	"github.com/mitchellh/go-homedir"
	"gopkg.in/yaml.v2"
)

func getSettingsFilePath() (string, error) {
	home, err := homedir.Dir()
	if err != nil {
		return "", err
	}
	return path.Join(home, ".kettle.yaml"), nil
}

func ReadSettings() (*Settings, error) {
	settingsFile, err := getSettingsFilePath()
	if err != nil {
		return nil, err
	}
	if _, err := os.Stat(settingsFile); os.IsNotExist(err) {
		// Return empty settings
		return &Settings{
			GoogleCloud: &GoogleCloudSettings{},
			AWS:         &AWSSettings{},
		}, nil
	}

	contents, err := ioutil.ReadFile(settingsFile)
	if err != nil {
		return nil, err
	}

	stg := &Settings{}
	if err := yaml.Unmarshal(contents, &stg); err != nil {
		return nil, err
	}
	return stg, nil
}

func WriteSettings(stg *Settings) error {
	settingsFile, err := getSettingsFilePath()
	if err != nil {
		return err
	}

	data, err := yaml.Marshal(stg)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(settingsFile, []byte(data), 0644)
	if err != nil {
		return err
	}
	return nil
}
