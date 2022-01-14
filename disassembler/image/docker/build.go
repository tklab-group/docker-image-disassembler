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

// CreateTarImageFromDockerfile builds docker image from dockerfile and exports the image as a tar file.
func CreateTarImageFromDockerfile(dfilePath string, tarPath string) (imageID string, err error) {
	iid, err := BuildImageFromCli([]string{"-f", dfilePath, "."})
	if err != nil {
		return "", err
	}

	err = RunDockerCmd("save", []string{iid, "-o", tarPath}, nil)
	if err != nil {
		return "", err
	}

	return iid, nil
}
