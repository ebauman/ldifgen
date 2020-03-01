package ous

import (
	"github.com/ebauman/ldifgen/pkg/generators/names"
	"github.com/sirupsen/logrus"
)

func Generate(topLevelCount int, maxDepth int, nameGenerator *names.NameGenerator) []string {
	ouMap := map[string]int{}

	ouSlice := make([]string, 0)
	for i := 0; i < topLevelCount; i++ {
		ouString := ""
		tryCount := 0
		for j := 0; j < maxDepth; j++ {
			for {
				if j == 0 {
					ouString = nameGenerator.Department()
				} else {
					ouString = nameGenerator.Department() + ",ou=" + ouString
				}

				if ouMap[ouString] == 1 {
					// already exists, try again!
					logrus.Infof("OU name collision: %s. Generating new OU", ouString)
					if tryCount > 5 {
						logrus.Fatalf("error generating OUs, too many name collisions")
					} else {
						tryCount++
					}
				} else {
					ouMap[ouString] = 1
					ouSlice = append(ouSlice, ouString)
					break
				}
			}
		}
	}

	return ouSlice
}
