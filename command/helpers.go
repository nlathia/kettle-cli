package command

import (
	"os"
	"os/exec"
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
