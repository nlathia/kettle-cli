package settings

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

func ReadSettings() (*Settings, error) {
	stg := &Settings{}
	if _, err := os.Stat(settingsFileName); os.IsNotExist(err) {
		// Return empty settings
		return stg, nil
	}

	contents, err := ioutil.ReadFile(settingsFileName)
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(contents, &stg); err != nil {
		return nil, err
	}
	return stg, nil
}

func WriteSettings(stg *Settings) error {
	data, err := yaml.Marshal(stg)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(settingsFileName, []byte(data), 0644)
	if err != nil {
		return err
	}
	return nil
}
