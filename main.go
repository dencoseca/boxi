package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

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

func stopContainers() {
	output, err := runCommand("docker", "ps", "-a", "--format", "{{.Names}}")
	if err != nil {
		log.Fatal(string(output))
	}

	containerNames := strings.Fields(string(output))
	if len(containerNames) == 0 {
		fmt.Println(colorise("No CONTAINERS to STOP", Warning))
		return
	}

	stoppedContainerCount := 0
	for _, container := range containerNames {
		output, err = runCommand("docker", "stop", container)
		if err != nil {
			fmt.Println(colorise(fmt.Sprintf("Failed to stop container %s: %s", container, output), Danger))
		} else {
			stoppedContainerCount++
		}
	}

	if stoppedContainerCount == 0 {
		fmt.Println(colorise("No CONTAINERS were STOPPED", Danger))
		return
	}

	fmt.Println(colorise(fmt.Sprintf("%s: Stopped %d container%s", colorise("STOPPING CONTAINERS", Success), stoppedContainerCount, pluralise(stoppedContainerCount))))
}

func removeContainers() {
	output, err := runCommand("docker", "ps", "-a", "--format", "{{.Names}}")
	if err != nil {
		log.Fatal(string(output))
	}

	containerNames := strings.Fields(string(output))
	if len(containerNames) == 0 {
		fmt.Println(colorise("No CONTAINERS to REMOVE", Warning))
		return
	}

	removedContainerCount := 0
	for _, container := range containerNames {
		output, err = runCommand("docker", "rm", container)
		if err != nil {
			fmt.Println(colorise(fmt.Sprintf("Failed to remove container %s: %s", container, output), Danger))
		} else {
			removedContainerCount++
		}
	}

	if removedContainerCount == 0 {
		fmt.Println(colorise("No CONTAINERS were REMOVED", Danger))
		return
	}

	fmt.Println(colorise(fmt.Sprintf("%s: Removed %d container%s", colorise("REMOVING CONTAINERS", Success), removedContainerCount, pluralise(removedContainerCount))))
}

func removeVolumes() {
	output, err := runCommand("docker", "volume", "ls", "-q")
	if err != nil {
		log.Fatal(string(output))
	}

	volumeIDs := strings.Fields(string(output))
	if len(volumeIDs) == 0 {
		fmt.Println(colorise("No VOLUMES to REMOVE", Warning))
		return
	}

	removedVolumeCount := 0
	for _, volume := range volumeIDs {
		output, err = runCommand("docker", "volume", "rm", volume)
		if err != nil {
			fmt.Println(colorise(fmt.Sprintf("Failed to remove volume %s: %s", volume, output), Danger))
		} else {
			removedVolumeCount++
		}
	}

	if removedVolumeCount == 0 {
		fmt.Println(colorise("No VOLUMES were REMOVED", Danger))
		return
	}

	fmt.Println(colorise(fmt.Sprintf("%s: Removed %d volume%s", colorise("REMOVING VOLUMES", Success), removedVolumeCount, pluralise(removedVolumeCount))))
}

func removeImages() {
	output, err := runCommand("docker", "images", "-q")
	if err != nil {
		log.Fatal(string(output))
	}

	imageIDs := strings.Fields(string(output))
	if len(imageIDs) == 0 {
		fmt.Println(colorise("No IMAGES to REMOVE", Warning))
		return
	}

	removedImageCount := 0
	for _, image := range imageIDs {
		output, err = runCommand("docker", "rmi", image)
		if err != nil {
			fmt.Println(colorise(fmt.Sprintf("Failed to remove image %s: %s", image, output), Danger))
		} else {
			removedImageCount++
		}
	}

	if removedImageCount == 0 {
		fmt.Println(colorise("No IMAGES were REMOVED", Danger))
		return
	}

	fmt.Println(colorise(fmt.Sprintf("%s: Removed %d image%s", colorise("REMOVING IMAGES", Success), removedImageCount, pluralise(removedImageCount))))
}

func pruneSystem() {
	output, err := runCommand("docker", "system", "prune", "-f")
	if err != nil {
		log.Fatal(string(output))
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
		fmt.Println(colorise("NOTHING to PRUNE", Warning))
		return
	}

	fmt.Println(fmt.Sprintf("%s: %s", colorise("Pruning SYSTEM", Success), colorise(reclaimedSpace)))
}
