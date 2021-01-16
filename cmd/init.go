package cmd

import (
	"github.com/operatorai/operator/clouds"
	"github.com/spf13/cobra"

	"github.com/operatorai/operator/command"
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
	settings := &config.Settings{}

	// Pick a cloud provider
	cloudProvider, err := command.PromptForValue("Cloud Provider", config.CloudProviderNames, false)
	if err != nil {
		return err
	}
	settings.CloudProvider = cloudProvider

	// Validate that the cloud provider is set up (e.g., the cli tool is installed)
	cloud, err := clouds.GetCloudProvider(cloudProvider)
	if err != nil {
		return err
	}
	if err := cloud.Setup(settings); err != nil {
		return err
	}

	// Pick a service
	availableDeploymentTypes := config.DeploymentNames[cloudProvider]
	deploymentType, err := command.PromptForValue("Deployment type", availableDeploymentTypes, false)
	if err != nil {
		return err
	}
	settings.DeploymentType = deploymentType

	// Pick a programming language
	availableLanguages := config.RuntimeNames[deploymentType]
	language, err := command.PromptForValue("Programming language", availableLanguages, false)
	if err != nil {
		return err
	}
	settings.Runtime = language

	// Save the settings
	config.WriteSettings(settings)
	return nil
}
