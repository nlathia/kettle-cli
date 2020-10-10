package cmd

import (
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/operatorai/operator/config"
)

func getEntryFunctionName(args []string, runtime string) string {
	switch {
	case strings.Contains(runtime, config.Python):
		return strcase.ToSnake(args[0])
	case strings.Contains(runtime, config.GoLang):
		return strcase.ToCamel(args[0])
	default:
		// Currently unreachable, as the `runtime` args
		// is checked before starting
		return args[0]
	}
}

func getFunctionName(args []string) string {
	// The cloud function name is derived from the directory name
	return strcase.ToKebab(args[0])
}

func getDirectoryPath(args []string) (string, error) {
	root, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return path.Join(root, getFunctionName(args)), nil
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

func getDirectoryExists(args []string) (bool, error) {
	functionPath, err := getDirectoryPath(args)
	if err != nil {
		return false, err
	}
	return pathExists(functionPath)
}

func executeCommand(command string, args []string) error {
	osCmd := exec.Command(command, args...)
	osCmd.Stderr = os.Stderr
	osCmd.Stdout = os.Stdout
	if err := osCmd.Run(); err != nil {
		return err
	}
	return nil
}

func executeCommandWithResult(command string, args []string) ([]byte, error) {
	osCmd := exec.Command(command, args...)
	osCmd.Stderr = os.Stderr
	output, err := osCmd.Output()
	if err != nil {
		return nil, err
	}
	return output, nil
}
