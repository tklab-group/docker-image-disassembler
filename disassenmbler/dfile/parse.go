package dfile

import (
	"github.com/moby/buildkit/frontend/dockerfile/parser"
	"io"
)

type Parsed struct {
	*parser.Result
}

func Parse(rwc io.Reader) (*Parsed, error) {
	result, err := parser.Parse(rwc)
	if err != nil {
		return nil, err
	}

	return &Parsed{result}, nil
}
