package clouds

import (
	"fmt"

	"github.com/operatorai/operator/config"
)

type GoogleCloudFunction struct{}

func (GoogleCloudFunction) Build(directory string, config *config.TemplateConfig) error {
	fmt.Println("ℹ️  Google Cloud Functions do not need to be built.")
	return nil
}

func (GoogleCloudFunction) Deploy(directory string, config *config.TemplateConfig) error {
	// Construct the gcloud command
	commandArgs := []string{
		"functions",
		"deploy",
		config.DirectoryName, // The cloud function is named the same as the directory
		"--runtime", config.Runtime,
		"--trigger-http", // We only currently support http triggers
		fmt.Sprintf("--entry-point=%s", config.FunctionName),

		// @TODO these should be configurable
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
	fmt.Println("🚢  Deploying ", config.DirectoryName, fmt.Sprintf("as an %s function", config.Type))
	fmt.Println("⏭  Entry point: ", config.FunctionName, fmt.Sprintf("(%s)", config.Runtime))
	return executeCommand("gcloud", commandArgs)
}
