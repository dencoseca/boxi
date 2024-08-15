package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type MessageType string

const (
	Danger  MessageType = "danger"
	Success MessageType = "success"
	Warning MessageType = "warning"
)

func printMessage(message string, msgType MessageType) {
	reset := "\x1B[0m"
	var style string
	switch msgType {
	case Danger:
		style = "\x1B[1;31m"
	case Success:
		style = "\x1B[1;32m"
	case Warning:
		style = "\x1B[1;33m"
	default:
		style = reset
	}
	fmt.Printf("%s%s%s\n", style, message, reset)
}

func runCommand(command string, args ...string) ([]byte, error) {
	cmd := exec.Command(command, args...)
	output, err := cmd.CombinedOutput()
	return output, err
}

func stopContainers() {
	output, err := runCommand("docker", "ps", "-a", "--format", "{{.Names}}")
	if err != nil {
		printMessage("Failed to list containers", Danger)
		return
	}
	containerNames := strings.Fields(string(output))
	if len(containerNames) > 0 {
		printMessage("STOPPING CONTAINERS", Success)
		for _, container := range containerNames {
			_, err = runCommand("docker", "stop", container)
			if err != nil {
				printMessage(fmt.Sprintf("Failed to stop container %s: %v", container, err), Danger)
			}
		}
	} else {
		printMessage("No CONTAINERS to STOP", Danger)
	}
}

func removeContainers() {
	output, err := runCommand("docker", "ps", "-a", "--format", "{{.Names}}")
	if err != nil {
		printMessage("Failed to list containers", Danger)
		return
	}
	containerNames := strings.Fields(string(output))
	if len(containerNames) > 0 {
		printMessage("REMOVING CONTAINERS", Success)
		for _, container := range containerNames {
			_, err = runCommand("docker", "rm", container)
			if err != nil {
				printMessage(fmt.Sprintf("Failed to remove container %s: %v", container, err), Danger)
			}
		}
	} else {
		printMessage("No CONTAINERS to REMOVE", Danger)
	}
}

func removeVolumes() {
	output, err := runCommand("docker", "volume", "ls", "-q")
	if err != nil {
		printMessage("Failed to list volumes", Danger)
		return
	}
	volumeNames := strings.Fields(string(output))
	if len(volumeNames) > 0 {
		printMessage("REMOVING VOLUMES", Success)
		for _, volume := range volumeNames {
			_, err = runCommand("docker", "volume", "rm", volume)
			if err != nil {
				printMessage(fmt.Sprintf("Failed to remove volume %s: %v", volume, err), Danger)
			}
		}
	} else {
		printMessage("No VOLUMES to REMOVE", Danger)
	}
}

func removeImages() {
	output, err := runCommand("docker", "images", "-q")
	if err != nil {
		printMessage("Failed to list images", Danger)
		return
	}
	imageIDs := strings.Fields(string(output))
	if len(imageIDs) > 0 {
		printMessage("REMOVING IMAGES", Success)
		for _, image := range imageIDs {
			_, err = runCommand("docker", "rmi", image)
			if err != nil {
				printMessage(fmt.Sprintf("Failed to remove image %s: %v", image, err), Danger)
			}
		}
	} else {
		printMessage("No IMAGES to REMOVE", Danger)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: boxi <command>")
		return
	}
	command := os.Args[1]
	switch command {
	case "sc":
		stopContainers()
	case "rc":
		removeContainers()
	case "rv":
		removeVolumes()
	case "ri":
		removeImages()
	case "cc":
		stopContainers()
		removeContainers()
	case "ca":
		stopContainers()
		removeContainers()
		removeVolumes()
	case "kill-it-with-fire-before-it-lays-eggs":
		stopContainers()
		removeContainers()
		removeVolumes()
		removeImages()
		printMessage("Pruning SYSTEM", Success)
		_, err := runCommand("docker", "system", "prune", "-f")
		if err != nil {
			printMessage("Failed to prune system", Danger)
		}
	default:
		fmt.Println("Unknown command:", command)
	}
}
