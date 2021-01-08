package config

import (
	"github.com/manifoldco/promptui"
	"github.com/spf13/viper"
)

func Collect(configChoices []*ConfigChoice) error {
	// Iterate on the flags first, which are quicker to validate
	for _, choice := range configChoices {
		if choice.FlagValue != "" {
			// The user has input a value as a flag; so we validate & store it
			if err := choice.ValidationFunc(choice.FlagValue); err != nil {
				return err
			}
			viper.Set(choice.Key, choice.FlagValue)
		}
	}

	// Iterate over all the choices
	for _, choice := range configChoices {
		if choice.FlagValue == "" {
			// The user has not input a value as a flag; we collect the
			// available values and show them as a prompt
			values, err := choice.CollectValuesFunc()
			if err != nil {
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
		return "", err
	}
	return values[result], nil
}
