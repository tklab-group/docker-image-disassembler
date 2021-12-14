package filetree

type FileNode struct {
	Tree     *FileTree
	Parent   *FileNode
	Name     string
	Info     *FileInfo
	Children map[string]*FileNode
	Path     string
}
