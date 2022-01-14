package cli

import (
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"os"
)

type Config struct {
	In  io.Reader
	Out io.Writer
	Err io.Writer
}

const RootCmdName = "docker-image-disassembler"

func newRoodCmd(config Config) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   RootCmdName, // TODO: Change the command name better.
		Short: "",          // TODO
		Long:  "",          // TODO
	}
	rootCmd.SetIn(config.In)
	rootCmd.SetOut(config.Out)
	rootCmd.SetErr(config.Err)

	rootCmd.AddCommand(
	// TODO
	)

	return rootCmd
}

func Execute(config Config) {
	if err := newRoodCmd(config).Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
