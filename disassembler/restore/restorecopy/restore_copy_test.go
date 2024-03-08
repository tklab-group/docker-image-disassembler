package restorecopy

import (
	dockerimage "github.com/docker/docker/image"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tklab-group/docker-image-disassembler/disassembler/filetree"
	"github.com/tklab-group/docker-image-disassembler/disassembler/image"
	"github.com/tklab-group/docker-image-disassembler/disassembler/image/docker"
	"github.com/tklab-group/docker-image-disassembler/disassembler/testutil"
	"os"
	"testing"
)

func TestRestoreCopiedObjects(t *testing.T) {
	imageTarName, _ := testutil.CreateTarImageFromDockerfile(t, "testdata/Dockerfile.copy")
	buf := testutil.ReadFileForBuffer(t, imageTarName)
	imageArchive, err := image.NewImageArchive(buf)
	require.NoError(t, err)
	copiedObjects, err := RestoreCopiedObjects(imageArchive)
	require.NoError(t, err)

	want := []*CopiedObject{
		{
			History: dockerimage.History{
				CreatedBy:  "COPY testdata/src/a a-copied # buildkit",
				EmptyLayer: false},
			Object: &filetree.FileNode{
				Name: "a-copied",
				Info: &filetree.FileInfo{
					Name:     "a-copied",
					Path:     "copied/a-file/a-copied",
					TypeFlag: 0x30,
					Linkname: "",
					Data:     []byte("aaa"),
					Size:     3,
					Mode:     0x1a4,
					Uid:      0,
					Gid:      0,
					IsDir:    false,
				},
				Children: map[string]*filetree.FileNode{},
			},
		},
		{
			History: dockerimage.History{
				CreatedBy:  "COPY testdata/src/b /copied/b-file # buildkit",
				EmptyLayer: false,
			},
			Object: &filetree.FileNode{
				Name: "b-file",
				Info: &filetree.FileInfo{
					Name:     "b-file",
					Path:     "copied/b-file",
					TypeFlag: 0x35,
					Linkname: "",
					Data:     []byte(""),
					Size:     0,
					Mode:     0x800001ed,
					Uid:      0,
					Gid:      0,
					IsDir:    true,
				},
				Children: map[string]*filetree.FileNode{
					"bb": {
						Name: "bb",
						Info: &filetree.FileInfo{
							Name:     "bb",
							Path:     "copied/b-file/bb",
							TypeFlag: 0x30,
							Linkname: "",
							Data:     []byte("bbb"),
							Size:     3,
							Mode:     0x1a4,
							Uid:      0,
							Gid:      0,
							IsDir:    false,
						},
						Children: map[string]*filetree.FileNode{},
					},
				},
			},
		},
		{
			History: dockerimage.History{
				CreatedBy:  "COPY testdata/src/*/cc c-file/ # buildkit",
				EmptyLayer: false,
			},
			Object: &filetree.FileNode{
				Name: "c-file",
				Info: &filetree.FileInfo{
					Name:     "c-file",
					Path:     "copied/c-file",
					TypeFlag: 0x35,
					Linkname: "",
					Data:     []byte(""),
					Size:     0,
					Mode:     0x800001ed,
					Uid:      0,
					Gid:      0,
					IsDir:    true,
				},
				Children: map[string]*filetree.FileNode{
					"cc": {
						Name: "cc",
						Info: &filetree.FileInfo{
							Name:     "cc",
							Path:     "copied/c-file/cc",
							TypeFlag: 0x30,
							Linkname: "",
							Data:     []byte("ccc"),
							Size:     3,
							Mode:     0x1a4,
							Uid:      0,
							Gid:      0,
							IsDir:    false,
						},
						Children: map[string]*filetree.FileNode{},
					},
				},
			},
		},
	}

	opts := []cmp.Option{
		cmpopts.IgnoreFields(CopiedObject{}, "LayerID"),
		cmpopts.IgnoreFields(dockerimage.History{}, "Created", "Author", "Comment"),
		cmpopts.IgnoreFields(filetree.FileNode{}, "Tree", "Parent"),
	}

	if diff := cmp.Diff(want, copiedObjects, opts...); diff != "" {
		t.Errorf("copiedObjects mismatch (-want +got):\n%s", diff)
	}

	for _, copiedObject := range copiedObjects {
		assert.NotEqual(t, "", copiedObject.LayerID)
	}
}

func TestRestoreCopiedObjects2(t *testing.T) {
	iid := "hello-world"
	err := docker.RunDockerCmd("pull", []string{iid}, nil)
	require.NoError(t, err)

	tmpDir := t.TempDir()
	tarFile, err := os.CreateTemp(tmpDir, "dockerimage-*.tar")
	require.NoError(t, err)

	err = docker.RunDockerCmd("save", []string{iid, "-o", tarFile.Name()}, nil)

	buf := testutil.ReadFileForBuffer(t, tarFile.Name())
	imageArchive, err := image.NewImageArchive(buf)
	require.NoError(t, err)

	copiedObjects, err := RestoreCopiedObjects(imageArchive)
	require.NoError(t, err)

	assert.Len(t, copiedObjects, 1)
}
