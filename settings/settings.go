package settings

import (
	"fmt"
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
		if DebugMode {
			fmt.Println("\tSettings file does not exist")
		}
		return &Settings{}, nil
	}

	contents, err := os.ReadFile(settingsFile)
	if err != nil {
		return nil, err
	}

	stg := &Settings{}
	if err := yaml.Unmarshal(contents, &stg); err != nil {
		return nil, err
	}
	if DebugMode {
		fmt.Printf("\tLoaded settings from: %s\n", settingsFile)
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

	err = os.WriteFile(settingsFile, []byte(data), 0644)
	if err != nil {
		return err
	}

	if DebugMode {
		fmt.Printf("\tSettings written to: %s\n", settingsFile)
		fmt.Println(string(data))
	}
	return nil
}
