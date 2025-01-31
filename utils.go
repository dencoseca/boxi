package main

import (
	"os/exec"
)

// pluralise returns the string "s" if the count is not equal to 1, otherwise an
// empty string.
func pluralise(count int) string {
	if count == 1 {
		return ""
	}

	return "s"
}

// runCommand executes a command with the given arguments and returns the
// combined output and error status.
func runCommand(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)
	output, err := cmd.CombinedOutput()

	return string(output), err
}
