package image

import (
	dockerimage "github.com/docker/docker/image"
)

// Config stores the image configuration.
type Config struct {
	dockerimage.Image
}
