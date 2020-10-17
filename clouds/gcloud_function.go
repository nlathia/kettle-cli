package clouds

import (
	"fmt"

	"github.com/operatorai/operator/config"
)

type GoogleCloudFunction struct{}

func (GoogleCloudFunction) Deploy(directory string, config *config.TemplateConfig) error {
	// Construct the gcloud command
	commandArgs := []string{
		"functions",
		"deploy",
		config.Name, // The cloud function is named the same as the directory
		"--runtime", config.Runtime,
		"--trigger-http", // We only currently support http triggers
		fmt.Sprintf("--entry-point=%s", config.FunctionName),
		fmt.Sprintf("--region=%s", config.DeploymentRegion),
		"--allow-unauthenticated",

		// @TODO these could be configurable
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
	fmt.Println("🚢  Deploying ", config.Name, "as a Google Cloud function")
	fmt.Println("⏭  Entry point: ", config.FunctionName, fmt.Sprintf("(%s)", config.Runtime))
	return executeCommand("gcloud", commandArgs)
}
