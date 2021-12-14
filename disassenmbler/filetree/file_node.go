package filetree

// FileNode represents a single file or directory.
type FileNode struct {
	Tree     *FileTree
	Parent   *FileNode
	Name     string
	Info     *FileInfo
	Children map[string]*FileNode
}

// NewFileNode creates a new FileNode relative to the given parent node with a payload.
func NewFileNode(parent *FileNode, info *FileInfo) *FileNode {
	node := &FileNode{
		Parent:   parent,
		Name:     info.Name,
		Info:     info,
		Children: map[string]*FileNode{},
	}

	if parent != nil {
		node.Tree = parent.Tree
	}

	return node
}
