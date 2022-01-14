package main

import (
	"github.com/tklab-group/docker-image-disassembler/cli"
	"os"
)

func main() {
	cli.Execute(cli.Config{
		In:  os.Stdin,
		Out: os.Stdout,
		Err: os.Stdout,
	})
}
