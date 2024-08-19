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
		log.Fatal(err)
	}

	containerNames := strings.Fields(string(output))
	if len(containerNames) == 0 {
		fmt.Println(colorise("No CONTAINERS to STOP", Danger))
		return
	}

	stoppedContainerCount := 0
	for _, container := range containerNames {
		_, err = runCommand("docker", "stop", container)
		if err != nil {
			fmt.Println(colorise(fmt.Sprintf("Failed to stop container %s: %v", container, err), Danger))
		} else {
			stoppedContainerCount++
		}
	}

	fmt.Println(colorise(fmt.Sprintf("%s: Stopped %d container%s", colorise("STOPPING CONTAINERS", Success), stoppedContainerCount, pluralise(stoppedContainerCount))))
}

func removeContainers() {
	output, err := runCommand("docker", "ps", "-a", "--format", "{{.Names}}")
	if err != nil {
		log.Fatal(err)
	}

	containerNames := strings.Fields(string(output))
	if len(containerNames) == 0 {
		fmt.Println(colorise("No CONTAINERS to REMOVE", Danger))
		return
	}

	removedContainerCount := 0
	for _, container := range containerNames {
		_, err = runCommand("docker", "rm", container)
		if err != nil {
			fmt.Println(colorise(fmt.Sprintf("Failed to remove container %s: %v", container, err), Danger))
		} else {
			removedContainerCount++
		}
	}

	fmt.Println(colorise(fmt.Sprintf("%s: Removed %d container%s", colorise("REMOVING CONTAINERS", Success), removedContainerCount, pluralise(removedContainerCount))))
}

func removeVolumes() {
	output, err := runCommand("docker", "volume", "ls", "-q")
	if err != nil {
		log.Fatal(err)
	}

	volumeIDs := strings.Fields(string(output))
	if len(volumeIDs) == 0 {
		fmt.Println(colorise("No VOLUMES to REMOVE", Danger))
		return
	}

	removedVolumeCount := 0
	for _, volume := range volumeIDs {
		_, err = runCommand("docker", "volume", "rm", volume)
		if err != nil {
			fmt.Println(colorise(fmt.Sprintf("Failed to remove volume %s: %v", volume, err), Danger))
		} else {
			removedVolumeCount++
		}
	}

	fmt.Println(colorise(fmt.Sprintf("%s: Removed %d volume%s", colorise("REMOVING VOLUMES", Success), removedVolumeCount, pluralise(removedVolumeCount))))
}

func removeImages() {
	output, err := runCommand("docker", "images", "-q")
	if err != nil {
		log.Fatal(err)
	}

	imageIDs := strings.Fields(string(output))
	if len(imageIDs) == 0 {
		fmt.Println(colorise("No IMAGES to REMOVE", Danger))
		return
	}

	removedImageCount := 0
	for _, image := range imageIDs {
		_, err = runCommand("docker", "rmi", image)
		if err != nil {
			fmt.Println(colorise(fmt.Sprintf("Failed to remove image %s: %v", image, err), Danger))
		} else {
			removedImageCount++
		}
	}

	fmt.Println(colorise(fmt.Sprintf("%s: Removed %d image%s", colorise("REMOVING IMAGES", Success), removedImageCount, pluralise(removedImageCount))))
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
		fmt.Println(colorise("NOTHING to PRUNE", Danger))
		return
	}

	fmt.Println(fmt.Sprintf("%s: %s", colorise("Pruning SYSTEM", Success), colorise(reclaimedSpace)))
}
