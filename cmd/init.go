package cmd

import (
	"encoding/json"
	"errors"
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

var initValues *config.TemplateConfig

func init() {
	rootCmd.AddCommand(initCmd)

	// Enable operator init to also work with flags
	initValues = &config.TemplateConfig{}
	initCmd.Flags().StringVar(&initValues.Type, "type", "", "The type of deployment to create")
	initCmd.Flags().StringVar(&initValues.Runtime, "runtime", "", "The function's runtime language")
	initCmd.Flags().StringVar(&initValues.DeploymentRegion, "region", "", "The region to deploy to")

	// Google Cloud specific flags
	initCmd.Flags().StringVar(&initValues.ProjectID, "project-id", "", "The gcloud project use")
}

func runInit(cmd *cobra.Command, args []string) {

	// A configChoice is defined as:
	type configChoice struct {
		// The label, which is shown in the prompt to the end user
		label string

		// The config key: the selection will be stored in viper using this
		key string

		// The flagValue, which has optionally been added by the user
		flagValue string

		// A function to collect values if the user does not provide one via a flag
		collectValuesFunc func() (map[string]string, error)

		// A function to validate the choice
		validationFunc func(string) error
	}

	configChoices := []configChoice{
		{
			// Pick a Google Cloud Project
			label:             "Google Cloud Project",
			key:               config.ProjectID,
			flagValue:         initValues.ProjectID,
			collectValuesFunc: getGoogleCloudProjects,
			validationFunc:    isActiveGoogleCloudProject,
		},
		{
			// Pick a deployment region
			label:             "Deployment Region",
			key:               config.DeploymentRegion,
			flagValue:         initValues.DeploymentRegion,
			collectValuesFunc: getGoogleCloudRegions,
			validationFunc:    isValidGoogleCloudRegion,
		},
		{
			// Pick the default deployment type
			label:     "Deployment type",
			key:       config.DeploymentType,
			flagValue: initValues.Type,
			collectValuesFunc: func() (map[string]string, error) {
				return config.DeploymentNames, nil
			},
			validationFunc: func(v string) error {
				if !config.DeploymentTypes.Contains(v) {
					return errors.New(fmt.Sprintf("unknown type: %s (%s)", v, config.DeploymentTypes))
				}
				return nil
			},
		},
		{
			// Pick the default programming language
			label:     "Programming language",
			key:       config.Runtime,
			flagValue: initValues.Runtime,
			collectValuesFunc: func() (map[string]string, error) {
				return config.RuntimeNames, nil
			},
			validationFunc: func(v string) error {
				if !config.Runtimes.Contains(v) {
					return errors.New(fmt.Sprintf("unknown runtime: %s (%s)", v, config.Runtimes))
				}
				return nil
			},
		},
	}

	// Iterate on the flags first, which are quicker to validate
	for _, choice := range configChoices {
		if choice.flagValue != "" {
			// The user has input a value as a flag; so we validate & store it
			if err := choice.validationFunc(choice.flagValue); err != nil {
				fmt.Printf("Error: %v", err)
				return
			}
			viper.Set(choice.key, choice.flagValue)
		}
	}

	// Iterate on all of the remaining choices second, since
	// it is slower to collect their values
	for _, choice := range configChoices {
		if choice.flagValue == "" {
			// The user has not input a value as a flag; we collect the
			// available values and show them as a prompt
			values, err := choice.collectValuesFunc()
			if err != nil {
				fmt.Printf("Error: %v", err)
				return
			}
			value, err := getValue(choice.label, values)
			if err != nil {
				return
			}
			viper.Set(choice.key, value)
		}
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

func getString(label string, validation func(string) error) (string, error) {
	prompt := promptui.Prompt{
		Label:    label,
		Validate: validation,
	}
	return prompt.Run()
}

func getGoogleCloudProjects() (map[string]string, error) {
	s := spinner.StartNew("Collecting Google Cloud projects...")
	defer s.Stop()

	// gcloud projects list --format="json"
	projectListLimit := 25
	commandArgs := []string{
		"projects",
		"list",
		"--format=\"json\"",
		fmt.Sprintf("--limit=%d", projectListLimit),
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
	if len(results) == projectListLimit {
		// Bail if the user has too many projects to 'reasonably' display
		return nil, errors.New(fmt.Sprintf("you have %d or more Google Cloud projects. "+
			"Please use operator init --project-id <id> to specify a project.", projectListLimit))
	}

	projectIDs := map[string]string{}
	for _, project := range results {
		displayName := fmt.Sprintf("%s (%s)", project.Name, project.ProjectID)
		projectIDs[displayName] = project.Name
	}
	return projectIDs, nil
}

func isActiveGoogleCloudProject(projectID string) error {
	s := spinner.StartNew(fmt.Sprintf("Checking Google Cloud project: %s...", projectID))
	defer s.Stop()

	// gcloud projects describe <id> --format="json"
	commandArgs := []string{
		"projects",
		"describe",
		projectID,
		"--format=\"json\"",
	}

	output, err := executeCommandWithResult("gcloud", commandArgs)
	if err != nil {
		return err
	}

	type gcloudResult struct {
		LifecycleState string `json:"lifecycleState"`
	}
	var result gcloudResult
	if err := json.Unmarshal(output, &result); err != nil {
		return err
	}

	if result.LifecycleState != "ACTIVE" {
		return errors.New("Project is not currently active")
	}
	return nil
}

func getGoogleCloudRegions() (map[string]string, error) {
	s := spinner.StartNew("Collecting Google Cloud regions...")
	defer s.Stop()

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

func isValidGoogleCloudRegion(locationID string) error {
	s := spinner.StartNew("Collecting Google Cloud regions...")
	defer s.Stop()

	// gcloud functions regions list --format="json"
	commandArgs := []string{
		"functions",
		"regions",
		"list",
		"--format=\"json\"",
	}

	output, err := executeCommandWithResult("gcloud", commandArgs)
	if err != nil {
		return err
	}

	type gcloudResults struct {
		LocationID string `json:"locationId"`
	}
	var results []gcloudResults
	if err := json.Unmarshal(output, &results); err != nil {
		return err
	}

	for _, region := range results {
		if region.LocationID == locationID {
			return nil
		}
	}
	return errors.New(fmt.Sprintf("unknown region: %s", locationID))
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
