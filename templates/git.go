package templates

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/operatorai/kettle-cli/cli"
)

func isGitRepository(templatePath string) bool {
	if strings.HasSuffix(templatePath, ".git") {
		if strings.HasPrefix(templatePath, "git") || strings.HasPrefix(templatePath, "http") {
			return true
		}
	}
	return false
}

func cloneRepository(url string) (string, error) {
	tempDirectory, err := ioutil.TempDir("", "kettle")
	if err != nil {
		return "", err
	}
	err = cli.Execute("git", []string{
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

	// Sparse checkout, to avoid cloning the entire kettle-templates
	// repository
	if err := cli.Execute("git", []string{
		"clone",
		"--depth", "1",
		"--filter=blob:none",
		"--sparse",
		"https://github.com/operatorai/kettle-templates",
		tempDirectory,
	}, "Searching for template..."); err != nil {
		return "", err
	}

	rootDir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	os.Chdir(tempDirectory)
	defer func() {
		// Return to the original root directory
		os.Chdir(rootDir)
	}()

	if err := cli.Execute("git", []string{
		"sparse-checkout",
		"init",
		"--cone",
	}, "Searching for template..."); err != nil {
		return "", err
	}
	if err := cli.Execute("git", []string{
		"sparse-checkout",
		"set",
		templateName,
	}, "Searching for template..."); err != nil {
		return "", err
	}

	// Sparse checkout returns empty if a directory does not exist
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
