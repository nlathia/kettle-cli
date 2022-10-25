package cli

import (
	"sort"
	"strings"

	"github.com/manifoldco/promptui"
)

const (
	PromptNoneOfTheseOption = "None of these (create a new one)"
)

func PromptForValue(label string, values map[string]string, addNoneOfThese bool) (string, error) {
	valueLabels := []string{}
	for valueLabel, _ := range values {
		valueLabels = append(valueLabels, valueLabel)
	}

	sort.Strings(valueLabels)
	if addNoneOfThese {
		valueLabels = append(valueLabels, PromptNoneOfTheseOption)
	}

	prompt := promptui.Select{
		Label: label,
		Items: valueLabels,
	}
	_, result, err := prompt.Run()
	if err != nil {
		return "", err
	}

	if addNoneOfThese && result == PromptNoneOfTheseOption {
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

func PromptForString(label string) (string, error) {
	prompt := promptui.Prompt{
		Label: label,
	}

	result, err := prompt.Run()
	if err != nil {
		return "", err
	}
	return result, nil
}
