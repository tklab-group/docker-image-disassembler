package restorecopy

import (
	"bytes"
	"fmt"
	"github.com/moby/buildkit/frontend/dockerfile/instructions"
	"github.com/moby/buildkit/frontend/dockerfile/parser"
	"github.com/tklab-group/docker-image-disassembler/disassembler/image"
	"path/filepath"
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

	copiedObject := &CopiedObject{
		LayerID: historyToLayer.LayerID,
		History: historyToLayer.History,
	}

	copiedPath := r.absPath(copyCmd.DestPath)
	fileNode := historyToLayer.Layer.FindNodeFromPath(copiedPath)
	if err != nil {
		return nil, fmt.Errorf("can't find copied object %s", copiedPath)
	}

	copiedObject.Object = fileNode.Copy(nil)

	return copiedObject, nil
}

func (r *restorer) absPath(path string) string {
	if strings.HasPrefix(path, "/") {
		return path
	}

	return filepath.Join(r.tmpWorkDir, path)
}
