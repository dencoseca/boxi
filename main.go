package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

const mainUsage = `
Usage: 	boxi <command> [subCommand]

Clear up Docker resources

Commands:
  con, container, containers    Container commands
  vol, volume, volumes          Volume commands
  img, image, images            Image commands
  wipe                          Clean up containers and volumes
  purge                         Clean up containers, volumes, images, networks and the build cache

Run 'boxi <command> --help' for more information.`

func main() {
	if len(os.Args) < 2 {
		fmt.Println(mainUsage)
		os.Exit(1)
	}

	mainCommand := os.Args[1]

	switch mainCommand {
	case "-h", "--help":
		fmt.Println(mainUsage)
	case "con", "container", "containers":
		handleContainers()
	case "vol", "volume", "volumes":
		handleVolumes()
	case "img", "image", "images":
		handleImages()
	case "wipe":
		wipe()
	case "purge":
		purge()
	}
}

const containerUsage = `
Usage: 	boxi [con|container|containers] <command>

Clear up Docker container resources

Commands:
  stop     Stop all running containers
  rm       Remove all stopped containers
  clean    Stop and remove all running containers`

func handleContainers() {
	if len(os.Args) < 3 {
		fmt.Println(containerUsage)
		os.Exit(1)
	}

	subCommand := os.Args[2]

	switch subCommand {
	case "-h", "--help":
		fmt.Println(containerUsage)
	case "stop":
		stopContainers()
	case "rm":
		removeContainers()
	case "clean":
		stopContainers()
		removeContainers()
	default:
		fmt.Println(containerUsage)
		os.Exit(1)
	}
}

const volumeUsage = `
Usage: 	boxi [vol|volume|volumes] <command>

Clear up Docker volume resources

Commands:
  rm    Remove all dangling volumes`

func handleVolumes() {
	if len(os.Args) < 3 {
		fmt.Println(volumeUsage)
		os.Exit(1)
	}

	subCommand := os.Args[2]

	switch subCommand {
	case "-h", "--help":
		fmt.Println(volumeUsage)
	case "rm":
		removeVolumes()
	default:
		fmt.Println(volumeUsage)
		os.Exit(1)
	}
}

const imageUsage = `
Usage: 	boxi [img|image|images] <command>

Clear up Docker image resources

Commands:
  rm    Remove all dangling images`

func handleImages() {
	if len(os.Args) < 3 {
		fmt.Println(imageUsage)
		os.Exit(1)
	}

	subCommand := os.Args[2]

	switch subCommand {
	case "-h", "--help":
		fmt.Println(imageUsage)
	case "rm":
		removeImages()
	default:
		fmt.Println(imageUsage)
		os.Exit(1)
	}
}

func wipe() {
	stopContainers()
	removeContainers()
	removeVolumes()
}

func purge() {
	stopContainers()
	removeContainers()
	removeVolumes()
	removeImages()
	pruneSystem()
}

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
		log.Fatal(err)
	}

	containerNames := strings.Fields(string(output))
	if len(containerNames) == 0 {
		printMessage("No CONTAINERS to STOP", Danger)
		return
	}

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

	message := fmt.Sprintf("Stopped %d container%s", stoppedContainerCount, pluralise(stoppedContainerCount))
	printMessage(message)
}

func removeContainers() {
	output, err := runCommand("docker", "ps", "-a", "--format", "{{.Names}}")
	if err != nil {
		log.Fatal(err)
	}

	containerNames := strings.Fields(string(output))
	if len(containerNames) == 0 {
		printMessage("No CONTAINERS to REMOVE", Danger)
		return
	}

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

	message := fmt.Sprintf("Removed %d container%s", removedContainerCount, pluralise(removedContainerCount))
	printMessage(message)
}

func removeVolumes() {
	output, err := runCommand("docker", "volume", "ls", "-q")
	if err != nil {
		log.Fatal(err)
	}

	volumeIDs := strings.Fields(string(output))
	if len(volumeIDs) == 0 {
		printMessage("No VOLUMES to REMOVE", Danger)
		return
	}

	printMessage("REMOVING VOLUMES", Success)
	removedVolumeCount := 0

	for _, volume := range volumeIDs {
		_, err = runCommand("docker", "volume", "rm", volume)
		if err != nil {
			printMessage(fmt.Sprintf("Failed to remove volume %s: %v", volume, err), Danger)
		} else {
			removedVolumeCount++
		}
	}

	message := fmt.Sprintf("Removed %d volume%s", removedVolumeCount, pluralise(removedVolumeCount))
	printMessage(message)
}

func removeImages() {
	output, err := runCommand("docker", "images", "-q")
	if err != nil {
		log.Fatal(err)
	}

	imageIDs := strings.Fields(string(output))
	if len(imageIDs) == 0 {
		printMessage("No IMAGES to REMOVE", Danger)
		return
	}

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

	message := fmt.Sprintf("Stopped %d image%s", removedImageCount, pluralise(removedImageCount))
	printMessage(message)
}

func pruneSystem() {
	output, err := runCommand("docker", "system", "prune", "-f")
	if err != nil {
		log.Fatal(err)
	}

	lines := strings.Split(string(output), "\n")
	var reclaimedSpace string

	for _, line := range lines {
		if strings.Contains(line, "Total reclaimed space") {
			reclaimedSpace = line
			break
		}
	}

	if reclaimedSpace == "Total reclaimed space: 0B" {
		printMessage("NOTHING to PRUNE", Danger)
		return
	}

	printMessage("Pruning SYSTEM", Success)
	printMessage(reclaimedSpace)
}
