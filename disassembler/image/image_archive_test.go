package image

import (
	"fmt"
	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tklab-group/docker-image-disassembler/disassembler/testutil"
	"testing"
)

func TestNewImageArchive(t *testing.T) {
	buf := testutil.ReadFileForBuffer(t, "testdata/dockerimage-add-file.tar")

	imageArchive, err := NewImageArchive(buf)
	assert.NoError(t, err)

	assert.NotNil(t, imageArchive.Manifest)
	assert.NotNil(t, imageArchive.Config)

	g := goldie.New(t, goldie.WithFixtureDir("testdata/golden"))
	g.AssertJson(t, "manifest", imageArchive.Manifest)
	g.AssertJson(t, "config", imageArchive.Config)

	layerCounts := len(imageArchive.Manifest.LayerTarPaths)
	lastLayer, err := imageArchive.GetFileTreeByLayerIndex(layerCounts - 1)
	assert.NoError(t, err)
	node := lastLayer.FindNodeFromPath("/a/bb/ccc/dddd/eeeee")
	assert.NotNil(t, node)

	g.Assert(t, "added-file", node.Info.Data)
}

func TestImageArchive_GetLatestFileNode(t *testing.T) {
	buf := testutil.ReadFileForBuffer(t, "testdata/dockerimage-add-file.tar")
	imageArchive, err := NewImageArchive(buf)
	require.NoError(t, err)

	assert.NotNil(t, imageArchive.GetLatestFileNode("/a/bb/ccc/dddd/eeeee"))
	assert.Nil(t, imageArchive.GetLatestFileNode("/a/bb/ccc/dddd/eeeee/f"))
}

func TestCollectWhiteoutFiles(t *testing.T) {
	imageTarName, _ := testutil.CreateTarImageFromDockerfile(t, "testdata/Dockerfile.whiteout-file")
	buf := testutil.ReadFileForBuffer(t, imageTarName)
	imageArchive, err := NewImageArchive(buf)
	require.NoError(t, err)

	g := goldie.New(t, goldie.WithFixtureDir("testdata/golden"))
	for i := 0; i < len(imageArchive.Manifest.LayerTarPaths); i++ {
		fileTree, err := imageArchive.GetFileTreeByLayerIndex(i)
		require.NoError(t, err)

		g.AssertJson(t, fmt.Sprintf("whiteout-files-layer%d", i), fileTree.WhiteoutFiles)
	}
}
