package main

import (
	"github.com/tklab-group/docker-image-disassembler/cli"
	"github.com/tklab-group/docker-image-disassembler/cli/config"
	"os"
)

func main() {
	cli.Execute(config.Config{
		In:  os.Stdin,
		Out: os.Stdout,
		Err: os.Stdout,
	})
}
