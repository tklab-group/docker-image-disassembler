package image

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"github.com/tklab-group/docker-image-disassembler/disassenmbler/filetree"
	"io"
	"path"
	"strings"
)

type ImageArchive struct {
	Manifest *Manifest
	Config   *Config
	LayerMap map[string]*filetree.FileTree
}

// NewImageArchive creates ImageArchive from a tar file.
func NewImageArchive(tarFile io.Reader) (*ImageArchive, error) {
	img := &ImageArchive{
		LayerMap: map[string]*filetree.FileTree{},
	}

	tarReader := tar.NewReader(tarFile)

	// Store discovered json files in a map, so we can read the image in one pass.
	jsonFiles := map[string][]byte{}

	for true {
		header, err := tarReader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}

			return nil, err
		}

		name := header.Name

		// Some layer tars can be relative layer symlinks to other layer tars.
		if header.Typeflag == tar.TypeSymlink || header.Typeflag == tar.TypeReg {
			if strings.HasSuffix(name, ".tar") {
				layerReader := tar.NewReader(tarReader)
				tree, err := processLayerTar(name, layerReader)
				if err != nil {
					return nil, err
				}

				img.LayerMap[tree.LayerName] = tree
			} else if strings.HasSuffix(name, ".tar.gz") || strings.HasSuffix(name, "tgz") {
				gzipReader, err := gzip.NewReader(tarReader)
				if err != nil {
					return nil, err
				}

				layerReader := tar.NewReader(gzipReader)
				tree, err := processLayerTar(name, layerReader)
				if err != nil {
					return nil, err
				}

				img.LayerMap[tree.LayerName] = tree
			} else if strings.HasSuffix(name, ".json") || strings.HasPrefix(name, "sha256:") {
				fileBuffer, err := io.ReadAll(tarReader)
				if err != nil {
					return nil, err
				}

				jsonFiles[name] = fileBuffer
			}
		}
	}

	manifestContent, exists := jsonFiles["manifest.json"]
	if !exists {
		return nil, fmt.Errorf("could not find image manifest")
	}

	manifest, err := newManifest(manifestContent)
	if err != nil {
		return nil, fmt.Errorf("failed to parse manifest.json: %w", err)
	}
	img.Manifest = manifest

	configContent, exists := jsonFiles[img.Manifest.ConfigPath]
	if !exists {
		return nil, fmt.Errorf("could not find image config")
	}

	config, err := newConfig(configContent)
	if err != nil {
		return nil, fmt.Errorf("failed to parse %s as config: %w", img.Manifest.ConfigPath, err)
	}
	img.Config = config

	return img, nil
}

// GetFileTreeByLayerIndex returns FileTree for the layer specified by index.
func (img *ImageArchive) GetFileTreeByLayerIndex(index int) (*filetree.FileTree, error) {
	if index < 0 || index >= len(img.Manifest.LayerTarPaths) {
		return nil, fmt.Errorf("index %d is out of range", index)
	}

	layerName := img.Manifest.LayerTarPaths[index]
	fileTree := img.LayerMap[layerName]

	return fileTree, nil
}

// GetLatestFileNode searches FileNode based on the path and returns the latest one.
// If the path doesn't exist in all layers, it returns nil.
func (img *ImageArchive) GetLatestFileNode(path string) *filetree.FileNode {
	for i := len(img.Manifest.LayerTarPaths) - 1; i >= 0; i-- {
		fileTree, _ := img.GetFileTreeByLayerIndex(i)
		fileNode := fileTree.FindNodeFromPath(path)
		if fileNode != nil {
			return fileNode
		}
	}

	return nil
}

func processLayerTar(name string, tarReader *tar.Reader) (*filetree.FileTree, error) {
	tree := filetree.NewFileTree()
	tree.LayerName = name

	fileInfos, err := getFileList(tarReader)
	if err != nil {
		return nil, err
	}

	for _, fileInfo := range fileInfos {
		err = tree.AddNode(fileInfo)
		if err != nil {
			return nil, err
		}
	}

	return tree, nil
}

func getFileList(tarReader *tar.Reader) ([]*filetree.FileInfo, error) {
	var files []*filetree.FileInfo

	for {
		header, err := tarReader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}

			return nil, err
		}

		// Always ensure relative path notations are not parsed as part of the filename.
		name := path.Clean(header.Name)
		if name == "." {
			continue
		}

		switch header.Typeflag {
		case tar.TypeXGlobalHeader:
			return nil, fmt.Errorf("unexptected tar file: (XGlobalHeader): type=%v name=%s", header.Typeflag, name)
		case tar.TypeXHeader:
			return nil, fmt.Errorf("unexptected tar file (XHeader): type=%v name=%s", header.Typeflag, name)
		default:
			file, err := filetree.NewFileInfoFromTarHeader(tarReader, header)
			if err != nil {
				return nil, err
			}

			files = append(files, file)
		}
	}

	return files, nil
}
