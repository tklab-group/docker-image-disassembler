package filetree

import (
	dockerarchive "github.com/docker/docker/pkg/archive"
	"strings"
)

type WhiteoutFile struct {
	Name         string
	OriginalName string
	FileInfo     *FileInfo
	WhiteoutType WhiteoutType
}

type WhiteoutType int

const (
	WhiteoutTypeBasic WhiteoutType = iota
	WhiteoutTypeLinkDir
	WhiteoutTypeOpaqueDir
	WhiteoutTypeOtherMetaPrefix
)

func NewWhiteoutFile(name string, info *FileInfo) *WhiteoutFile {
	whiteoutFile := &WhiteoutFile{
		OriginalName: name,
		FileInfo:     info,
	}

	if !strings.HasPrefix(name, dockerarchive.WhiteoutPrefix) {
		return nil
	}

	if strings.HasPrefix(name, dockerarchive.WhiteoutMetaPrefix) {
		switch name {
		case dockerarchive.WhiteoutLinkDir:
			whiteoutFile.WhiteoutType = WhiteoutTypeLinkDir
		case dockerarchive.WhiteoutOpaqueDir:
			whiteoutFile.WhiteoutType = WhiteoutTypeOpaqueDir
		default:
			whiteoutFile.WhiteoutType = WhiteoutTypeOtherMetaPrefix
		}
		whiteoutFile.Name = name
	} else {
		whiteoutFile.Name = strings.TrimPrefix(name, dockerarchive.WhiteoutPrefix)
		whiteoutFile.WhiteoutType = WhiteoutTypeBasic
	}

	return whiteoutFile
}
