package config

import (
	"fmt"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

func setConfigPath() {
	// Find home directory.
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

func Read() {
	setConfigPath()
	viper.SetDefault(DeploymentType, GoogleCloudRun)
	viper.SetDefault(Runtime, Python)
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		// fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func Write() {
	setConfigPath()
	viper.WriteConfig()
}
