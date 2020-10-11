package templates

import (
	"os"
	"path"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/operatorai/operator/config"
)

// The entry function's case will vary based on the language;
// Right now, we're only supporting Python so we use ToSnake()
// But, for example, for Go we could use ToCamel()
func CreateEntryFunctionName(args []string, runtime string) string {
	switch {
	case strings.Contains(runtime, config.Python):
		return strcase.ToSnake(args[0])
	default:
		// Currently unreachable, as the `runtime` args
		// is checked before starting
		return args[0]
	}
}

// The cloud function name is derived from the directory name
func CreateFunctionName(args []string) string {
	return strcase.ToKebab(args[0])
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
