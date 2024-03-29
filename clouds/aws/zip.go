package aws

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/operatorai/kettle-cli/cli"
	"github.com/operatorai/kettle-cli/config"
)

const (
	deploymentArchiveName = "deployment.zip"
	goBuildFileName       = "main"
)

func createDeploymentArchive(cfg *config.Config) (string, error) {
	// Remove any existing deployment package
	if err := removeDeploymentArchive(cfg); err != nil {
		return "", err
	}

	// Create a path to the deployment archive
	rootDir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	deploymentFile := path.Join(rootDir, deploymentArchiveName)

	switch {
	case strings.HasPrefix(cfg.Config.Runtime, "python"):
		// https://docs.aws.amazon.com/lambda/latest/dg/python-package.html
		if err := addPythonLambdaToArchive(deploymentFile, cfg); err != nil {
			return "", err
		}
	case strings.HasPrefix(cfg.Config.Runtime, "go"):
		// https://docs.aws.amazon.com/lambda/latest/dg/golang-package.html
		if err := addGoLambdaToArchive(deploymentFile, cfg); err != nil {
			return "", err
		}
	}
	return deploymentFile, nil
}

func removeDeploymentArchive(cfg *config.Config) error {
	if err := removeFile(deploymentArchiveName); err != nil {
		return err
	}
	if strings.HasPrefix(cfg.Config.Runtime, "go") {
		if err := removeFile(goBuildFileName); err != nil {
			return err
		}
	}
	return nil
}

func removeFile(fileName string) error {
	if _, err := os.Stat(fileName); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	return os.Remove(fileName)
}

func addPythonLambdaToArchive(deploymentFile string, cfg *config.Config) error {
	// Add the contents of the lambda function directory
	err := cli.Execute("zip", []string{
		"-g",
		deploymentArchiveName,
		"-r",
		".",
	}, "Adding code to the deployment archive")
	if err != nil {
		return err
	}

	// Python builds need to add the site-packages contents
	var sitePackages string
	switch cfg.Config.PythonManager {
	case "pyenv":
		sitePackages, err = getPyenvSitePackagesDirectory(cfg.Config.Runtime)
		if err != nil {
			return err
		}
	case "conda":
		sitePackages, err = getCondaSitePackagesDirectory(cfg.Config.Runtime)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown python_manager: %s", cfg.Config.PythonManager)
	}

	if _, err := os.Stat(sitePackages); !os.IsNotExist(err) {
		// Change to the directory where the site-packages are stored
		// So that we can add them to the zip file as a directory
		rootDir, err := os.Getwd()
		if err != nil {
			return err
		}

		os.Chdir(sitePackages)
		err = cli.Execute("zip", []string{
			"-r",
			deploymentFile,
			".",
		}, "Adding site-packages to the deployment archive")
		if err != nil {
			return err
		}

		// Return to root directory to deploy the .zip file
		os.Chdir(rootDir)
	}
	return nil
}

func getPyenvSitePackagesDirectory(pythonVersion string) (string, error) {
	pyenvRoot, err := cli.ExecuteWithResult("pyenv", []string{
		"root",
	}, "Finding pyenv root")
	if err != nil {
		return "", err
	}

	pyenvLocal, err := cli.ExecuteWithResult("pyenv", []string{
		"local",
	}, "Finding pyenv local version")
	if err != nil {
		return "", err
	}

	fmt.Println(fmt.Sprintf("🔒  Adding site-packages from the pyenv '%s' environment.", string(pyenvLocal)))
	return fmt.Sprintf("%s/versions/%s/lib/%s/site-packages/",
		strings.Trim(string(pyenvRoot), "\n"),
		strings.Trim(string(pyenvLocal), "\n"),
		pythonVersion,
	), nil
}

func getCondaSitePackagesDirectory(pythonVersion string) (string, error) {
	condaRoot, err := cli.ExecuteWithResult("conda", []string{
		"info",
		"--base",
	}, "Finding conda root")
	if err != nil {
		return "", err
	}

	// Assumes that the conda env is active
	condaLocal := os.Getenv("CONDA_DEFAULT_ENV")
	fmt.Println(fmt.Sprintf("🔒  Adding site-packages from the conda '%s' environment.", condaLocal))
	if condaLocal == "base" {
		useBaseConda := cli.PromptToConfirm("The conda base environment is active. Continue")
		if !useBaseConda {
			return "", errors.New("please activate the conda environment for your project before deploying")
		}
	}

	return fmt.Sprintf("%s/envs/%s/lib/%s/site-packages/",
		strings.Trim(string(condaRoot), "\n"),
		strings.Trim(condaLocal, "\n"),
		pythonVersion,
	), nil
}

func addGoLambdaToArchive(deploymentFile string, cfg *config.Config) error {
	// go get github.com/aws/aws-lambda-go/lambda
	err := cli.Execute("go", []string{
		"get",
		"./...",
	}, "Running go get ./...")
	if err != nil {
		return err
	}

	// Build the function for linux
	err = cli.Execute("env", []string{
		"GOOS=linux",
		"CGO_ENABLED=0",
		"go",
		"build",
		"-o", goBuildFileName,
	}, "Building Go binary for GOOS=linux")
	if err != nil {
		return err
	}

	// zip function.zip main
	return cli.Execute("zip", []string{
		deploymentFile,
		"main",
	}, "Adding Go binary to deployment archive")
}
