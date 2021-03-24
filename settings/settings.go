package settings

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

func ReadSettings() (*Settings, error) {
	if _, err := os.Stat(settingsFileName); os.IsNotExist(err) {
		// Return empty settings
		return &Settings{
			GoogleCloud: &GoogleCloudSettings{},
			AWS:         &AWSSettings{},
		}, nil
	}

	contents, err := ioutil.ReadFile(settingsFileName)
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
