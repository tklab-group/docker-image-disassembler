package checkpkg

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_checkPackageInformation(t *testing.T) {
	got, err := checkPackageInformation("testdata/Dockerfile")
	require.NoError(t, err)

	require.Len(t, got, 2)
	assert.Equal(t, got[0].name, "tzdata")
	assert.Equal(t, got[0].versionInDfile, "")
	assert.NotEqual(t, got[0].versionInImage, "")

	assert.Equal(t, got[1].name, "wget")
	assert.Equal(t, got[1].versionInDfile, "1.21-1ubuntu3")
	assert.Equal(t, got[1].versionInImage, "1.21-1ubuntu3")
}

func Test_getAptPkgInfoInImageFromDfile(t *testing.T) {
	got, err := getAptPkgInfoInImageFromDfile("testdata/Dockerfile")
	require.NoError(t, err)
	assert.True(t, len(got) > 2)
	assert.NotEqual(t, got["tzdata"], "")
	assert.Equal(t, got["wget"], "1.21-1ubuntu3")
}

func Test_outPackageVersionDiff(t *testing.T) {
	packageInfos := []packageInfo{
		{
			name:           "tzdata",
			versionInDfile: "",
			versionInImage: "2021e-0ubuntu0.21.04",
		},
		{
			name:           "wget",
			versionInDfile: "1.21-1ubuntu3",
			versionInImage: "1.21-1ubuntu3",
		},
	}

	buf := bytes.Buffer{}
	err := outPackageVersionDiff(&buf, packageInfos)
	require.NoError(t, err)
	assert.Equal(t, "tzdata => tzdata=2021e-0ubuntu0.21.04\n", buf.String())
}
