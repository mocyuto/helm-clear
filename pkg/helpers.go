package pkg

import (
	"fmt"
	"os"
	"os/exec"
)

func execute(args []string) ([]byte, error) {
	cmd := exec.Command(os.Getenv("HELM_BIN"), args...)
	output, err := cmd.Output()
	if exitError, ok := err.(*exec.ExitError); ok {
		return output, fmt.Errorf("%s: %s", exitError.Error(), string(exitError.Stderr))
	}
	return output, err
}
