package disassembler

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_getAptPkgInfoInImageFromDfile(t *testing.T) {
	got, err := GetAptPkgInfoInImageFromDfile("testdata/Dockerfile")
	require.NoError(t, err)
	assert.True(t, len(got) > 2)
	assert.NotEqual(t, got["tzdata"], "")
	assert.Equal(t, got["wget"], "1.21.3-1ubuntu1")
}
