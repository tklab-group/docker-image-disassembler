package image

import (
	"fmt"
	dockerimage "github.com/docker/docker/image"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
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

func TestImageArchive_GetHistoryToLayers(t *testing.T) {
	buf := testutil.ReadFileForBuffer(t, "testdata/dockerimage-add-file.tar")
	imageArchive, err := NewImageArchive(buf)
	require.NoError(t, err)

	historyToLayers, err := imageArchive.GetHistoryToLayers()
	require.NoError(t, err)

	want := []*HistoryToLayer{
		{
			History: dockerimage.History{
				CreatedBy:  "/bin/sh -c #(nop) ADD file:9233f6f2237d79659a9521f7e390df217cec49f1a8aa3a12147bbca1956acdb9 in / ",
				Comment:    "",
				EmptyLayer: false,
			},
			LayerID: "1e33021cb9e4fc8e989445d68584dd2386b15e902fc4e95fba64a5a917759466",
		},
		{
			History: dockerimage.History{
				CreatedBy:  "RUN /bin/sh -c mkdir -p /a/bb/ccc/dddd # buildkit",
				Comment:    "buildkit.dockerfile.v0",
				EmptyLayer: false,
			},
			LayerID: "4a0c493ec850512a4f72774b55b08d61fe6a0ab80cda2165cab36f1e20b1429c",
		},
		{
			History: dockerimage.History{
				CreatedBy:  "RUN /bin/sh -c echo eeeee \u003e /a/bb/ccc/dddd/eeeee # buildkit",
				Comment:    "buildkit.dockerfile.v0",
				EmptyLayer: false,
			},
			LayerID: "9027e3575c9e8b582f792d5cd342c3cb0144a170413c98fd092432f782a6d952",
		},
	}

	opts := []cmp.Option{
		cmpopts.IgnoreFields(dockerimage.History{}, "Created"),
		cmpopts.IgnoreFields(HistoryToLayer{}, "Layer"),
	}

	if diff := cmp.Diff(want, historyToLayers, opts...); diff != "" {
		t.Errorf("historyToLayers mismatch (-want +got):\n%s", diff)
	}

	for _, historyToLayer := range historyToLayers {
		assert.NotNil(t, historyToLayer.Layer)
	}
}
