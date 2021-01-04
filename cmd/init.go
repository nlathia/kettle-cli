package cmd

import (
	"errors"
	"fmt"
	"strings"

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

The init command allows you to set up your default preferences.`,
	Run: runInit,
}

var configChoices = []*preferences.ConfigChoice{
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
		Label:           "Deployment type",
		Key:             config.DeploymentType,
		FlagKey:         "type",
		FlagDescription: "The type of deployment (function, run, lambda)",
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
	// This currently adds all available flags, without checking for consistency
	// E.g., using --cloud aws and a GCP config flag would still parse but would
	// fail on validation

	configChoiceLists := [][]*preferences.ConfigChoice{
		configChoices,           // Add global flags
		clouds.GCPConfigChoices, // Add GCP-specific flags
		clouds.AWSConfigChoices, // Add AWS-specific flags
	}

	for _, configChoiceList := range configChoiceLists {
		for _, configChoice := range configChoiceList {
			if configChoice.FlagKey != "" {
				initCmd.Flags().StringVar(&configChoice.FlagValue, configChoice.FlagKey, "", configChoice.FlagDescription)
			}
		}
	}
}

func runInit(cmd *cobra.Command, args []string) {
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

	// Save the config
	config.Write()
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
