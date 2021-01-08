package aws

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/manifoldco/promptui"
)

func executeCommand(command string, args []string, quiet bool) error {
	osCmd := exec.Command(command, args...)
	if !quiet {
		osCmd.Stderr = os.Stderr
		osCmd.Stdout = os.Stdout
	}
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

// getValue shows a prompt (using a map's keys) to the user and returns
// the value that is indexed at that key
func getValue(label string, values map[string]string) (string, error) {
	valueLabels := []string{}
	for valueLabel, _ := range values {
		valueLabels = append(valueLabels, valueLabel)
	}

	prompt := promptui.Select{
		Label: label,
		Items: valueLabels,
	}
	_, result, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return "", err
	}
	return values[result], nil
}

// mapContainsValue returns an error if a map doesn't contain a specific value
func mapContainsValue(value string, mapValues map[string]string) error {
	values := []string{}
	for _, mapValue := range mapValues {
		if mapValue == value {
			return nil
		}
		values = append(values, mapValue)
	}
	return errors.New(fmt.Sprintf("unknown value: %s (%s)", value, strings.Join(values, ", ")))
}
