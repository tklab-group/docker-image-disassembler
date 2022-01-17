package checkpkg

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tklab-group/docker-image-disassembler/cli/cmdname"
	"github.com/tklab-group/docker-image-disassembler/cli/config"
	"github.com/tklab-group/docker-image-disassembler/disassembler/dfile"
	"github.com/tklab-group/docker-image-disassembler/disassembler/image"
	"github.com/tklab-group/docker-image-disassembler/disassembler/image/docker"
	"github.com/tklab-group/docker-image-disassembler/disassembler/pkginfo"
	"io"
	"os"
	"strings"
)

func Cmd(config config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use: fmt.Sprintf("%s Dockerfile", cmdname.CheckPkgCmdName),
		Short: `check-pkg prints the difference of the packages versions between Dockerfile and the built image by it.
check-pkg requires 1 argument to specify the target Dockerfile.`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			packageInfos, err := checkPackageInformation(args[0])
			if err != nil {
				return err
			}

			err = outPackageVersionDiff(config.Out, packageInfos)
			if err != nil {
				return err
			}
			return nil
		},
	}
	cmd.SetIn(config.In)
	cmd.SetOut(config.Out)
	cmd.SetErr(config.Err)

	return cmd
}

type packageInfo struct {
	name           string
	versionInDfile string
	versionInImage string
}

func checkPackageInformation(dfilePath string) ([]packageInfo, error) {
	file, err := os.ReadFile(dfilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s: %w", dfilePath, err)
	}

	buf := bytes.NewBuffer(file)
	parsed, err := dfile.Parse(buf)
	if err != nil {
		return nil, fmt.Errorf("failed to parse %s as Dockerfile: %w", dfilePath, err)
	}

	aptPkgInfoMap, err := getAptPkgInfoInImageFromDfile(dfilePath)
	if err != nil {
		return nil, err
	}

	packageInfos := make([]packageInfo, 0)
	for _, node := range parsed.AST.Children {
		cmdInfos := dfile.NewCmdInfos(node.Next.Value)
		for _, cmdInfo := range cmdInfos {
			// only supports apt.
			if !(cmdInfo.MainCmd == "apt-get" && len(cmdInfo.Args) > 1 && cmdInfo.Args[0] == "install") {
				continue
			}

			for i, arg := range cmdInfo.Args {
				if i < 1 || strings.HasPrefix(arg, "-") {
					continue
				}
				split := strings.Split(arg, "=")
				packageName := split[0]
				info := packageInfo{
					name:           packageName,
					versionInImage: aptPkgInfoMap[packageName]}
				if len(split) == 2 {
					info.versionInDfile = split[1]
				}
				packageInfos = append(packageInfos, info)
			}
		}
	}

	return packageInfos, nil
}

func getAptPkgInfoInImageFromDfile(dfilePath string) (map[string]string, error) {
	imageTarFile, err := os.CreateTemp("/tmp", "dockerimage-*.tar")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp file: %w", err)
	}
	err = imageTarFile.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to close temp file: %w", err)
	}
	defer os.Remove(imageTarFile.Name())

	_, err = docker.CreateTarImageFromDockerfile(dfilePath, imageTarFile.Name())
	if err != nil {
		return nil, fmt.Errorf("failed to build image: %w", err)
	}
	imageTarFile, err = os.Open(imageTarFile.Name())
	reader := bufio.NewReader(imageTarFile)

	imageArchive, err := image.NewImageArchive(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to parse image tar file: %w", err)
	}

	aptPkgFile := imageArchive.GetLatestFileNode(pkginfo.AptPkgFilePath)
	if aptPkgFile == nil {
		return nil, fmt.Errorf("faild to get %s in the image", pkginfo.AptPkgFilePath)
	}

	buf := bytes.NewBuffer(aptPkgFile.Info.Data)
	aptPkgInfos, err := pkginfo.ReadAptPkgInfos(buf)

	aptPkgInfoMap := map[string]string{}
	for _, aptPkgInfo := range aptPkgInfos {
		aptPkgInfoMap[aptPkgInfo.Package] = aptPkgInfo.Version
	}

	return aptPkgInfoMap, nil
}

func outPackageVersionDiff(out io.Writer, packageInfos []packageInfo) error {
	for _, info := range packageInfos {
		if info.versionInDfile == info.versionInImage {
			continue
		}
		before := info.name
		if info.versionInDfile != "" {
			before = strings.Join([]string{info.name, info.versionInDfile}, "=")
		}
		after := strings.Join([]string{info.name, info.versionInImage}, "=")
		_, err := fmt.Fprintf(out, "%s => %s\n", before, after)
		if err != nil {
			return fmt.Errorf("faild to output: %w", err)
		}
	}

	return nil
}
