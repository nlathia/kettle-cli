package command

import (
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/manifoldco/promptui"
	"github.com/operatorai/kettle/config"
)

func getSpinner(statusMessage string) *spinner.Spinner {
	s := spinner.New(spinner.CharSets[39], 100*time.Millisecond)
	s.Suffix = fmt.Sprintf("  %s...", statusMessage)
	s.Start()
	return s
}

func Execute(command string, args []string, statusMessage string) error {
	_, err := ExecuteWithResult(command, args, statusMessage)
	return err
}

func ExecuteWithResult(command string, args []string, statusMessage string) ([]byte, error) {
	osCmd := exec.Command(command, args...)
	if config.DebugMode {
		fmt.Println("\n", command, strings.Join(args, " "))
		osCmd.Stderr = os.Stderr
	} else {
		s := getSpinner(statusMessage)
		defer s.Stop()
	}

	output, err := osCmd.Output()
	if err != nil {
		return nil, err
	}
	return output, nil
}

// PromptForValue shows a prompt (using a map's keys) to the user and returns
// the value that is indexed at that key
func PromptForValue(label string, values map[string]string, addNoneOfThese bool) (string, error) {
	valueLabels := []string{}
	for valueLabel, _ := range values {
		valueLabels = append(valueLabels, valueLabel)
	}
	sort.Strings(valueLabels)
	if addNoneOfThese {
		valueLabels = append(valueLabels, config.PromptNoneOfTheseOption)
	}

	prompt := promptui.Select{
		Label: label,
		Items: valueLabels,
	}
	_, result, err := prompt.Run()
	if err != nil {
		return "", err
	}

	if addNoneOfThese && result == config.PromptNoneOfTheseOption {
		return "", nil
	}
	return values[result], nil
}

func PromptToConfirm(label string) bool {
	prompt := promptui.Prompt{
		Label:     label,
		IsConfirm: true,
	}

	result, err := prompt.Run()
	if err != nil {
		return false
	}

	if strings.ToLower(result) == "y" {
		return true
	}
	return false
}

func PromptForKeyValue(label string, values map[string]string) (string, string, error) {
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
		return "", "", err
	}
	return result, values[result], nil
}
