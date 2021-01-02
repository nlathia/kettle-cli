package preferences

import (
	"fmt"

	"github.com/manifoldco/promptui"
	"github.com/spf13/viper"
)

func Collect(configChoices []ConfigChoice) error {
	for _, choice := range configChoices {
		if choice.FlagValue == "" {
			// The user has not input a value as a flag; we collect the
			// available values and show them as a prompt
			values, err := choice.CollectValuesFunc()
			if err != nil {
				fmt.Printf("Error: %v", err)
				return err
			}
			value, err := getValue(choice.Label, values)
			if err != nil {
				return err
			}
			viper.Set(choice.Key, value)
		}
	}
	return nil
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
