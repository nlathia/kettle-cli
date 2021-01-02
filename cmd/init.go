package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/operatorai/operator/clouds"
	"github.com/operatorai/operator/config"
	"github.com/operatorai/operator/preferences"
)

// initCmd represents the command to set up and store preferences for the CLI tool
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Set up the operator CLI tool",
	Long: `The operator CLI tool supports multiple types of deployments: Google Cloud Functions,
	 Cloud Run Containers, and AWS Lambda functions.

The init command allows you to set up your preferences.`,
	Run: runInit,
}

var configChoices = []preferences.ConfigChoice{
	{
		// Pick a cloud provider
		Label:           "Cloud Provider",
		Key:             config.CloudProvider,
		FlagKey:         "cloud",
		FlagDescription: "The cloud provider to use",
		CollectValuesFunc: func() (map[string]string, error) {
			return config.CloudProviderNames, nil
		},
		ValidationFunc: func(v string) error {
			return mapContainsValue(v, config.CloudProviderNames)
		},
	},
	{
		// Pick a deployment type; assumes that the 'Pick a cloud provider'
		// step has already run and can not be set via a flag (to make things simpler)
		Label: "Deployment type",
		Key:   config.DeploymentType,
		CollectValuesFunc: func() (map[string]string, error) {
			selectedCloud := viper.GetString(config.CloudProvider)
			if selectedCloud != "" {
				return config.DeploymentNames[selectedCloud], nil
			}
			return nil, errors.New(fmt.Sprintf("unknown cloud: %s", selectedCloud))
		},
		ValidationFunc: func(v string) error {
			selectedCloud := viper.GetString(config.CloudProvider)
			if selectedCloud != "" {
				return mapContainsValue(v, config.DeploymentNames[selectedCloud])
			}
			return errors.New(fmt.Sprintf("unknown cloud: %s", selectedCloud))
		},
	},
	{
		// Pick the default programming language
		Label:           "Programming language",
		Key:             config.Runtime,
		FlagKey:         "runtime",
		FlagDescription: "The function's runtime language",
		CollectValuesFunc: func() (map[string]string, error) {
			return config.RuntimeNames, nil
		},
		ValidationFunc: func(v string) error {
			return mapContainsValue(v, config.RuntimeNames)
		},
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Enable operator init to also work with flags
	for _, configChoice := range configChoices {
		if configChoice.FlagKey != "" {
			initCmd.Flags().StringVar(&configChoice.FlagValue, configChoice.FlagKey, "", configChoice.FlagDescription)
		}
	}

	// initCmd.Flags().StringVar(&initValues.Runtime, "runtime", "")
	// Google Cloud specific flags
	// initCmd.Flags().StringVar(&initValues.DeploymentRegion, "region", "", "The region to deploy to")
	// initCmd.Flags().StringVar(&initValues.ProjectID, "project-id", "", "The gcloud project use")
}

func runInit(cmd *cobra.Command, args []string) {
	// Iterate on the flags first, which are quicker to validate
	for _, choice := range configChoices {
		if choice.FlagValue != "" {
			// The user has input a value as a flag; so we validate & store it
			if err := choice.ValidationFunc(choice.FlagValue); err != nil {
				fmt.Printf("Error: %v", err)
				return
			}
			viper.Set(choice.Key, choice.FlagValue)
		}
	}

	// Collect the remaining global preferences
	err := preferences.Collect(configChoices)
	if err != nil {
		fmt.Printf("Error: %v", err)
		return
	}

	// Run the cloud-specific setup
	selectedDeploymentType := viper.GetString(config.DeploymentType)
	cloudProvider, err := clouds.GetCloudProvider(selectedDeploymentType)
	if err != nil {
		fmt.Printf("Error: %v", err)
		return
	}

	err = cloudProvider.Setup()
	if err != nil {
		fmt.Printf("Error: %v", err)
		return
	}

	// Does not use SafeWrite - overwrites everything
	config.Write()
}

// mapContainsValue returns an error if a map doesn't contain a specific value
func mapContainsValue(value string, mapValues map[string]string) error {
	values := make([]string, len(mapValues))
	for _, mapValue := range mapValues {
		if mapValue == value {
			return nil
		}
		values = append(values, mapValue)
	}
	return errors.New(fmt.Sprintf("unknown value: %s (%s)", value, values))
}
