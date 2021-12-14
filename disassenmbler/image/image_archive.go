package image

import "github.com/tklab-group/docker-image-disassembler/disassenmbler/filetree"

type ImageArchive struct {
	Manifest *Manifest
	Config   *Config
	LayerMap map[string]*filetree.FileTree
}
