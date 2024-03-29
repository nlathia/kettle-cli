package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/operatorai/kettle-cli/clouds"
	"github.com/operatorai/kettle-cli/config"
	"github.com/operatorai/kettle-cli/settings"
	"github.com/operatorai/kettle-cli/templates"
)

var (
	environment string

	deployCmd = &cobra.Command{
		Use:   "deploy",
		Short: "Ship a project you have created from a kettle template",
		Long: `🚢 The kettle CLI tool can automatically deploy
 your projects to your cloud provider.`,
		Args: validateDeployArgs,
		RunE: runDeploy,
	}
)

func init() {
	rootCmd.AddCommand(deployCmd)
	deployCmd.Flags().StringVarP(&environment, "env", "e", "", "Environment to deploy to (GCP only)")
}

func validateDeployArgs(cmd *cobra.Command, args []string) error {
	// Validate that args exist
	if len(args) == 0 {
		return errors.New("please specify a path or directory name")
	}
	return nil
}

// runDeploy creates or updates a cloud function
func runDeploy(cmd *cobra.Command, args []string) error {
	// Construct the path we want to deploy from
	deploymentPath, err := templates.GetProject(args)
	if err != nil {
		return formatError(err)
	}

	// Read the template's config
	templateConfig, err := config.ReadConfig(deploymentPath)
	if err != nil {
		return formatError(err)
	}

	// Read global settings
	cloudSettings, err := settings.ReadSettings()
	if err != nil {
		return formatError(err)
	}

	// Get the cloud provider & service type
	cloudProvider, err := clouds.GetCloudProvider(templateConfig.Config.CloudProvider)
	if err != nil {
		return formatError(err)
	}

	// Set up the provider (if not done so already)
	if err := cloudProvider.Setup(cloudSettings, false); err != nil {
		return formatError(err)
	}

	service, err := cloudProvider.GetService(templateConfig.Config.DeploymentType)
	if err != nil {
		return formatError(err)
	}

	// Store the current directory before changing away from it
	rootDir, err := os.Getwd()
	if err != nil {
		return formatError(err)
	}

	// Change to the directory where the function to deploy is implemented
	// and run the deployment command
	os.Chdir(deploymentPath)
	defer func() {
		// Return to the original root directory
		os.Chdir(rootDir)
	}()

	// Deploy
	if err := service.Deploy(deploymentPath, templateConfig, cloudSettings, environment); err != nil {
		return formatError(err)
	}

	// Write the settings & config back (they may have been changed)
	if err := settings.WriteSettings(cloudSettings); err != nil {
		if settings.DebugMode {
			fmt.Println(err.Error())
		}
	}
	if err := config.WriteConfig(deploymentPath, templateConfig); err != nil {
		if settings.DebugMode {
			fmt.Println(err.Error())
		}
	}

	fmt.Println("✅  Deployed!")
	return nil
}
