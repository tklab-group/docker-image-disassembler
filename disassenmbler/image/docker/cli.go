package docker

import (
	"fmt"
	"github.com/tklab-group/docker-image-disassembler/disassenmbler/utils"
	"io"
	"os"
	"os/exec"
)

type RunDockerCmdStds struct {
	Stdout io.Writer
	Stderr io.Writer
	Stdin  io.Reader
}

var defaultRunDockerCmdStds = RunDockerCmdStds{
	Stdout: os.Stdout,
	Stderr: os.Stderr,
	Stdin:  os.Stdin,
}

// RunDockerCmd runs a given Docker command.
// args and stds are optional.
func RunDockerCmd(cmdStr string, args []string, stds *RunDockerCmdStds) error {
	if !isDockerClientBinaryAvailable() {
		return fmt.Errorf("cannot find docker client executable")
	}

	allArgs := utils.CleanArgs(append([]string{cmdStr}, args...))

	cmd := exec.Command("docker", allArgs...)
	cmd.Env = os.Environ()

	if stds == nil {
		stds = &defaultRunDockerCmdStds
	}
	cmd.Stdout = stds.Stdout
	cmd.Stderr = stds.Stderr
	cmd.Stdin = stds.Stdin

	return cmd.Run()
}

func isDockerClientBinaryAvailable() bool {
	_, err := exec.LookPath("docker")
	return err == nil
}
