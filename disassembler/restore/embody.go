package restore

import (
	"fmt"
	"github.com/tklab-group/docker-image-disassembler/disassembler/filetree"
	"os"
	"path/filepath"
)

// EmbodyFileTree embodies FileTree under targetPath.
func EmbodyFileTree(targetPath string, fileTree *filetree.FileTree) error {
	if !existDir(targetPath) {
		err := os.MkdirAll(targetPath, 0777)
		if err != nil {
			return fmt.Errorf("faild to create directory by target path: %w", err)
		}
	}

	for _, fileNode := range fileTree.Root.Children {
		err := EmbodyFileNode(targetPath, fileNode)
		if err != nil {
			return err
		}
	}

	return nil
}

// EmbodyFileNode embodies FileNode and its children nodes.
func EmbodyFileNode(parentPath string, fileNode *filetree.FileNode) error {
	if fileNode == nil {
		return nil
	}

	tmpPath := filepath.Join(parentPath, fileNode.Name)

	// Embody the node.
	// Skip if fileNode is dummy node.
	if fileNode.Name != "" {
		if fileNode.Info.IsDir {
			// Ignore its file mode in container.
			err := os.Mkdir(tmpPath, 0777)
			if err != nil {
				return fmt.Errorf("faild to embody FileNode as directory: %w", err)
			}
		} else {
			// Ignore its file mode in container.
			err := os.WriteFile(tmpPath, fileNode.Info.Data, 0777)
			if err != nil {
				return fmt.Errorf("faild to embody FileNode as file: %w", err)
			}
		}
	}

	for childName, child := range fileNode.Children {
		err := EmbodyFileNode(tmpPath, child)
		if err != nil {
			return fmt.Errorf("faild to embody %s: %w", filepath.Join(tmpPath, childName), err)
		}
	}

	return nil
}

func existDir(dirPath string) bool {
	f, err := os.Stat(dirPath)
	if err == nil && f.IsDir() {
		return true
	}

	return false
}
