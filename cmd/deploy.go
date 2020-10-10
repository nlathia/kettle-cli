package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/operatorai/operator/config"
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

var templateConfig *config.TemplateConfig

func init() {
	rootCmd.AddCommand(deployCmd)
}

func deployArgs(cmd *cobra.Command, args []string) error {
	functionPath, err := getDirectoryPath(args)
	if err != nil {
		return err
	}

	// Validate that the function path argument exists
	exists, err := pathExists(functionPath)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("error validating path")
	}

	// Validate that a config file exists
	configFilePath := config.GetConfigFilePath(functionPath)
	exists, err = pathExists(configFilePath)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("operator.config file is missing")
	}

	// Read the config
	err = config.ReadConfig(configFilePath, templateConfig)
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
	functionPath, err := getDirectoryPath(args)
	if err != nil {
		return err
	}
	os.Chdir(functionPath)

	fmt.Println("🚢  Deploying ", templateConfig.DirectoryName, fmt.Sprintf("as an %s function", templateConfig.Type))
	fmt.Println("⏭  Entry point: ", templateConfig.FunctionName, fmt.Sprintf("(%s)", templateConfig.Runtime))

	// @TODO this needs to differentiate between a Cloud Function and Cloud Run
	// @Future support for other clouds/AWS
	// Construct the gcloud command
	commandArgs := []string{
		"functions",
		"deploy",
		templateConfig.DirectoryName,
		"--runtime", templateConfig.Runtime,
		fmt.Sprintf("--trigger-%s", templateConfig.Type),
		fmt.Sprintf("--entry-point=%s", templateConfig.FunctionName),
		"--region=europe-west2",
		"--allow-unauthenticated",
		// "--ignore-file=IGNORE_FILE",
		// "--egress-settings=EGRESS_SETTINGS",
		// "--ingress-settings=INGRESS_SETTINGS",
		// "--memory=MEMORY",
		// "--service-account=SERVICE_ACCOUNT",
		// "--source=SOURCE",
		// "--stage-bucket=STAGE_BUCKET",
		// "--timeout=TIMEOUT",
		// "--update-labels=[KEY=VALUE,…]",
		// "--env-vars-file=FILE_PATH",
		// "--max-instances=MAX_INSTANCES",
	}

	err = executeCommand("gcloud", commandArgs)
	if err != nil {
		return err
	}

	// Return to the original root directory
	os.Chdir(rootDir)
	return nil
}
