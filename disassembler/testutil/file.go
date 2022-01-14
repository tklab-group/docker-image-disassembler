package testutil

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func ReadFileForBuffer(t *testing.T, filepath string) *bytes.Buffer {
	t.Helper()

	file, err := os.ReadFile(filepath)
	require.NoError(t, err)
	return bytes.NewBuffer(file)
}
