package listPkg

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tklab-group/docker-image-disassembler/cli/cmdname"
	"github.com/tklab-group/docker-image-disassembler/cli/config"
	"github.com/tklab-group/docker-image-disassembler/disassembler"
)

func Cmd(config config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("%s imageID", cmdname.ListPkgCmdName),
		Short: `list-pkg prints packages and their versions in the image`,
		Long:  `Currently only packages installed with apt are supported`,
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			imageID := args[0]

			packages := map[string]map[string]string{}

			aptPackages, err := disassembler.GetAptPkgInfoInImageFromImageID(imageID)
			if err != nil {
				return err
			}

			packages["apt"] = aptPackages

			b, err := json.MarshalIndent(packages, "", "\t")
			if err != nil {
				return err
			}

			fmt.Fprintf(config.Out, "%s", string(b))

			return nil
		},
	}

	cmd.SetIn(config.In)
	cmd.SetOut(config.Out)
	cmd.SetErr(config.Err)

	return cmd
}
