package ous

import (
	"github.com/ebauman/ldifgen/pkg/generators/names"
)

func Generate(topLevelCount int, maxDepth int, nameGenerator *names.NameGenerator) []string {
	ouSlice := make([]string, 0)
	for i := 0; i < topLevelCount; i++ {
		ouString := ""
		for j := 0; j < maxDepth; j++ {
			if j == 0 {
				ouString = nameGenerator.Department()
			} else {
				ouString = nameGenerator.Department() + ",ou=" + ouString
			}
			ouSlice = append(ouSlice, ouString)
		}
	}

	return ouSlice
}