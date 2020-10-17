package templates

import (
	"log"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/operatorai/operator/config"
)

func removePunctuation(input, replaceWith string) (string, error) {
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		return "", err
	}
	return reg.ReplaceAllString(input, replaceWith), nil
}

// The entry function's case will vary based on the language;
// Right now, we're only supporting Python so we use ToSnake()
func CreateEntryFunctionName(args []string, runtime string) string {
	switch {
	case strings.Contains(runtime, config.Python):
		entryName, err := removePunctuation(args[0], "_")
		if err != nil {
			log.Fatal(err)
		}
		return strcase.ToSnake(entryName)
	default:
		// Currently unreachable, as the `runtime` args
		// is checked before starting
		return args[0]
	}
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
