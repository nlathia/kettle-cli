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

	// Create a path to the deployment archive
	rootDir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	deploymentFile := path.Join(rootDir, deploymentArchiveName)

	switch {
	case strings.HasPrefix(cfg.Runtime, "python"):
		// https://docs.aws.amazon.com/lambda/latest/dg/python-package.html
		if err := addPythonLambdaToArchive(deploymentFile, cfg); err != nil {
			return "", err
		}
	case strings.HasPrefix(cfg.Runtime, "go"):
		// https://docs.aws.amazon.com/lambda/latest/dg/golang-package.html
		if err := addGoLambdaToArchive(deploymentFile, cfg); err != nil {
			return "", err
		}
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

func addPythonLambdaToArchive(deploymentFile string, cfg *config.TemplateConfig) error {
	// Add the contents of the lambda function directory
	err := command.Execute("zip", []string{
		"-g",
		deploymentArchiveName,
		"-r",
		".",
	})
	if err != nil {
		return err
	}

	// Python builds need to add the site-packages contents
	sitePackages, err := getPyenvSitePackagesDirectory(cfg.Runtime)
	if err != nil {
		return err
	}
	if _, err := os.Stat(sitePackages); !os.IsNotExist(err) {
		// Change to the directory where the site-packages are stored
		// So that we can add them to the zip file as a directory
		rootDir, err := os.Getwd()
		if err != nil {
			return err
		}

		os.Chdir(sitePackages)
		err = command.Execute("zip", []string{
			"-r",
			deploymentFile,
			".",
		})
		if err != nil {
			return err
		}

		// Return to root directory to deploy the .zip file
		os.Chdir(rootDir)
	}
	return nil
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

func addGoLambdaToArchive(deploymentFile string, cfg *config.TemplateConfig) error {
	// go get github.com/aws/aws-lambda-go/lambda
	err := command.Execute("go", []string{
		"get",
		"./...",
	})
	if err != nil {
		return err
	}

	// Build the function for linux
	err = command.Execute("env", []string{
		"GOOS=linux",
		"go",
		"build",
		"-o", "main",
		"./...",
	})
	if err != nil {
		return err
	}

	// zip function.zip main
	return command.Execute("zip", []string{
		deploymentFile,
		"main",
	})
}
