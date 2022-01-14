package docker

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestRunDockerCmd(t *testing.T) {
	outBuf := bytes.Buffer{}
	errBuf := bytes.Buffer{}
	stds := &RunDockerCmdStds{
		Stdout: &outBuf,
		Stderr: &errBuf,
		Stdin:  nil,
	}

	err := RunDockerCmd("version", nil, stds)
	if assert.NoError(t, err) {
		assert.Empty(t, errBuf.String())
		outStr := outBuf.String()
		assert.True(t, strings.Contains(outStr, "Client:"))
		assert.True(t, strings.Contains(outStr, "Server:"))
	}
}
