package pkginfo

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestReadAptPkgInfos(t *testing.T) {
	expected := []*AptPkgInfo{
		{Package: "adduser", Version: "3.118ubuntu5"},
		{Package: "apt", Version: "2.2.4ubuntu0.1"},
	}

	b, err := os.ReadFile("testdata/apt_pkg")
	assert.NoError(t, err)
	buf := bytes.NewBuffer(b)

	got, err := ReadAptPkgInfos(buf)
	assert.NoError(t, err)
	assert.Equal(t, expected, got)
}
