package filetree

import (
	"fmt"
	dockerarchive "github.com/docker/docker/pkg/archive"
	"strings"
)

type FileTree struct {
	Root *FileNode
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
