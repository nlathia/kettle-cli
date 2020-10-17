package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"github.com/janeczku/go-spinner"
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

	s := spinner.StartNew("Collecting Google Cloud projects and regions...")
	cloudProjects, err := getGoogleCloudProjects()
	if err != nil {
		fmt.Printf("Unable to query for active projects: %v", err)
		return
	}
	if len(cloudProjects) == 0 {
		fmt.Printf("Could not find any active Google projects")
		return
	}

	deploymentRegions, err := getGoogleCloudRegions()
	if err != nil {
		fmt.Printf("Unable to query for deployment regions: %v", err)
		return
	}
	if len(deploymentRegions) == 0 {
		fmt.Printf("Could not find any active Google projects")
		return
	}
	s.Stop()

	configChoices := []configChoice{
		{
			// Pick a Google Cloud Project
			label:  "Google Cloud Project",
			values: cloudProjects,
			key:    config.ProjectID,
		},
		{
			// Pick a deployment region
			label:  "Deployment Region",
			values: deploymentRegions,
			key:    config.DeploymentRegion,
		},
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

func getGoogleCloudProjects() (map[string]string, error) {
	// Construct the gcloud command
	// gcloud projects list --format="json"
	commandArgs := []string{
		"projects",
		"list",
		"--format=\"json\"",
	}

	output, err := executeCommandWithResult("gcloud", commandArgs)
	if err != nil {
		return nil, err
	}

	type gcloudResults struct {
		ProjectID string `json:"projectId"`
		Name      string `json:"name"`
	}
	var results []gcloudResults
	if err := json.Unmarshal(output, &results); err != nil {
		return nil, err
	}

	projectIDs := map[string]string{}
	for _, project := range results {
		projectIDs[project.Name] = project.ProjectID
	}
	return projectIDs, nil
}

func getGoogleCloudRegions() (map[string]string, error) {
	// Construct the gcloud command
	// gcloud functions regions list --format="json"
	commandArgs := []string{
		"functions",
		"regions",
		"list",
		"--format=\"json\"",
	}

	output, err := executeCommandWithResult("gcloud", commandArgs)
	if err != nil {
		return nil, err
	}

	type gcloudResults struct {
		DisplayName string `json:"displayName"`
		LocationID  string `json:"locationId"`
	}
	var results []gcloudResults
	if err := json.Unmarshal(output, &results); err != nil {
		return nil, err
	}

	regions := map[string]string{}
	for _, region := range results {
		displayName := fmt.Sprintf("%s (%s)", region.DisplayName, region.LocationID)
		regions[displayName] = region.LocationID
	}
	return regions, nil
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
