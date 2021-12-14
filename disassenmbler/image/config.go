package image

import (
	dimage "github.com/docker/docker/image"
)

// Config stores the image configuration.
type Config struct {
	dimage.Image
}
