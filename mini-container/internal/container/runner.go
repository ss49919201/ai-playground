package container

import (
	"os/exec"
)

func RunCommand(program string, args []string) error {
	cmd := exec.Command(program, args...)
	return cmd.Run()
}

func RunCommandWithOutput(program string, args []string) (string, error) {
	cmd := exec.Command(program, args...)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}