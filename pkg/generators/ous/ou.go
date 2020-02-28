package ous

import (
	"errors"
	"github.com/ebauman/ldifgen/pkg/generators/names"
)

type OUGenerator struct {
	MaxDepth int
	NameGenerator *names.NameGenerator
}

func New(maxDepth int, nameGenerator *names.NameGenerator) (*OUGenerator, error) {
	if nameGenerator == nil {
		return nil, errors.New("undefined name generator")
	}

	if maxDepth == 0 {
		return nil, errors.New("max depth must be >= 1")
	}

	return &OUGenerator{MaxDepth:maxDepth, NameGenerator:nameGenerator}, nil
}

func (oug OUGenerator) Generate() []string {
	ouSlice := make([]string, 0)
	for i := 0; i < oug.MaxDepth; i++ {
		ou := oug.NameGenerator.Department()
		ouSlice = append(ouSlice, ou)
	}

	return ouSlice
}