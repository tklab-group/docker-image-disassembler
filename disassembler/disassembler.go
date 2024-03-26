package disassembler

import (
	"bufio"
	"bytes"
	"fmt"
	"os"

	"github.com/tklab-group/docker-image-disassembler/disassembler/image"
	"github.com/tklab-group/docker-image-disassembler/disassembler/image/docker"
	"github.com/tklab-group/docker-image-disassembler/disassembler/pkginfo"
)

// TODO: Enhance to other package managers

func GetAptPkgInfoInImageFromImageID(imageID string) (map[string]string, error) {
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
		return nil, fmt.Errorf("failed to save image: %w", err)
	}

	imageTarFile, err = os.Open(imageTarFile.Name())
	if err != nil {
		return nil, fmt.Errorf("failed to open temp file: %w", err)
	}

	reader := bufio.NewReader(imageTarFile)

	imageArchive, err := image.NewImageArchive(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to parse image tar file: %w", err)
	}

	return GetAptPkgInfoInImageFromImageArchive(imageArchive)
}

func GetAptPkgInfoInImageFromDfile(dfilePath string) (map[string]string, error) {
	imageTarFile, err := os.CreateTemp("/tmp", "dockerimage-*.tar")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp file: %w", err)
	}
	err = imageTarFile.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to close temp file: %w", err)
	}
	defer os.Remove(imageTarFile.Name())

	_, err = docker.CreateTarImageFromDockerfile(dfilePath, imageTarFile.Name())
	if err != nil {
		return nil, fmt.Errorf("failed to build image: %w", err)
	}
	imageTarFile, err = os.Open(imageTarFile.Name())
	reader := bufio.NewReader(imageTarFile)

	imageArchive, err := image.NewImageArchive(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to parse image tar file: %w", err)
	}

	return GetAptPkgInfoInImageFromImageArchive(imageArchive)
}

func GetAptPkgInfoInImageFromImageArchive(imageArchive *image.ImageArchive) (map[string]string, error) {
	aptPkgFile := imageArchive.GetLatestFileNode(pkginfo.AptPkgFilePath)
	if aptPkgFile == nil {
		return nil, fmt.Errorf("faild to get %s in the image", pkginfo.AptPkgFilePath)
	}

	buf := bytes.NewBuffer(aptPkgFile.Info.Data)
	aptPkgInfos, err := pkginfo.ReadAptPkgInfos(buf)
	if err != nil {
		return nil, fmt.Errorf("faild to read apt package file: %w", err)
	}

	aptPkgInfoMap := map[string]string{}
	for _, aptPkgInfo := range aptPkgInfos {
		aptPkgInfoMap[aptPkgInfo.Package] = aptPkgInfo.Version
	}

	return aptPkgInfoMap, nil
}
