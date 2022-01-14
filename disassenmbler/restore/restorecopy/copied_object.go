package restorecopy

import (
	dockerimage "github.com/docker/docker/image"
	"github.com/tklab-group/docker-image-disassembler/disassenmbler/filetree"
)

type CopiedObject struct {
	History dockerimage.History
	Object  *filetree.FileNode
}
