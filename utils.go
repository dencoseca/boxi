package main

import (
	"fmt"
	"os/exec"
)

type MessageType int

const (
	Danger MessageType = iota
	Success
	Warning
)

func pluralise(count int) string {
	if count == 1 {
		return ""
	}
	return "s"
}

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

func runCommand(command string, args ...string) ([]byte, error) {
	cmd := exec.Command(command, args...)
	output, err := cmd.CombinedOutput()
	return output, err
}
