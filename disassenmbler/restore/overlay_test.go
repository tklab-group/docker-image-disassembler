package restore

import (
	"github.com/kr/pretty"
	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/assert"
	"github.com/tklab-group/docker-image-disassembler/disassenmbler/filetree"
	"testing"
)

func TestOverlayFileTree(t *testing.T) {
	base := &filetree.FileTree{
		Root: &filetree.FileNode{
			Children: map[string]*filetree.FileNode{
				"a": {
					Name: "a",
					Info: &filetree.FileInfo{
						Name:  "a",
						Path:  "a",
						IsDir: true,
						Mode:  0444,
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
					Children: map[string]*filetree.FileNode{
						"bb": {
							Name: "bb",
							Info: &filetree.FileInfo{
								Name:  "bb",
								Path:  "b/bb",
								Data:  []byte("bbb"),
								IsDir: false,
							},
							Children: map[string]*filetree.FileNode{},
						},
					},
				},
			},
		},
	}

	ref := &filetree.FileTree{
		Root: &filetree.FileNode{
			Children: map[string]*filetree.FileNode{
				"a": {
					Name: "a",
					Info: &filetree.FileInfo{
						Name:  "a",
						Path:  "a",
						IsDir: true,
						Mode:  0777,
					},
					Children: map[string]*filetree.FileNode{
						"aa": {
							Name: "aa",
							Info: &filetree.FileInfo{
								Name:  "aa",
								Path:  "a/aa",
								Data:  []byte("modified"),
								IsDir: false,
							},
							Children: map[string]*filetree.FileNode{},
						},
						"aa2": {
							Name: "aa2",
							Info: &filetree.FileInfo{
								Name:  "aa2",
								Path:  "a/aa2",
								Data:  []byte("added"),
								IsDir: false,
							},
							Children: map[string]*filetree.FileNode{},
						},
					},
				},
			},
		},
		WhiteoutFiles: []*filetree.WhiteoutFile{
			{
				Name:         "bb",
				OriginalName: ".wh.bb",
				FileInfo: &filetree.FileInfo{
					Name:  ".wh.bb",
					Path:  "b/.wh.bb",
					IsDir: false,
				},
				WhiteoutType: filetree.WhiteoutTypeBasic,
			},
		},
	}

	expected := &filetree.FileTree{
		Root: &filetree.FileNode{
			Children: map[string]*filetree.FileNode{
				"a": {
					Name: "a",
					Info: &filetree.FileInfo{
						Name:  "a",
						Path:  "a",
						IsDir: true,
						Mode:  0777,
					},
					Children: map[string]*filetree.FileNode{
						"aa": {
							Name: "aa",
							Info: &filetree.FileInfo{
								Name:  "aa",
								Path:  "a/aa",
								Data:  []byte("modified"),
								IsDir: false,
							},
							Children: map[string]*filetree.FileNode{},
						},
						"aa2": {
							Name: "aa2",
							Info: &filetree.FileInfo{
								Name:  "aa2",
								Path:  "a/aa2",
								Data:  []byte("added"),
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

	setParentAndTree(base, nil, base.Root)
	setParentAndTree(ref, nil, ref.Root)
	setParentAndTree(expected, nil, expected.Root)

	OverlayFileTree(base, ref)
	assert.Equal(t, expected, base)

	g := goldie.New(t, goldie.WithFixtureDir("testdata/golden"))
	g.Assert(t, "overlay-filetree", []byte(pretty.Sprint(base)))
}

func setParentAndTree(tree *filetree.FileTree, parent *filetree.FileNode, node *filetree.FileNode) {
	node.Tree = tree
	node.Parent = parent

	for _, children := range node.Children {
		setParentAndTree(tree, node, children)
	}
}
