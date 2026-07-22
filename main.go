package main

import (
	"os"

	"github.com/davidborzek/docker-exporter/cmd"
)

var version = "dev"

func main() {
	cmd.Main(version, os.Args)
}
