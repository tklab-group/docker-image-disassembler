package restore

import (
	"github.com/tklab-group/docker-image-disassembler/disassenmbler/filetree"
	"path/filepath"
)

// OverlayFileTree overlays ref over base.
func OverlayFileTree(base *filetree.FileTree, ref *filetree.FileTree) {
	overlayFileNode(base.Root, ref.Root)

	for _, whiteoutFile := range ref.WhiteoutFiles {
		if whiteoutFile.WhiteoutType != filetree.WhiteoutTypeBasic {
			continue
		}

		parentPath := filepath.Dir(whiteoutFile.FileInfo.Path)
		parentNode := base.FindNodeFromPath(parentPath)
		delete(parentNode.Children, whiteoutFile.Name)
	}
}

func overlayFileNode(base *filetree.FileNode, ref *filetree.FileNode) {
	if base == nil || ref == nil {
		return
	}

	if ref.Info != nil {
		base.Info = ref.Info.Copy()
	}

	for refChildName, refChild := range ref.Children {
		baseChild, ok := base.Children[refChildName]
		if ok {
			overlayFileNode(baseChild, refChild)
		} else {
			base.Children[refChildName] = refChild.Copy(base)
		}
	}
}
