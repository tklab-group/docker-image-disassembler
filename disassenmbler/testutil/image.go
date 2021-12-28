package testutil

import (
	"github.com/stretchr/testify/require"
	"github.com/tklab-group/docker-image-disassembler/disassenmbler/image/docker"
	"os"
	"testing"
)

// CreateTarImageFromDockerfile builds docker image from dockerfile and exports the image as a tar file.
// The tar file is automatically removed when the test complete.
func CreateTarImageFromDockerfile(t *testing.T, path string) string {
	iid, err := docker.BuildImageFromCli([]string{"-f", path, "."})
	require.NoError(t, err)

	tmpDir := t.TempDir()
	imageTar, err := os.CreateTemp(tmpDir, "dockerimage-*.tar")
	require.NoError(t, err)

	err = docker.RunDockerCmd("save", []string{iid, "-o", imageTar.Name()}, nil)
	require.NoError(t, err)

	return imageTar.Name()
}