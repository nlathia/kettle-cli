package templates

import (
	"os"
	"path"
	"strings"
)

// Returns a path that is relative to the current working directory
func GetRelativeDirectory(directoryName string) (string, error) {
	root, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return path.Join(root, directoryName), nil
}

func PathExists(path string) (bool, error) {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// isGitRepository returns true if templatePath appears to be a git repo
func isGitRepository(templatePath string) bool {
	if strings.HasSuffix(templatePath, ".git") {
		if strings.HasSuffix(templatePath, "git") || strings.HasSuffix(templatePath, "http") {
			return true
		}
	}
	return false
}
