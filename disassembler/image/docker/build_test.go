package docker

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBuildImageFromCli(t *testing.T) {
	iid, err := BuildImageFromCli([]string{"-f", " testdata/Dockerfile.echo-hello", "."})
	assert.NoError(t, err)
	assert.NotEmpty(t, iid)

	outBuf := bytes.Buffer{}
	stds := &RunDockerCmdStds{Stdout: &outBuf}
	err = RunDockerCmd("run", []string{"--rm", iid}, stds)
	assert.NoError(t, err)
	assert.Equal(t, "hello\n", outBuf.String())
}
