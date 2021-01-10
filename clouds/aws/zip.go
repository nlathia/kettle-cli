package aws

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/operatorai/operator/command"
	"github.com/operatorai/operator/config"
)

const (
	deploymentArchiveName = "deployment.zip"
)

func createDeploymentArchive(cfg *config.TemplateConfig) (string, error) {
	// Remove any existing deployment package
	if err := removeDeploymentArchiveIfExists(); err != nil {
		return "", err
	}

	// Store the current working directory before navigating away
	rootDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// Create the zip file, starting with the contents
	// of the current working directory
	deploymentFile := path.Join(rootDir, deploymentArchiveName)
	err = command.Execute("zip", []string{
		"-g",
		deploymentArchiveName,
		"-r",
		".",
	}, true)
	if err != nil {
		return "", err
	}

	if !strings.HasPrefix(cfg.Runtime, "python") {
		return deploymentArchiveName, nil
	}

	// Python builds need to add the site-packages contents
	sitePackages, err := getPyenvSitePackagesDirectory(cfg.Runtime)
	if err != nil {
		return "", err
	}

	if _, err := os.Stat(sitePackages); !os.IsNotExist(err) {
		// Change to the directory where the site-packages are stored
		// So that we can add them to the zip file as a directory
		os.Chdir(sitePackages)
		err = command.Execute("zip", []string{
			"-r",
			deploymentFile,
			".",
		}, true)
		if err != nil {
			return "", err
		}

		// Return to root directory to deploy the .zip file
		os.Chdir(rootDir)
	}
	return deploymentFile, nil
}

// removeDeploymentArchiveIfExists removes the deployment.zip file if present
func removeDeploymentArchiveIfExists() error {
	if _, err := os.Stat(deploymentArchiveName); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	return os.Remove(deploymentArchiveName)
}

func getPyenvSitePackagesDirectory(pythonVersion string) (string, error) {
	pyenvRoot, err := command.ExecuteWithResult("pyenv", []string{"root"})
	if err != nil {
		return "", err
	}

	pyenvLocal, err := command.ExecuteWithResult("pyenv", []string{"local"})
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/versions/%s/lib/%s/site-packages/",
		strings.Trim(string(pyenvRoot), "\n"),
		strings.Trim(string(pyenvLocal), "\n"),
		pythonVersion,
	), nil
}
