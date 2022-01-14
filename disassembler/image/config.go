package image

import (
	"encoding/json"
	dockerimage "github.com/docker/docker/image"
)

// Config stores the image configuration.
type Config struct {
	dockerimage.Image
}

func newConfig(configBytes []byte) (*Config, error) {
	var config Config
	err := json.Unmarshal(configBytes, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
