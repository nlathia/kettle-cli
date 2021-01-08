package cmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/operatorai/operator/config"
)

// initCmd represents the command to set up and store preferences for the CLI tool
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Set up the operator CLI tool",
	Long: `The operator CLI tool supports multiple types of deployments: Google Cloud Functions, 
Cloud Run Containers, and AWS Lambda functions.

The init command allows you to set up your default preferences.`,
	RunE: runInit,
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func runInit(cmd *cobra.Command, args []string) error {
	configChoices := []*config.ConfigChoice{
		{
			// Pick a cloud provider
			Label: "Cloud Provider",
			Key:   config.CloudProvider,
			CollectValuesFunc: func() (map[string]string, error) {
				return config.CloudProviderNames, nil
			},
		},
		{
			// Pick a deployment type; assumes that the 'Pick a cloud provider'
			// step has already run or has been set via a flag
			Label: "Deployment type",
			Key:   config.DeploymentType,
			CollectValuesFunc: func() (map[string]string, error) {
				selectedCloud := viper.GetString(config.CloudProvider)
				if selectedCloud != "" {
					return config.DeploymentNames[selectedCloud], nil
				}
				return nil, errors.New(fmt.Sprintf("unknown cloud: %s", selectedCloud))
			},
		},
		{
			// Pick the default programming language; assumes that the 'Pick a deployment type'
			// step has already run or has been set via a flag
			Label: "Programming language",
			Key:   config.Runtime,
			CollectValuesFunc: func() (map[string]string, error) {
				selectedType := viper.GetString(config.DeploymentType)
				if selectedType != "" {
					return config.RuntimeNames[selectedType], nil
				}
				return nil, errors.New(fmt.Sprintf("unknown deployment type: %s", selectedType))
			},
		},
	}

	// Collect the remaining global preferences
	err := config.Collect(configChoices)
	if err != nil {
		return err
	}

	// Save the config
	config.Write()
	return nil
}

// mapContainsValue returns an error if a map doesn't contain a specific value
func mapContainsValue(value string, mapValues map[string]string) error {
	values := []string{}
	for _, mapValue := range mapValues {
		if mapValue == value {
			return nil
		}
		values = append(values, mapValue)
	}
	return errors.New(fmt.Sprintf("unknown value: %s (%s)", value, strings.Join(values, ", ")))
}
