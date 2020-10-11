package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/operatorai/operator/config"
)

// initCmd represents the command to set up and store preferences for the CLI tool
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Set up the operator CLI tool",
	Long: `The operator CLI tool supports multiple types of deployments: Google Cloud Functions
and Cloud Run Containers.

The init command allows you to set up your preferences.`,
	Run: runInit,
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func runInit(cmd *cobra.Command, args []string) {
	// A configChoice is defined as:
	// 1. The label, which is shown in the prompt
	// 2. The values (keys are shown in the prompt, values are stored in config)
	// 3. The config key (the selection will be stored in viper using this)
	type configChoice struct {
		label  string
		values map[string]string
		key    string
	}

	configChoices := []configChoice{
		{
			// Pick the default deployment type
			label:  "Deployment type",
			values: config.DeploymentNames,
			key:    config.DeploymentType,
		},
		{
			// Pick the default programming language
			label:  "Programming language",
			values: config.RuntimeNames,
			key:    config.Runtime,
		},
	}

	// Iterate on all of the choices
	for _, choice := range configChoices {
		value, err := getValue(choice.label, choice.values)
		if err != nil {
			fmt.Printf("Unknown value: %v\n", value)
			return
		}
		viper.Set(choice.key, value)
	}

	// Set the derived settings
	cloud, exists := config.CloudProviders[viper.GetString(config.DeploymentType)]
	if !exists {
		fmt.Printf("Unknown provider for: %v\n", viper.GetString(config.DeploymentType))
		return
	}
	viper.Set(config.CloudProvider, cloud)
	if cloud == config.GoogleCloud {
		// gcloud config get-value project
		projectID, err := getGoogleCloudProject()
		if err != nil {
			fmt.Printf("Unable to query for active project: %v", err)
			return
		}
		viper.Set(config.ProjectID, projectID)
	}

	// Does not use SafeWrite - overwrites everything
	config.Write()
}

func getValue(label string, values map[string]string) (string, error) {
	valueLabels := []string{}
	for valueLabel, _ := range values {
		valueLabels = append(valueLabels, valueLabel)
	}

	prompt := promptui.Select{
		Label: label,
		Items: valueLabels,
	}
	_, result, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return "", err
	}

	return values[result], nil
}

func getGoogleCloudProject() (string, error) {
	// Construct the gcloud command
	// gcloud config get-value project
	commandArgs := []string{
		"config",
		"get-value",
		"project",
	}

	fmt.Println("🔍  Querying for active gcloud project...")
	output, err := executeCommandWithResult("gcloud", commandArgs)
	if err != nil {
		return "", err
	}

	projectID := string(output)
	fmt.Println(fmt.Sprintf("✅  Using project: %s", projectID))
	return strings.Trim(string(output), "\n"), nil
}

func executeCommandWithResult(command string, args []string) ([]byte, error) {
	osCmd := exec.Command(command, args...)
	osCmd.Stderr = os.Stderr
	output, err := osCmd.Output()
	if err != nil {
		return nil, err
	}
	return output, nil
}
