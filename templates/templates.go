package templates

import (
	"fmt"
	"os"
	"path"

	"github.com/operatorai/kettle-cli/config"
)

func GetTemplate(templatePath string) (string, bool, error) {
	// Match on a local path first
	exists, err := pathExists(templatePath)
	if err != nil {
		return "", false, err
	}
	if exists {
		return templatePath, false, nil
	}

	// Match against a github repo & clone the repo to a tmp directory
	if isGitRepository(templatePath) {
		tempDirectory, err := cloneRepository(templatePath)
		return tempDirectory, true, err
	}

	// Look for the template in the kettle-templates monorepo
	tempDirectory, err := searchTemplates(templatePath)
	if err != nil {
		return "", false, err
	}
	return tempDirectory, true, nil
}

func GetProject(args []string) (string, error) {
	// Deploys from the current working directory
	rootDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	rootDir = path.Clean(rootDir)
	exists, err := config.HasConfigFile(rootDir)
	if err != nil {
		return "", err
	}
	if exists {
		return rootDir, nil
	}

	// Deploys from a directory relative to the current working directory
	deploymentPath, err := getRelativeDirectory(args[0])
	exists, err = config.HasConfigFile(deploymentPath)
	if err != nil {
		return "", err
	}
	if exists {
		return deploymentPath, nil
	}

	return "", fmt.Errorf("could not find template config file in %s", args[0])
}

func NewProjectPath(path string) (string, error) {
	directoryPath, err := getRelativeDirectory(path)
	if err != nil {
		return "", err
	}

	// Validate that the function path does *not* already exist
	exists, err := pathExists(directoryPath)
	if err != nil {
		return "", err
	}
	if exists {
		return "", fmt.Errorf("directory already exists: %s", directoryPath)
	}
	return directoryPath, nil
}
