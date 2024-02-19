package checkpkg

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tklab-group/docker-image-disassembler/cli/cmdname"
	"github.com/tklab-group/docker-image-disassembler/cli/config"
	"github.com/tklab-group/docker-image-disassembler/disassembler"
	"github.com/tklab-group/docker-image-disassembler/disassembler/dfile"
)

const flagNameImageID = "imageID"

func Cmd(config config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("%s Dockerfile", cmdname.CheckPkgCmdName),
		Short: "Print the difference of the packages versions between Dockerfile and the built image by it",
		Long: `check-pkg prints the difference of the packages versions between Dockerfile and the built image by it.
check-pkg requires 1 argument to specify the target Dockerfile.`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			flagImageID := cmd.Flag(flagNameImageID)
			var imageID *string
			if flagImageID != nil && flagImageID.Value.String() != "" {
				_imageID := flagImageID.Value.String()
				imageID = &_imageID
			}

			packageInfos, err := checkPackageInformation(args[0], imageID)
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

	cmd.Flags().String(flagNameImageID, "", "specify docker image identifier to compare with Dockerfile")

	return cmd
}

type packageInfo struct {
	name           string
	versionInDfile string
	versionInImage string
}

func checkPackageInformation(dfilePath string, imageID *string) ([]packageInfo, error) {
	file, err := os.ReadFile(dfilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s: %w", dfilePath, err)
	}

	buf := bytes.NewBuffer(file)
	parsed, err := dfile.Parse(buf)
	if err != nil {
		return nil, fmt.Errorf("failed to parse %s as Dockerfile: %w", dfilePath, err)
	}

	var aptPkgInfoMap map[string]string
	if imageID != nil {
		aptPkgInfoMap, err = disassembler.GetAptPkgInfoInImageFromImageID(*imageID)
		if err != nil {
			return nil, err
		}
	} else {
		aptPkgInfoMap, err = disassembler.GetAptPkgInfoInImageFromDfile(dfilePath)
		if err != nil {
			return nil, err
		}
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
