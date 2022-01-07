package restore

import (
	"bytes"
	"fmt"
	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tklab-group/docker-image-disassembler/disassenmbler/filetree"
	"github.com/tklab-group/docker-image-disassembler/disassenmbler/image"
	"github.com/tklab-group/docker-image-disassembler/disassenmbler/image/docker"
	"github.com/tklab-group/docker-image-disassembler/disassenmbler/testutil"
	"os/exec"
	"regexp"
	"testing"
)

func TestRestoreLastLayer(t *testing.T) {
	imageTarName, iid := testutil.CreateTarImageFromDockerfile(t, "testdata/Dockerfile.restore-layer")
	imageArchive, err := image.NewImageArchive(testutil.ReadFileForBuffer(t, imageTarName))
	require.NoError(t, err)

	base := filetree.NewFileTree()
	for i := 0; i < len(imageArchive.Manifest.LayerTarPaths); i++ {
		layerTree, err := imageArchive.GetFileTreeByLayerIndex(i)
		require.NoError(t, err)
		OverlayFileTree(base, layerTree)
	}

	targetPath := t.TempDir()
	err = EmbodyFileTree(targetPath, base)
	require.NoError(t, err)

	output, err := exec.Command("/bin/sh", "-c", fmt.Sprintf("cd %s/a && tree", targetPath)).Output()
	require.NoError(t, err)

	g := goldie.New(t, goldie.WithFixtureDir("testdata/golden"))
	g.Assert(t, "restore-layer-exec-tree", output)

	outBufFromImage := bytes.Buffer{}
	stds := &docker.RunDockerCmdStds{Stdout: &outBufFromImage}
	err = docker.RunDockerCmd("run", []string{"--rm", iid, "tree", "var", "-i"}, stds)
	require.NoError(t, err)

	output, err = exec.Command("/bin/sh", "-c", fmt.Sprintf("cd %s && tree var -i", targetPath)).Output()
	require.NoError(t, err)

	assertTreeOutput(t, outBufFromImage.String(), string(output))
}

// assertTreeOutput ignore description related to symbolic link.
func assertTreeOutput(t *testing.T, expected string, actual string) bool {
	t.Helper()

	// Remove description for symbolic link.
	symlinkReg := regexp.MustCompile(`\s->\s.+`)
	expected = symlinkReg.ReplaceAllString(expected, "")
	actual = symlinkReg.ReplaceAllString(actual, "")

	// Trim description of files and dirs counts.
	countReg := regexp.MustCompile(`\d+\sdirectories,\s\d+\sfiles`)
	expected = countReg.ReplaceAllString(expected, "")
	actual = countReg.ReplaceAllString(actual, "")

	return assert.Equal(t, expected, actual)
}
