package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/operatorai/operator/config"
	"github.com/operatorai/operator/templates"
)

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Ship a function you have created",
	Long: `The operator CLI tool can automatically deploy
 a cloud function or GCP run project that you created with this tool.
	   
 The deploy command wraps the gsutil commands to simplify deployment.`,
	Args: deployArgs,
	RunE: runDeploy,
}

var deploymentConfig *config.TemplateConfig
var deploymentPath string

func init() {
	rootCmd.AddCommand(deployCmd)
}

func deployArgs(cmd *cobra.Command, args []string) error {
	// Validate that args exist
	if len(args) == 0 {
		return errors.New("please specify a directory name")
	}

	var err error
	deploymentPath, err = templates.GetRelativeDirectory(args[0])
	if err != nil {
		return err
	}

	// Validate that the function path exists
	exists, err := templates.PathExists(deploymentPath)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("error validating path")
	}

	// Validate that a config file exists
	configFilePath := config.GetConfigFilePath(deploymentPath)
	exists, err = templates.PathExists(configFilePath)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("%s file is missing", config.DeploymentConfig)
	}

	// Read the config
	err = config.ReadConfig(configFilePath, deploymentConfig)
	if err != nil {
		return err
	}
	return nil
}

// runDeploy creates or updates a cloud function
// https://cloud.google.com/sdk/gcloud/reference/functions/deploy
func runDeploy(cmd *cobra.Command, args []string) error {
	// We assume we are in the directory that is one level above the one with the functions
	// Store the current directory before changing away from it
	rootDir, err := os.Getwd()
	if err != nil {
		return err
	}

	// Change to the directory where the function to deploy is implemented
	// `gcloud functions deploy` assumes we are in this directory
	// functionPath, err := getDirectoryPath(args)
	// if err != nil {
	// 	return err
	// }
	// os.Chdir(functionPath)

	// Return to the original root directory
	os.Chdir(rootDir)
	return nil
}
