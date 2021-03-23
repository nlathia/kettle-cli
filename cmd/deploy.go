package cmd

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/spf13/cobra"

	"github.com/operatorai/kettle-cli/clouds"
	"github.com/operatorai/kettle-cli/templates"
)

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Ship a project you have created from a kettle template",
	Long: `🚢 The kettle CLI tool can automatically deploy
 your projects to your cloud provider.`,
	Args: validateDeployArgs,
	RunE: runDeploy,
}

func init() {
	rootCmd.AddCommand(deployCmd)
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
	deploymentPath, err := getDeploymentPath(args)
	if err != nil {
		return formatError(err)
	}

	// Read the template's config
	templateConfig, err := templates.ReadConfig(deploymentPath)
	if err != nil {
		return formatError(err)
	}

	// Get the cloud provider & service type
	cloudProvider, err := clouds.GetCloudProvider(templateConfig.Config.CloudProvider)
	if err != nil {
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
	if err := service.Deploy(deploymentPath, templateConfig); err != nil {
		return formatError(err)
	}

	// Write the settings back (they may have been changed)
	// @TODO
	// _ = config.WriteSettings(deploymentConfig.Settings)

	fmt.Println("✅  Deployed!")
	return nil
}

func getDeploymentPath(args []string) (string, error) {
	// Deploys from the current working directory
	rootDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	rootDir = path.Clean(rootDir)
	exists, err := templates.HasConfigFile(rootDir)
	if err != nil {
		return "", err
	}
	if exists {
		return rootDir, nil
	}

	// Deploys from a directory relative to the current working directory
	deploymentPath, err := templates.GetRelativeDirectory(args[0])
	exists, err = templates.HasConfigFile(deploymentPath)
	if err != nil {
		return "", err
	}
	if exists {
		return deploymentPath, nil
	}

	return "", fmt.Errorf("could not find template config file in %s", args[0])
}
