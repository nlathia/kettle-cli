package settings

import (
	"io/ioutil"
	"path"

	"gopkg.in/yaml.v2"
)

func GetSettingsFilePath(directoryPath string) string {
	return path.Join(directoryPath, settingsFileName)
}

func WriteSettings(stg *Settings, directoryPath string) error {
	data, err := yaml.Marshal(stg)
	if err != nil {
		return err
	}

	filePath := GetSettingsFilePath(directoryPath)
	err = ioutil.WriteFile(filePath, []byte(data), 0644)
	if err != nil {
		return err
	}
	return nil
}

func ReadSettings(directoryPath string) (*Settings, error) {
	filePath := GetSettingsFilePath(directoryPath)
	contents, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	stg := Settings{}
	if err := yaml.Unmarshal(contents, &stg); err != nil {
		return nil, err
	}
	return &stg, nil
}
