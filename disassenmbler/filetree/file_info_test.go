package filetree

import (
	"archive/tar"
	"bytes"
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"testing"
)

func TestNewFileInfoFromTarHeader(t *testing.T) {
	expected := []*FileInfo{
		{
			Name:     "directory",
			Path:     "directory",
			TypeFlag: 53,
			Linkname: "",
			Data:     []byte{},
			Size:     0,
			Mode:     2147484141,
			Uid:      501,
			Gid:      20,
			IsDir:    true,
		},
		{
			Name:     "file.txt",
			Path:     "directory/file.txt",
			TypeFlag: 48,
			Linkname: "",
			Data:     []byte("aaa\nbb\nc\n"),
			Size:     9,
			Mode:     420,
			Uid:      501,
			Gid:      20,
			IsDir:    false,
		},
	}

	b, err := ioutil.ReadFile("testdata/new-file-info-from-tar-header.tar")
	buf := bytes.NewBuffer(b)
	assert.NoError(t, err)

	tarReader := tar.NewReader(buf)
	var index int
	for {
		header, err := tarReader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			assert.NoError(t, err)
		}

		if index >= len(expected) {
			t.Fatalf("index %d is out of range", index)
		}

		got, err := NewFileInfoFromTarHeader(tarReader, header)
		assert.NoError(t, err)
		assert.Equal(t, expected[index], got)

		index++
	}
}
