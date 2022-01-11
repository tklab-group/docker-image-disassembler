package restorecopy

import "github.com/tklab-group/docker-image-disassembler/disassenmbler/image"

// RestoreCopiedObjects restores copied objects from ImageArchive.
func RestoreCopiedObjects(imageArchive *image.ImageArchive) ([]*CopiedObject, error) {
	r := newRestorer()
	return r.restore(imageArchive)
}
