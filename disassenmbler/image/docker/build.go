package docker

import (
	"os"
)

// BuildImageFromCli returns the id for the built image.
func BuildImageFromCli(buildArgs []string) (string, error) {
	iidfile, err := os.CreateTemp("/tmp", "iid")
	if err != nil {
		return "", err
	}
	defer os.Remove(iidfile.Name())

	allArgs := append([]string{"--iidfile", iidfile.Name()}, buildArgs...)
	err = RunDockerCmd("build", allArgs, nil)
	if err != nil {
		return "", err
	}

	imageId, err := os.ReadFile(iidfile.Name())
	if err != nil {
		return "", err
	}

	return string(imageId), nil
}
