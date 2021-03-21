package templates

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/operatorai/kettle/command"
)

func GetTemplate(templatePath string) (string, bool, error) {
	// Match on a local path first
	exists, err := PathExists(templatePath)
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
	return "", false, errors.New(fmt.Sprintf("%s not found", templatePath))
}

func cloneRepository(url string) (string, error) {
	tempDirectory, err := ioutil.TempDir("", "kettle")
	if err != nil {
		return "", err
	}
	err = command.Execute("git", []string{
		"clone",
		url,
		tempDirectory,
	}, "Cloning template...")
	return tempDirectory, err
}

func searchTemplates(templateName string) (string, error) {
	tempDirectory, err := ioutil.TempDir("", "kettle-templates")
	if err != nil {
		return "", err
	}
	rootDir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	defer func() {
		// Return to the original root directory
		os.Chdir(rootDir)
	}()

	// Sparse checkout, to avoid cloning the entire kettle-templates
	// repository. This will return empty if the templateName does
	// not exist
	if err := command.Execute("git", []string{
		"clone",
		"--depth", "1",
		"--filter=blob:none",
		"--sparse",
		"https://github.com/nlathia/kettle-templates.git",
		tempDirectory,
	}, "Searching for template..."); err != nil {
		return "", err
	}
	if err := command.Execute("cd", []string{
		tempDirectory,
	}, "Searching for template..."); err != nil {
		return "", err
	}
	if err := command.Execute("git", []string{
		"sparse-checkout",
		"init",
		"--cone",
	}, "Searching for template..."); err != nil {
		return "", err
	}
	if err := command.Execute("git", []string{
		"sparse-checkout",
		"set",
		templateName,
	}, "Searching for template..."); err != nil {
		return "", err
	}

	tempDirectory = path.Join(tempDirectory, templateName)
	exists, err := PathExists(tempDirectory)
	if err != nil {
		return "", err
	}
	if !exists {
		return "", errors.New(fmt.Sprintf("%s not found", templateName))
	}
	return tempDirectory, nil
}
