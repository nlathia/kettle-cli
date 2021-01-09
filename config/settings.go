package config

import (
	"fmt"
	"os"
	"path"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

func getSettingsPath() string {
	// Find the home directory.
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return path.Join(home, ".operator.yaml")
}

func setSettingsPath() {
	// Find the home directory.
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Search config in home directory with name ".operator" (without extension).
	viper.AddConfigPath(home)
	viper.SetConfigName(".operator")
	viper.SetConfigType("yaml")
}

func ReadSettings() (*TemplateConfig, error) {
	setSettingsPath()
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	return &TemplateConfig{
		CloudProvider:  viper.GetString(CloudProvider),
		DeploymentType: viper.GetString(DeploymentType),
		Runtime:        viper.GetString(Runtime),
	}, nil
}

func WriteSettings() {
	configPath := getSettingsPath()
	if err := viper.SafeWriteConfigAs(configPath); err != nil {
		if os.IsNotExist(err) {
			_ = viper.WriteConfigAs(configPath)
		}
	}
	viper.WriteConfigAs(configPath)
}
