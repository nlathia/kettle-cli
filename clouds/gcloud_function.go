package clouds

import (
	"fmt"

	"github.com/operatorai/operator/config"
)

type GoogleCloudFunction struct{}

func (GoogleCloudFunction) Build(directory string, config *config.TemplateConfig) error {
	return nil
}

func (GoogleCloudFunction) Deploy(directory string, config *config.TemplateConfig) error {

	// Construct the gcloud command
	commandArgs := []string{
		"functions",
		"deploy",
		config.DirectoryName,
		"--runtime", config.Runtime,
		fmt.Sprintf("--trigger-%s", config.Type),
		fmt.Sprintf("--entry-point=%s", config.FunctionName),
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

	fmt.Println("🚢  Deploying ", templateConfig.DirectoryName, fmt.Sprintf("as an %s function", templateConfig.Type))
	fmt.Println("⏭  Entry point: ", templateConfig.FunctionName, fmt.Sprintf("(%s)", templateConfig.Runtime))
	return executeCommand("gcloud", commandArgs)
}
