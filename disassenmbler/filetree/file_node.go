package filetree

import (
	dockerarchive "github.com/docker/docker/pkg/archive"
	"strings"
)

// FileNode represents a single file or directory.
type FileNode struct {
	Tree     *FileTree
	Parent   *FileNode
	Name     string
	Info     *FileInfo
	Children map[string]*FileNode
}

// NewFileNode creates a new FileNode relative to the given parent node with a payload.
func NewFileNode(parent *FileNode, name string, info *FileInfo) *FileNode {
	node := &FileNode{
		Parent:   parent,
		Name:     name,
		Info:     info,
		Children: map[string]*FileNode{},
	}

	if parent != nil {
		node.Tree = parent.Tree
	}

	return node
}

// AddChild creates a new node relative to the current FileNode.
func (node *FileNode) AddChild(name string, info *FileInfo) *FileNode {
	// Ignore the file has WhiteoutPrefix here since the file isn't a usual file now.
	if strings.HasPrefix(name, dockerarchive.WhiteoutPrefix) {
		return nil
	}

	child, ok := node.Children[name]
	if ok {
		if child.Info == nil {
			child.Info = info
		}
	} else {
		child = NewFileNode(node, name, info)
		node.Children[name] = child
	}

	return child
}

// Copy duplicates FileNode with new parent.
func (node *FileNode) Copy(parent *FileNode) *FileNode {
	newNode := NewFileNode(parent, node.Name, node.Info.Copy())
	newNode.Tree = parent.Tree
	for name, child := range node.Children {
		newNode.Children[name] = child.Copy(newNode)
	}

	return newNode
}
