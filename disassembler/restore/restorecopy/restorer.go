package restorecopy

import (
	"bytes"
	"fmt"
	"github.com/moby/buildkit/frontend/dockerfile/instructions"
	"github.com/moby/buildkit/frontend/dockerfile/parser"
	"github.com/tklab-group/docker-image-disassembler/disassembler/image"
	"path/filepath"
	"regexp"
	"strings"
)

type restorer struct {
	tmpWorkDir string
}

func newRestorer() *restorer {
	return &restorer{tmpWorkDir: "/"}
}

func (r *restorer) restore(imageArchive *image.ImageArchive) ([]*CopiedObject, error) {
	historyToLayers, err := imageArchive.GetHistoryToLayers()
	if err != nil {
		return nil, err
	}

	copiedObjects := make([]*CopiedObject, 0)
	for _, historyToLayer := range historyToLayers {
		command := strings.ReplaceAll(historyToLayer.History.CreatedBy, "# buildkit", "")

		// Support format like `/bin/sh -c #(nop) COPY ... in ...`
		if strings.HasPrefix(command, "/bin/sh -c") {
			if isExtraFormatCopy(command) {
				copiedObject, err := r.handleExtraFormatCopy(command, historyToLayer)
				if err != nil {
					return nil, err
				}

				copiedObjects = append(copiedObjects, copiedObject)
			}

			continue
		}

		cmdBuf := bytes.NewBuffer([]byte(command))
		parsed, err := parser.Parse(cmdBuf)
		if err != nil {
			return nil, err
		}

		parsedNode := parsed.AST.Children[0]
		switch strings.ToLower(parsedNode.Value) {
		case "workdir":
			err = r.handleWorkdir(parsedNode)
			if err != nil {
				return nil, err
			}
		case "copy":
			copiedObject, err := r.handleCopy(parsedNode, historyToLayer)
			if err != nil {
				return nil, err
			}
			copiedObjects = append(copiedObjects, copiedObject)
		}
	}

	return copiedObjects, nil
}

func (r *restorer) handleWorkdir(node *parser.Node) error {
	cmd, err := instructions.ParseCommand(node)
	if err != nil {
		return err
	}

	workDirCmd, ok := cmd.(*instructions.WorkdirCommand)
	if !ok {
		return fmt.Errorf("cann't parse as WorkDirCommand: %s", cmd.Name())
	}

	r.tmpWorkDir = r.absPath(workDirCmd.Path)

	return nil
}

func (r *restorer) handleCopy(node *parser.Node, historyToLayer *image.HistoryToLayer) (*CopiedObject, error) {
	cmd, err := instructions.ParseCommand(node)
	if err != nil {
		return nil, err
	}

	copyCmd, ok := cmd.(*instructions.CopyCommand)
	if !ok {
		return nil, fmt.Errorf("cann't parse as CopyCommand: %s", cmd.Name())
	}

	return r.genCopiedObject(copyCmd.DestPath, historyToLayer)
}

func (r *restorer) handleExtraFormatCopy(command string, historyToLayer *image.HistoryToLayer) (*CopiedObject, error) {
	destPath, err := getDestPathFromExtraFormatCopy(command)
	if err != nil {
		return nil, err
	}
	destPath = strings.TrimSpace(destPath)

	return r.genCopiedObject(destPath, historyToLayer)
}

func (r *restorer) genCopiedObject(destPath string, historyToLayer *image.HistoryToLayer) (*CopiedObject, error) {
	copiedObject := &CopiedObject{
		LayerID: historyToLayer.LayerID,
		History: historyToLayer.History,
	}

	copiedPath := r.absPath(destPath)
	fileNode := historyToLayer.Layer.FindNodeFromPath(copiedPath)
	if fileNode == nil {
		return nil, fmt.Errorf("can't find copied object %s", copiedPath)
	}

	copiedObject.Object = fileNode.Copy(nil)

	return copiedObject, nil
}

func isExtraFormatCopy(command string) bool {
	reg := regexp.MustCompile(`/bin/sh -c #\(nop\)\s+COPY .+ in .+`)
	return reg.MatchString(command)
}

func getDestPathFromExtraFormatCopy(command string) (string, error) {
	reg := regexp.MustCompile(`/bin/sh -c #\(nop\)\s+COPY .+ in (.+)`)
	r := reg.FindAllStringSubmatch(command, -1)
	if len(r) != 1 && len(r[0]) != 2 {
		return "", fmt.Errorf("faild to get destination path: `%s`", command)
	}

	return r[0][1], nil
}

func (r *restorer) absPath(path string) string {
	if strings.HasPrefix(path, "/") {
		return path
	}

	return filepath.Join(r.tmpWorkDir, path)
}
