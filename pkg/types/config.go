package types

import (
	"strings"
)

type GenerateConfig struct {
	Users                    int
	Groups                   int
	OUs                      int
	OUDepth                  int
	Domain                   []string
	UserClasses              []string
	GroupClasses             []string
	OUClasses                []string
	UserChangeType			 string
	GroupChangeType			 string
	OUChangeType			 string
	BuzzwordDataset          string
	DepartmentDataset        string
	FirstNameDataset         string
	LastNameDataset          string
	GroupsDataset            string
	GroupMembershipAttribute string
}

type RenderConfig struct {
	UserClasses              []string
	GroupClasses             []string
	OUClasses                []string
	UserChangeType			 string
	GroupChangeType			 string
	OUChangeType			 string
	Users                    []*User
	Domain                   []string
	OUs                      []string
	Groups                   []*Group
	GroupMembershipAttribute string
	Time                     string
}

func (c RenderConfig) DC() string {
	var domain = ""
	domain = "dc=" + strings.Join(c.Domain, ",dc=")
	return domain
}

func (c RenderConfig) TrimOU(ou string) string {
	return strings.Split(ou, ",")[0]
}
