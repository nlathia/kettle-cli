package clouds

import (
	"fmt"
	"strings"
)

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
