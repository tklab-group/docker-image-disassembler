package filetree

import (
	"fmt"
	dockerarchive "github.com/docker/docker/pkg/archive"
	"path"
	"strings"
)

type FileTree struct {
	Root          *FileNode
	LayerName     string
	WhiteoutFiles []*WhiteoutFile
}

// NewFileTree creates an empty FileTree.
func NewFileTree() *FileTree {
	tree := &FileTree{
		WhiteoutFiles: make([]*WhiteoutFile, 0),
	}

	// Root is a dummy node.
	root := &FileNode{
		Tree:     tree,
		Parent:   nil,
		Name:     "",
		Info:     nil,
		Children: map[string]*FileNode{},
	}
	tree.Root = root
	return tree
}

// AddNode adds new node to the tree.
func (tree *FileTree) AddNode(info *FileInfo) error {
	if info.Path == "." {
		return fmt.Errorf("cannot ad relative path '%s'", info.Path)
	}
	nodeNames := strings.Split(strings.Trim(info.Path, "/"), "/")
	node := tree.Root
	for i, name := range nodeNames {
		if n, ok := node.Children[name]; ok {
			node = n
		} else {
			// Not to add node under the whiteout node.
			if strings.HasPrefix(name, dockerarchive.WhiteoutPrefix) {
				whiteoutFile := NewWhiteoutFile(name, info)
				if whiteoutFile == nil {
					return fmt.Errorf("could not add whiteout file: '%s' (path: '%s')", name, info.Path)
				}

				tree.WhiteoutFiles = append(tree.WhiteoutFiles, whiteoutFile)
				return nil
			}

			// Just adding intermediary node.
			node = node.AddChild(name, nil)
			if node == nil {
				return fmt.Errorf("could not add child node: '%s' (path: '%s')", name, info.Path)
			}
		}

		// The last node is targeted for the info.
		if i == len(nodeNames)-1 {
			node.Info = info
		}
	}

	return nil
}

// FindNodeFromPath returns FileNode specified by the path.
// If not exist, it returns nil.
func (tree *FileTree) FindNodeFromPath(pathStr string) *FileNode {
	nodeNames := strings.Split(strings.Trim(path.Clean(pathStr), "/"), "/")
	node := tree.Root
	for _, name := range nodeNames {
		n, exist := node.Children[name]
		if !exist {
			return nil
		}
		node = n
	}

	return node
}
