package dfile

import (
	"github.com/kr/pretty"
	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/require"
	"github.com/tklab-group/docker-image-disassembler/disassembler/testutil"
	"testing"
)

func TestParse(t *testing.T) {
	buf := testutil.ReadFileForBuffer(t, "testdata/Dockerfile")
	parsed, err := Parse(buf)
	require.NoError(t, err)

	g := goldie.New(t, goldie.WithFixtureDir("testdata/golden"))
	g.Assert(t, "parsed", []byte(pretty.Sprint(parsed)))
}
