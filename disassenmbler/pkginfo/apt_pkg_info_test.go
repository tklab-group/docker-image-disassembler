package pkginfo

import (
	"bytes"
	"fmt"
	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tklab-group/docker-image-disassembler/disassenmbler/image"
	"github.com/tklab-group/docker-image-disassembler/disassenmbler/image/docker"
	"html/template"
	"os"
	"strings"
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

// TestAptDockerfileReproduction tests the reproduction of Dockerfile containing the same apt packages.
// TODO: Rename and move.
func TestAptDockerfileReproduction(t *testing.T) {
	baseIid, err := docker.BuildImageFromCli([]string{"-f", " testdata/Dockerfile.apt", "."})
	require.NoError(t, err)

	tmpDir := t.TempDir()
	baseImageTar, err := os.CreateTemp(tmpDir, "dockerimage-*.tar")
	require.NoError(t, err)

	err = docker.RunDockerCmd("save", []string{baseIid, "-o", baseImageTar.Name()}, nil)
	require.NoError(t, err)

	b, err := os.ReadFile(baseImageTar.Name())
	require.NoError(t, err)
	buf := bytes.NewBuffer(b)
	imageArchive, err := image.NewImageArchive(buf)
	require.NoError(t, err)

	lastLayerName := imageArchive.Manifest.LayerTarPaths[len(imageArchive.Manifest.LayerTarPaths)-1]
	lastLayerFileTree := imageArchive.LayerMap[lastLayerName]
	aptPkgFile := lastLayerFileTree.FindNodeFromPath(AptPkgFilePath)
	require.NotNil(t, aptPkgFile)

	buf = bytes.NewBuffer(aptPkgFile.Info.Data)
	aptPkgInfos, err := ReadAptPkgInfos(buf)
	require.NoError(t, err)

	installList := make([]string, 0)
	for _, aptPkgInfo := range aptPkgInfos {
		s := fmt.Sprintf("%s=%s", aptPkgInfo.Package, aptPkgInfo.Version)
		installList = append(installList, s)
	}

	dockerfileTemplate := `FROM {{ .baseImage }}
RUN apt-get update \
	&& apt-get install -y --allow-downgrades \
	` + strings.Join(installList, " \\\n\t")

	tpl, err := template.New("").Parse(dockerfileTemplate)
	require.NoError(t, err)
	m := map[string]interface{}{
		"baseImage": "ubuntu:hirsute-20211107",
	}
	dockerfileBuf := bytes.Buffer{}
	err = tpl.Execute(&dockerfileBuf, m)
	require.NoError(t, err)

	g := goldie.New(t, goldie.WithFixtureDir("testdata/golden"))
	g.Assert(t, "apt.dockerfile", dockerfileBuf.Bytes())

	generatedIid, err := docker.BuildImageFromCli([]string{"-f", " testdata/golden/apt.dockerfile.golden", "."})
	require.NoError(t, err)

	fromBase := checkInstalledAptPackages(t, baseIid)
	fromGenerated := checkInstalledAptPackages(t, generatedIid)
	assert.Equal(t, fromBase, fromGenerated)
}

func checkInstalledAptPackages(t *testing.T, iid string) string {
	t.Helper()
	outBuf := bytes.Buffer{}
	stds := &docker.RunDockerCmdStds{Stdout: &outBuf}
	err := docker.RunDockerCmd("run", []string{"--rm", iid, "dpkg", "-l"}, stds)
	require.NoError(t, err)
	return outBuf.String()
}
