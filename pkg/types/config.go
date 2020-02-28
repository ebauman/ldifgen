package types

import (
	"fmt"
	"strings"
)

type GenerateConfig struct {
	Users                    int
	Groups                   int
	OUs                      int
	Domain                   string
	UserClasses              []string
	GroupClasses             []string
	GroupMembershipAttribute string
}

type RenderConfig struct {
	UserClasses              []string
	GroupClasses             []string
	OUClasses                []string
	Users                    []*User
	Domain                   []string
	OUs                      []string
	GroupMembershipAttribute string
}

func (c RenderConfig) DC() string {
	var domain = ""
	for index, domainPart := range c.Domain {
		domain += fmt.Sprintf("dc=%s", domainPart)
		if index != len(c.Domain)-1 {
			// this isn't the last segment
			domain += ","
		}
	}
	return domain
}

func (c RenderConfig) TrimOU(ou string) string {
	return strings.Split(ou, ",")[0]
}
