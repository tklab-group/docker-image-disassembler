package image

import (
	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/assert"
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
	lastLayer, ok := imageArchive.LayerMap[imageArchive.Manifest.LayerTarPaths[layerCounts-1]]
	assert.True(t, ok)
	node := lastLayer.FindNodeFromPath("/a/bb/ccc/dddd/eeeee")
	assert.NotNil(t, node)

	g.Assert(t, "added-file", node.Info.Data)
}
