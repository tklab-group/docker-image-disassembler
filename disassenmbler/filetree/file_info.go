package filetree

import "os"

type FileInfo struct {
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
