package restorecopy

import (
	"bufio"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tklab-group/docker-image-disassembler/cli/cmdname"
	"github.com/tklab-group/docker-image-disassembler/cli/config"
	"github.com/tklab-group/docker-image-disassembler/disassembler/image"
	"github.com/tklab-group/docker-image-disassembler/disassembler/image/docker"
	"github.com/tklab-group/docker-image-disassembler/disassembler/restore"
	"github.com/tklab-group/docker-image-disassembler/disassembler/restore/restorecopy"
	"io"
	"os"
	"path/filepath"
)

func Cmd(config config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("%s imageID targetPath", cmdname.RestoreCopyCmdName),
		Short: `restore-copy extracts copied files from the image and embodies them at the target path.`,
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			copiedObjects, err := restoreCopiedObjects(args[0])
			if err != nil {
				return err
			}

			err = embodyCopiedObjects(args[1], copiedObjects)
			if err != nil {
				return err
			}

			if len(copiedObjects) == 0 {
				return fmt.Errorf("there is no copy instruction")
			}

			err = outLayerID(config.Out, copiedObjects)
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

func restoreCopiedObjects(imageID string) ([]*restorecopy.CopiedObject, error) {
	imageTarFile, err := os.CreateTemp("/tmp", "dockerimage-*.tar")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp file: %w", err)
	}
	err = imageTarFile.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to close temp file: %w", err)
	}
	defer os.Remove(imageTarFile.Name())

	err = docker.RunDockerCmd("save", []string{imageID, "-o", imageTarFile.Name()}, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch image by `%s`: %w", imageID, err)
	}

	imageTarFile, err = os.Open(imageTarFile.Name())
	reader := bufio.NewReader(imageTarFile)

	imageArchive, err := image.NewImageArchive(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to parse image tar file: %w", err)
	}

	copiedObjects, err := restorecopy.RestoreCopiedObjects(imageArchive)
	if err != nil {
		return nil, fmt.Errorf("failed to restore copied objects: %w", err)
	}

	return copiedObjects, nil
}

func embodyCopiedObjects(targetPath string, copiedObjects []*restorecopy.CopiedObject) error {
	if !existDir(targetPath) {
		return fmt.Errorf("cannot use as a target directory: %s", targetPath)
	}

	for _, copiedObject := range copiedObjects {
		pathForLayer := filepath.Join(targetPath, copiedObject.LayerID)
		err := os.Mkdir(pathForLayer, 0777)
		if err != nil {
			return fmt.Errorf("failed to create a directory: %w", err)
		}

		err = restore.EmbodyFileNode(pathForLayer, copiedObject.Object)
		if err != nil {
			return fmt.Errorf("failed to embody copied objects: %w", err)
		}
	}

	return nil
}

func outLayerID(out io.Writer, copiedObjects []*restorecopy.CopiedObject) error {
	for _, copiedObject := range copiedObjects {
		_, err := fmt.Fprintf(out, "%s: `%s`\n", copiedObject.LayerID, copiedObject.History.CreatedBy)
		if err != nil {
			return fmt.Errorf("faild to output: %w", err)
		}
	}

	return nil
}

func existDir(dirPath string) bool {
	f, err := os.Stat(dirPath)
	if err == nil && f.IsDir() {
		return true
	}

	return false
}
