package groups

import (
	"errors"
	"fmt"
	"github.com/ebauman/ldifgen/pkg/generators/names"
	"github.com/ebauman/ldifgen/pkg/types"
	"github.com/sirupsen/logrus"
	"math/rand"
	"strings"
)

type GroupGenerator struct {
	NameGenerator *names.NameGenerator
	OUList        []string
	Members       []string
	Domain        []string
	UserList      []*types.User
}

func New(domain []string, ng *names.NameGenerator, ouList []string, userList []*types.User) (*GroupGenerator, error) {
	if ng == nil {
		return nil, errors.New("undefined name generator")
	}

	return &GroupGenerator{NameGenerator: ng, OUList: ouList, Domain: domain, UserList: userList}, nil
}

func (gg *GroupGenerator) Generate(members []*types.User) (*types.Group, error) {
	newGroup := types.Group{}
	newGroup.CommonName = gg.NameGenerator.Group()
	ouIndex := rand.Intn(len(gg.OUList))
	newGroup.OrganizationalUnit = gg.OUList[ouIndex]
	newGroup.DistinguishedName = "cn=" + newGroup.CommonName + ",ou=" + newGroup.OrganizationalUnit + ",dc=" + strings.Join(gg.Domain, ",dc=")
	newGroup.Members = make([]string, 0)

	for _, m := range members {
		newGroup.Members = append(newGroup.Members, m.DistinguishedName)
	}

	return &newGroup, nil
}

func (gg *GroupGenerator) GenerateN(count int) ([]*types.Group, error) {
	userChunkSize := len(gg.UserList) / count

	groupList := make([]*types.Group, 0)
	groupMap := map[string]int{}
	userChunkPos := 0
	for i := 0; i < count; i++ {
		tryCount := 0
		members := gg.UserList[userChunkPos:((i + 1) * userChunkSize)]
		for {
			tempGroup, err := gg.Generate(members)
			if err != nil {
				return nil, fmt.Errorf("error generating group: %v", err)
			}

			if groupMap[tempGroup.DistinguishedName] == 1 {
				logrus.Infof("group name collision. generating new group")
				if tryCount > 5 {
					return nil, fmt.Errorf("too many name collisions generating groups")
				}
				tryCount++
			} else {
				groupMap[tempGroup.DistinguishedName] = 1
				groupList = append(groupList, tempGroup)
				userChunkPos = (i + 1) * userChunkSize
				break
			}
		}
	}

	return groupList, nil
}
