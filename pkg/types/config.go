package types

import "fmt"

type GenerateConfig struct {
	Users int
	Groups int
	OUs int
	Domain string
	UserClasses []string
	GroupClasses []string
	GroupMembershipAttribute string
}

type RenderConfig struct {
	UserClasses              []string
	GroupClasses             []string
	Users                    []*User
	Domain                   []string
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
