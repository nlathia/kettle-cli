package cmd

import (
	"errors"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/spf13/cobra"

	"github.com/operatorai/operator/clouds"
	"github.com/operatorai/operator/config"
	"github.com/operatorai/operator/templates"
)

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Ship a function you have created",
	Long: `The operator CLI tool can automatically deploy
 a cloud function or GCP run project that you created with this tool.
	   
 The deploy command wraps the gsutil commands to simplify deployment.`,
	Args: validateDeployArgs,
	RunE: runDeploy,
}

var deploymentConfig *config.TemplateConfig
var deploymentPath string

func init() {
	rootCmd.AddCommand(deployCmd)
}

func validateDeployArgs(cmd *cobra.Command, args []string) error {
	// Validate that args exist
	if len(args) == 0 {
		return errors.New("please specify a path or directory name")
	}

	// Validate that the function path exists
	var err error
	deploymentPath, err = getDeploymentPath(args)
	if err != nil {
		return err
	}

	// Read the config
	deploymentConfig, err = config.ReadConfig(deploymentPath)
	if err != nil {
		return err
	}

	// Validate that the cloud provider is set up locally
	cloud, err := clouds.GetCloudProvider(deploymentConfig.CloudProvider)
	if err != nil {
		return err
	}
	if err := cloud.Setup(); err != nil {
		return err
	}
	return nil
}

// runDeploy creates or updates a cloud function
func runDeploy(cmd *cobra.Command, args []string) error {
	// Get the cloud provider & service type
	cloudProvider, err := clouds.GetCloudProvider(deploymentConfig.CloudProvider)
	if err != nil {
		return err
	}

	service, err := cloudProvider.GetService(deploymentConfig.DeploymentType)
	if err != nil {
		return err
	}

	// Store the current directory before changing away from it
	rootDir, err := os.Getwd()
	if err != nil {
		return err
	}

	// Change to the directory where the function to deploy is implemented
	// and run the deployment command
	os.Chdir(deploymentPath)
	if err := service.Deploy(deploymentPath, deploymentConfig); err != nil {
		return err
	}

	// Store that this function has been deployed
	deploymentConfig.Deployed = time.Now().UTC().String()
	config.WriteConfig(deploymentConfig, deploymentPath)

	// Return to the original root directory
	os.Chdir(rootDir)

	fmt.Println("✅  Deployed!")
	return nil
}

func getDeploymentPath(args []string) (string, error) {
	// operator deploy .
	// Deploys from the current working directory
	rootDir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	rootDir = path.Clean(rootDir)
	exists, err := directoryHasConfigFile(rootDir)
	if err != nil {
		return "", err
	}
	if exists {
		return rootDir, nil
	}

	// operator deploy /path/to/some/directory
	// Deploys from a fully formed path
	exists, err = directoryHasConfigFile(args[0])
	if err != nil {
		return "", err
	}
	if exists {
		return args[0], nil
	}

	// operator deploy some-directory
	// Deploys from a directory relative to the current working directory
	deploymentPath, err := templates.GetRelativeDirectory(args[0])
	exists, err = directoryHasConfigFile(deploymentPath)
	if err != nil {
		return "", err
	}
	if exists {
		return deploymentPath, nil
	}

	return "", fmt.Errorf("could not find %s file", config.DeploymentConfig)
}

func directoryHasConfigFile(directory string) (bool, error) {
	exists, err := templates.PathExists(directory)
	if err != nil {
		return false, err
	}
	if !exists {
		return false, nil
	}

	configFilePath := config.GetConfigFilePath(directory)
	exists, err = templates.PathExists(configFilePath)
	if err != nil {
		return false, err
	}
	if exists {
		return true, nil
	}
	return false, nil
}
