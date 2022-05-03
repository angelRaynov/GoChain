package main

import (
	"github.com/angelRaynov/GoChain/cli"
	"os"
)

func main() {
	defer os.Exit(0)

	cmd := cli.CommandLine{}

	cmd.Run()
}
