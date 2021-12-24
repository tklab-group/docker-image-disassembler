package image

import (
	"bytes"
	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestNewImageArchive(t *testing.T) {
	tarFile, err := os.ReadFile("testdata/dockerimage-add-file.tar")
	assert.NoError(t, err)
	buf := bytes.NewBuffer(tarFile)

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
