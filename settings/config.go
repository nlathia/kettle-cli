package config

// func GetConfigFilePath(directoryPath string) string {
// 	return path.Join(directoryPath, DeploymentConfig)
// }

// func WriteConfig(cfg *TemplateConfig, directoryPath string) error {
// 	data, err := yaml.Marshal(cfg)
// 	if err != nil {
// 		return err
// 	}

// 	filePath := GetConfigFilePath(directoryPath)
// 	err = ioutil.WriteFile(filePath, []byte(data), 0644)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// func ReadConfig(directoryPath string) (*TemplateConfig, error) {
// 	filePath := GetConfigFilePath(directoryPath)
// 	contents, err := ioutil.ReadFile(filePath)
// 	if err != nil {
// 		return nil, err
// 	}

// 	values := TemplateConfig{}
// 	if err := yaml.Unmarshal(contents, &values); err != nil {
// 		return nil, err
// 	}
// 	return &values, nil
// }