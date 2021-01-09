package command

import (
	"fmt"
	"os"
	"os/exec"
	"sort"

	"github.com/manifoldco/promptui"
)

func Execute(command string, args []string, quiet bool) error {
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

func ExecuteWithResult(command string, args []string) ([]byte, error) {
	osCmd := exec.Command(command, args...)
	osCmd.Stderr = os.Stderr
	output, err := osCmd.Output()
	if err != nil {
		return nil, err
	}
	return output, nil
}

// PromptForValue shows a prompt (using a map's keys) to the user and returns
// the value that is indexed at that key
func PromptForValue(label string, values map[string]string) (string, error) {
	valueLabels := []string{}
	for valueLabel, _ := range values {
		valueLabels = append(valueLabels, valueLabel)
	}
	sort.Strings(valueLabels)

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
