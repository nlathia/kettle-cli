package templates

import (
	"log"
	"os"
	"path"
	"regexp"

	"github.com/iancoleman/strcase"
)

func removePunctuation(input, replaceWith string) (string, error) {
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		return "", err
	}
	return reg.ReplaceAllString(input, replaceWith), nil
}

// The cloud function name is derived from the directory name
func CreateFunctionName(args []string) string {
	functionName, err := removePunctuation(args[0], "-")
	if err != nil {
		log.Fatal(err)
	}
	return strcase.ToKebab(functionName)
}

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
