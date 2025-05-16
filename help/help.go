package help

import (
	"fmt"
	"os"
)

var MainUsage = `
Usage: 	boxi <command> [subCommand]

Clear up Docker resources

Commands:
  con, container, containers    Container commands
  vol, volume, volumes          Volume commands
  img, image, images            Image commands
  info                          Display Docker system information
  wipe                          Clean up containers and volumes
  purge                         Clean up containers, volumes, images, networks and the build cache

Run 'boxi <command> --help' for more information.`

var ContainerUsage = `
Usage: 	boxi [con|container|containers] <command>

Clear up Docker container resources

Commands:
  stop     Stop all running containers
  rm       Remove all stopped containers
  clean    Stop and remove all running containers`

var VolumeUsage = `
Usage: 	boxi [vol|volume|volumes] <command>

Clear up Docker volume resources

Commands:
  rm    Remove all dangling volumes`

var ImageUsage = `
Usage: 	boxi [img|image|images] <command>

Clear up Docker image resources

Commands:
  rm    Remove all dangling images
  rmf   Force remove all dangling images`

type Usage int

const (
	Main Usage = iota
	Container
	Volume
	Image
)

// PrintHelpAndExit prints the usage instructions for the specified command and
// exits the program with the given exit code.
func PrintHelpAndExit(usage Usage, code ...int) {
	exitCode := 0
	if len(code) > 0 {
		exitCode = code[0]
	}

	output := ""
	switch usage {
	case Main:
		output = MainUsage
	case Container:
		output = ContainerUsage
	case Volume:
		output = VolumeUsage
	case Image:
		output = ImageUsage
	}

	fmt.Println(output)
	os.Exit(exitCode)
}
