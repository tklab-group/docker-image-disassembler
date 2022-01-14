package restorecopy

import (
	dockerimage "github.com/docker/docker/image"
	"github.com/tklab-group/docker-image-disassembler/disassembler/filetree"
)

type CopiedObject struct {
	History dockerimage.History
	Object  *filetree.FileNode
}
