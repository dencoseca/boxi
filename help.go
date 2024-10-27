package main

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

const containerUsage = `
Usage: 	boxi [con|container|containers] <command>

Clear up Docker container resources

Commands:
  stop     Stop all running containers
  rm       Remove all stopped containers
  clean    Stop and remove all running containers`

const volumeUsage = `
Usage: 	boxi [vol|volume|volumes] <command>

Clear up Docker volume resources

Commands:
  rm    Remove all dangling volumes`

const imageUsage = `
Usage: 	boxi [img|image|images] <command>

Clear up Docker image resources

Commands:
  rm    Remove all dangling images
  rmf   Force remove all dangling images`
