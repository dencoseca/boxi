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

func printMessage(message string, msgType ...MessageType) {
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
		stoppedContainerCount := 0
		for _, container := range containerNames {
			_, err = runCommand("docker", "stop", container)
			if err != nil {
				printMessage(fmt.Sprintf("Failed to stop container %s: %v", container, err), Danger)
			} else {
				stoppedContainerCount++
			}
		}
		message := fmt.Sprintf("Stopped %d container%s", stoppedContainerCount, func(n int) string {
			if n == 1 {
				return ""
			}
			return "s"
		}(stoppedContainerCount))
		printMessage(message)
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
		removedContainerCount := 0
		for _, container := range containerNames {
			_, err = runCommand("docker", "rm", container)
			if err != nil {
				printMessage(fmt.Sprintf("Failed to remove container %s: %v", container, err), Danger)
			} else {
				removedContainerCount++
			}
		}
		message := fmt.Sprintf("Removed %d container%s", removedContainerCount, func(n int) string {
			if n == 1 {
				return ""
			}
			return "s"
		}(removedContainerCount))
		printMessage(message)
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
		removedVolumeCount := 0
		for _, volume := range volumeNames {
			_, err = runCommand("docker", "volume", "rm", volume)
			if err != nil {
				printMessage(fmt.Sprintf("Failed to remove volume %s: %v", volume, err), Danger)
			} else {
				removedVolumeCount++
			}
		}
		message := fmt.Sprintf("Removed %d volume%s", removedVolumeCount, func(n int) string {
			if n == 1 {
				return ""
			}
			return "s"
		}(removedVolumeCount))
		printMessage(message)
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
		removedImageCount := 0
		for _, image := range imageIDs {
			_, err = runCommand("docker", "rmi", image)
			if err != nil {
				printMessage(fmt.Sprintf("Failed to remove image %s: %v", image, err), Danger)
			} else {
				removedImageCount++
			}
		}
		message := fmt.Sprintf("Stopped %d image%s", removedImageCount, func(n int) string {
			if n == 1 {
				return ""
			}
			return "s"
		}(removedImageCount))
		printMessage(message)
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
	case "purge":
		stopContainers()
		removeContainers()
		removeVolumes()
		removeImages()
		pruneSystem()
	default:
		fmt.Println("Unknown command:", command)
	}
}

func pruneSystem() {
	output, err := runCommand("docker", "system", "prune", "-f")
	if err != nil {
		printMessage("Pruning SYSTEM", Success)
		printMessage("Failed to prune system", Danger)
	} else {
		lines := strings.Split(string(output), "\n")
		var reclaimedSpace string
		for _, line := range lines {
			if strings.Contains(line, "Total reclaimed space") {
				reclaimedSpace = line
				break
			}
		}
		if reclaimedSpace != "Total reclaimed space: 0B" {
			printMessage("Pruning SYSTEM", Success)
			printMessage(reclaimedSpace)
		} else {
			printMessage("NOTHING to PRUNE", Danger)
		}
	}
}
