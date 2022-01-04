package dfile

import (
	"bytes"
	"github.com/kr/pretty"
	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/require"
	"github.com/tklab-group/docker-image-disassembler/disassenmbler/image"
	"github.com/tklab-group/docker-image-disassembler/disassenmbler/pkginfo"
	"github.com/tklab-group/docker-image-disassembler/disassenmbler/testutil"
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	buf := testutil.ReadFileForBuffer(t, "testdata/Dockerfile")
	parsed, err := Parse(buf)
	require.NoError(t, err)

	g := goldie.New(t, goldie.WithFixtureDir("testdata/golden"))
	g.Assert(t, "parsed", []byte(pretty.Sprint(parsed)))
}

func TestNoteLackOfPackageVersion_Apt(t *testing.T) {
	buf := testutil.ReadFileForBuffer(t, "testdata/Dockerfile")
	parsed, err := Parse(buf)
	require.NoError(t, err)

	type PackageInfo struct {
		Package string
		Version string
	}

	// Extract apt packages information from dockerfile.
	dfilePackageInfos := make([]PackageInfo, 0)
	for _, node := range parsed.AST.Children {
		if strings.ToLower(node.Value) != "run" {
			continue
		}
		cmdInfos := NewCmdInfos(node.Next.Value)
		for _, cmdInfo := range cmdInfos {
			if !(cmdInfo.MainCmd == "apt-get" && len(cmdInfo.Args) > 1 && cmdInfo.Args[0] == "install") {
				continue
			}

			for i, arg := range cmdInfo.Args {
				if i < 1 || strings.HasPrefix(arg, "-") {
					continue
				}
				split := strings.Split(arg, "=")
				packageInfo := PackageInfo{Package: split[0]}
				if len(split) == 2 {
					packageInfo.Version = split[1]
				}
				dfilePackageInfos = append(dfilePackageInfos, packageInfo)
			}
		}
	}

	// Extract apt packages information from docker image.
	imageTarName, _ := testutil.CreateTarImageFromDockerfile(t, "testdata/Dockerfile")
	imagTarBuf := testutil.ReadFileForBuffer(t, imageTarName)
	imageArchive, err := image.NewImageArchive(imagTarBuf)
	require.NoError(t, err)

	aptPkgFile := imageArchive.GetLatestFileNode(pkginfo.AptPkgFilePath)
	require.NotNil(t, aptPkgFile)

	buf = bytes.NewBuffer(aptPkgFile.Info.Data)
	aptPkgInfos, err := pkginfo.ReadAptPkgInfos(buf)
	require.NoError(t, err)

	aptPkgInfoMap := map[string]string{}
	for _, aptPkgInfo := range aptPkgInfos {
		aptPkgInfoMap[aptPkgInfo.Package] = aptPkgInfo.Version
	}

	type CorrectPackageInfo struct {
		Package  string
		Actual   string
		Expected string
	}

	// Find diff between information from docker file and from docker image.
	correctPackageInfos := make([]CorrectPackageInfo, 0)
	for _, dfilePackageInfo := range dfilePackageInfos {
		versionInImage, ok := aptPkgInfoMap[dfilePackageInfo.Package]
		require.True(t, ok)
		if dfilePackageInfo.Version != versionInImage {
			correctPackageInfo := CorrectPackageInfo{
				Package:  dfilePackageInfo.Package,
				Actual:   dfilePackageInfo.Version,
				Expected: versionInImage,
			}
			correctPackageInfos = append(correctPackageInfos, correctPackageInfo)
		}
	}

	g := goldie.New(t, goldie.WithFixtureDir("testdata/golden"))
	g.AssertJson(t, "lack-version-info", correctPackageInfos)
}
