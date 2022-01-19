package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
)

var targetDir string

func init() {
	flag.StringVar(&targetDir, "targetDir", "/", "target directory to exec ls in the container")
}

func main() {
	flag.Parse()
	args := flag.Args()

	err := run(args, os.Stdout)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}

func run(args []string, out io.Writer) error {
	if len(args) != 1 {
		return fmt.Errorf("1 arg is required")
	}

	imageID := args[0]

	cmdOnContainer := fmt.Sprintf("cd %s && find . -type f -ls | awk '{ print $11, $7 }'", targetDir)
	cmd := exec.Command("docker", "run", "--rm", imageID, "/bin/sh", "-c", cmdOnContainer)
	cmd.Stdout = out

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("faild to exec command: %w", err)
	}

	return nil
}
