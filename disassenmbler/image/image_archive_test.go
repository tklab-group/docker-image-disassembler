package image

import (
	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tklab-group/docker-image-disassembler/disassenmbler/testutil"
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
