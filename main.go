package main

import (
	"fmt"
	"github.com/dencoseca/boxi/styles"
	"log"
	"os"
	"strings"
)

// main is the entry point of the application. It processes command-line
// arguments and invokes corresponding handlers. Supported commands include
// container, volume, image operations, and system cleanup actions like wipe and
// purge.
func main() {
	if len(os.Args) < 2 {
		fmt.Println(mainUsage)
		os.Exit(1)
	}

	mainCommand := os.Args[1]

	switch mainCommand {
	case "-h", "help", "--help":
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
	default:
		fmt.Println(mainUsage)
		os.Exit(1)
	}
}

// handleContainers processes container-related commands such as stop, rm, and
// clean, based on user input.
func handleContainers() {
	if len(os.Args) < 3 {
		fmt.Println(containerUsage)
		os.Exit(1)
	}

	subCommand := os.Args[2]

	switch subCommand {
	case "-h", "help", "--help":
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

// handleVolumes processes volume-related commands such as rm, based on user
// input.
func handleVolumes() {
	if len(os.Args) < 3 {
		fmt.Println(volumeUsage)
		os.Exit(1)
	}

	subCommand := os.Args[2]

	switch subCommand {
	case "-h", "help", "--help":
		fmt.Println(volumeUsage)
	case "rm":
		removeVolumes()
	default:
		fmt.Println(volumeUsage)
		os.Exit(1)
	}
}

// handleImages processes image-related commands such as rm, based on user input.
func handleImages() {
	if len(os.Args) < 3 {
		fmt.Println(imageUsage)
		os.Exit(1)
	}

	subCommand := os.Args[2]

	switch subCommand {
	case "-h", "help", "--help":
		fmt.Println(imageUsage)
	case "rm":
		removeImages()
	case "rmf":
		removeImages(true)
	default:
		fmt.Println(imageUsage)
		os.Exit(1)
	}
}

// wipe stops and removes all Docker containers and removes all Docker volumes.
func wipe() {
	stopContainers()
	removeContainers()
	removeVolumes()
}

// purge stops and removes all Docker containers, removes all Docker volumes,
// removes all Docker images, and prunes the Docker system including networks and
// the build cache.
func purge() {
	stopContainers()
	removeContainers()
	removeVolumes()
	removeImages(true)
	pruneSystem()
}

// stopContainers stops all running Docker containers and logs the result for
// each container stopped, or any errors encountered.
func stopContainers() {
	output, err := runCommand("docker", "ps", "--format", "{{.Names}}")
	if err != nil {
		log.Fatal(output)
	}

	containerNames := strings.Fields(output)
	if len(containerNames) == 0 {
		styles.Yellow("No CONTAINERS to STOP")
		return
	}

	stoppedContainerCount := 0
	for _, container := range containerNames {
		output, err = runCommand("docker", "stop", container)
		if err != nil {
			styles.Red("Failed to stop container %s: %s", container, output)
		} else {
			stoppedContainerCount++
		}
	}

	if stoppedContainerCount == 0 {
		styles.Red("No CONTAINERS were STOPPED")
		return
	}

	styles.Green("%s CONTAINER%s STOPPED", stoppedContainerCount, strings.ToUpper(pluralise(stoppedContainerCount)))
}

// removeContainers removes all stopped Docker containers and logs the result for
// each container removed or any errors encountered.
func removeContainers() {
	output, err := runCommand("docker", "ps", "-a", "--format", "{{.Names}}")
	if err != nil {
		log.Fatal(output)
	}

	containerNames := strings.Fields(output)
	if len(containerNames) == 0 {
		styles.Yellow("No CONTAINERS to REMOVE")
		return
	}

	removedContainerCount := 0
	for _, container := range containerNames {
		output, err = runCommand("docker", "rm", container)
		if err != nil {
			styles.Red("Failed to remove container %s: %s", container, output)
		} else {
			removedContainerCount++
		}
	}

	if removedContainerCount == 0 {
		styles.Red("No CONTAINERS were REMOVED")
		return
	}

	styles.Green("%s CONTAINER%s REMOVED", removedContainerCount, strings.ToUpper(pluralise(removedContainerCount)))
}

// removeVolumes removes all dangling Docker volumes and logs the result for each
// volume removed or any errors encountered.
func removeVolumes() {
	output, err := runCommand("docker", "volume", "ls", "-q")
	if err != nil {
		log.Fatal(output)
	}

	volumeIDs := strings.Fields(output)
	if len(volumeIDs) == 0 {
		styles.Yellow("No VOLUMES to REMOVE")
		return
	}

	removedVolumeCount := 0
	for _, volume := range volumeIDs {
		output, err = runCommand("docker", "volume", "rm", volume)
		if err != nil {
			styles.Red("Failed to remove volume %s: %s", volume, output)
		} else {
			removedVolumeCount++
		}
	}

	if removedVolumeCount == 0 {
		styles.Yellow("No VOLUMES were REMOVED")
		return
	}

	styles.Green("%s VOLUME%s REMOVED", removedVolumeCount, strings.ToUpper(pluralise(removedVolumeCount)))
}

// removeImages removes all Docker images, logging the result for each image
// removed or any errors encountered.
func removeImages(force ...bool) {
	forceFlag := false
	if len(force) > 0 {
		forceFlag = true
	}

	output, err := runCommand("docker", "images", "-q")
	if err != nil {
		log.Fatal(output)
	}

	imageIDs := strings.Fields(output)
	if len(imageIDs) == 0 {
		styles.Yellow("No IMAGES to REMOVE")
		return
	}

	removedImageCount := 0
	for _, image := range imageIDs {
		args := []string{"rmi", image}
		if forceFlag {
			args = append(args, "-f")
		}

		output, err = runCommand("docker", args...)
		if err != nil {
			styles.Red("Failed to remove image %s: %s", image, output)
		} else {
			removedImageCount++
		}
	}

	if removedImageCount == 0 {
		styles.Red("No IMAGES were REMOVED")
		return
	}

	styles.Green("%s IMAGE%s REMOVED", removedImageCount, strings.ToUpper(pluralise(removedImageCount)))
}

// pruneSystem prunes the Docker system, logging the total space reclaimed or a
// warning if nothing was pruned.
func pruneSystem() {
	output, err := runCommand("docker", "system", "prune", "-f")
	if err != nil {
		log.Fatal(output)
	}

	lines := strings.Split(output, "\n")
	var reclaimedSpace string

	for _, line := range lines {
		if strings.Contains(line, "Total reclaimed space") {
			reclaimedSpace = line
			break
		}
	}

	if reclaimedSpace == "Total reclaimed space: 0B" {
		styles.Yellow("NOTHING to PRUNE")
		return
	}

	styles.Green("PRUNED %s", reclaimedSpace)
}
