package filetree

import (
	"archive/tar"
	"io"
	"os"
	"path/filepath"
)

// FileInfo contains tar metadata for a specific FileNode.
type FileInfo struct {
	Name     string
	Path     string
	TypeFlag byte
	Linkname string
	Data     []byte // Data contains actual data in the file.
	Size     int64
	Mode     os.FileMode
	Uid      int
	Gid      int
	IsDir    bool
}

// NewFileInfoFromTarHeader extracts the metadata from a tar header and file contents.
func NewFileInfoFromTarHeader(reader *tar.Reader, header *tar.Header) (*FileInfo, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	return &FileInfo{
		Name:     filepath.Base(header.Name),
		Path:     filepath.Clean(header.Name),
		TypeFlag: header.Typeflag,
		Linkname: header.Linkname,
		Data:     data,
		Size:     header.FileInfo().Size(),
		Mode:     header.FileInfo().Mode(),
		Uid:      header.Uid,
		Gid:      header.Gid,
		IsDir:    header.FileInfo().IsDir(),
	}, nil
}

// Copy duplicates FileInfo.
func (info *FileInfo) Copy() *FileInfo {
	var fileInfo FileInfo
	fileInfo = *info
	return &fileInfo
}
