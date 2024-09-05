# Boxi

## Overview

Boxi is a command line tool to provide an easy way to remove docker resources, by resource type and on mass, without
having
to list or specify named resources.

## Getting started

### Prerequisites

Assuming a Mac OS environment...

Go SDK installed:

```shell
brew install go
```

GOPATH/bin folder in your system PATH variable:

```shell
export PATH="$PATH:$(go env GOPATH)/bin"
```

### Installation

```shell
go install github.com/dencoseca/boxi@latest
```

### Usage

```shell
boxi --help
```
