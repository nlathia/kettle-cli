package cli

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/operatorai/kettle-cli/settings"
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
	if settings.DebugMode {
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
