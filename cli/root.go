package cli

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tklab-group/docker-image-disassembler/cli/checkpkg"
	"github.com/tklab-group/docker-image-disassembler/cli/cmdname"
	"github.com/tklab-group/docker-image-disassembler/cli/config"
	"github.com/tklab-group/docker-image-disassembler/cli/restorecopy"
	"os"
)

func newRoodCmd(config config.Config) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   cmdname.RootCmdName,
		Short: "", // TODO
		Long:  "", // TODO
	}
	rootCmd.SetIn(config.In)
	rootCmd.SetOut(config.Out)
	rootCmd.SetErr(config.Err)

	rootCmd.AddCommand(
		checkpkg.Cmd(config),
		restorecopy.Cmd(config),
	)

	return rootCmd
}

func Execute(config config.Config) {
	if err := newRoodCmd(config).Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
