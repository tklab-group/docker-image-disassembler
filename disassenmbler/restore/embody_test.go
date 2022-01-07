package restore

import (
	"fmt"
	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tklab-group/docker-image-disassembler/disassenmbler/filetree"
	"github.com/tklab-group/docker-image-disassembler/disassenmbler/testutil"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestEmbodyFileTree(t *testing.T) {
	targetPath := t.TempDir()

	fileTree := &filetree.FileTree{
		Root: &filetree.FileNode{
			Name: ".",
			Children: map[string]*filetree.FileNode{
				"a": {
					Name: "a",
					Info: &filetree.FileInfo{
						Name:  "a",
						Path:  "a",
						IsDir: true,
					},
					Children: map[string]*filetree.FileNode{
						"aa": {
							Name: "aa",
							Info: &filetree.FileInfo{
								Name:  "aa",
								Path:  "a/aa",
								Data:  []byte("aaa"),
								IsDir: false,
							},
							Children: map[string]*filetree.FileNode{},
						},
					},
				},
				"b": {
					Name: "b",
					Info: &filetree.FileInfo{
						Name:  "b",
						Path:  "b",
						IsDir: true,
					},
					Children: map[string]*filetree.FileNode{},
				},
			},
		},
	}

	err := EmbodyFileTree(targetPath, fileTree)
	require.NoError(t, err)

	output, err := exec.Command("/bin/sh", "-c", fmt.Sprintf("cd %s && tree", targetPath)).Output()
	require.NoError(t, err)

	g := goldie.New(t, goldie.WithFixtureDir("testdata/golden"))
	g.Assert(t, "embody-filetree-exec-tree", output)

	buf := testutil.ReadFileForBuffer(t, filepath.Join(targetPath, "a/aa"))
	assert.Equal(t, "aaa", buf.String())
}
