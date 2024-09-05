package main

import (
	"fmt"
	"os/exec"
)

// MessageType represents the type or category of a message in the system.
type MessageType int

const (
	Danger MessageType = iota
	Success
	Warning
)

// colorise applies a color style to a given message string based on the provided
// MessageType. If no MessageType is provided, it uses the default style.
func colorise(message string, msgType ...MessageType) string {
	reset := "\x1B[0m"
	var style string

	if len(msgType) > 0 {
		switch msgType[0] {
		case Danger:
			style = "\x1B[1;31m"
		case Success:
			style = "\x1B[1;32m"
		case Warning:
			style = "\x1B[1;33m"
		default:
			style = reset
		}
	} else {
		style = reset
	}

	return fmt.Sprintf("%s%s%s", style, message, reset)
}

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
