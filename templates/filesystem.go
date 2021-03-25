package templates

import (
	"os"
	"path"
)

func getRelativeDirectory(directoryName string) (string, error) {
	root, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return path.Join(root, directoryName), nil
}

func pathExists(path string) (bool, error) {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
