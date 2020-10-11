package clouds

import (
	"os"
	"os/exec"
)

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
