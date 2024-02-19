# Docker Image Disassembler

`docker-image-disassembler` is a CLI tool for extracting package versions and copied files from Docker image.

For package versions, it currently supports only packages installed via `apt`.

> [!WARNING]
> Currently it doesn't support new Docker image format.
>
> Use `docker-image-disassembler` with Docker Engine older version than `v25.0`.

## Usage

### `list-pkg`
```shell
docker-image-disassembler list-pkg imageID [flags]
```

Prints the packages and their versions in the image as JSON format.

### `check-pkg`
```shell
docker-image-disassembler check-pkg Dockerfile [flags]
```

Prints the difference of the packages versions between Dockerfile and the built image by it.

### `restore-copy`
```shell
docker-image-disassembler restore-copy imageID targetPath [flags]
```

Extracts files, copied to the image by COPY instruction, from the image and embodies them at the target path.

## Installation

### Go tools

```shell
go install github.com/tklab-group/docker-image-disassembler@latest
```

### Binary

Download from [release page](https://github.com/tklab-group/docker-image-disassembler/releases/latest).