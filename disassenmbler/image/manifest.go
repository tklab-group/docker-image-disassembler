package image

import "encoding/json"

// Manifest stores information from manifest.json in the image compressed as a tar file.
type Manifest struct {
	ConfigPath    string   `json:"Config"`
	RepoTags      []string `json:"RepoTags"`
	LayerTarPaths []string `json:"Layers"`
}

func newManifest(manifestBytes []byte) (*Manifest, error) {
	var manifest []Manifest
	err := json.Unmarshal(manifestBytes, &manifest)
	if err != nil {
		return nil, err
	}
	return &manifest[0], nil
}
