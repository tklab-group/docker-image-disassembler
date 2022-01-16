package testutil

import (
	"github.com/stretchr/testify/require"
	"github.com/tklab-group/docker-image-disassembler/disassembler/image/docker"
	"os"
	"testing"
)

// CreateTarImageFromDockerfile builds docker image from dockerfile and exports the image as a tar file.
// The tar file is automatically removed when the test complete.
func CreateTarImageFromDockerfile(t *testing.T, path string) (imageTarName string, imageID string) {
	tmpDir := t.TempDir()
	tarFile, err := os.CreateTemp(tmpDir, "dockerimage-*.tar")
	require.NoError(t, err)

	iid, err := docker.CreateTarImageFromDockerfile(path, tarFile.Name())
	require.NoError(t, err)

	return tarFile.Name(), iid
}
