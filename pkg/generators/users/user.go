package users

import (
	"errors"
	"fmt"
	"github.com/ebauman/ldifgen/pkg/generators/names"
	"github.com/ebauman/ldifgen/pkg/types"
	"github.com/sirupsen/logrus"
	"math/rand"
	"strings"
)

type UserGenerator struct {
	NameGenerator *names.NameGenerator
	OUList        []string
	Domain        []string
}

func New(domain []string, nameGenerator *names.NameGenerator, ouList []string) (*UserGenerator, error) {
	if nameGenerator == nil {
		return nil, errors.New("undefined name generator")
	}

	return &UserGenerator{NameGenerator: nameGenerator, OUList: ouList, Domain: domain}, nil
}

func (ug *UserGenerator) Generate() (*types.User, error) {
	newUser := &types.User{}
	newUser.GivenName = ug.NameGenerator.FirstName()
	newUser.Surname = ug.NameGenerator.LastName()
	newUser.CommonName = newUser.GivenName + " " + newUser.Surname
	newUser.Description = "This is the description for " + newUser.CommonName + " " + newUser.Surname
	newUser.OU = ug.OUList[rand.Intn(len(ug.OUList))]
	newUser.DistinguishedName = "cn=" + newUser.CommonName + ",ou=" + newUser.OU + ",dc=" + strings.Join(ug.Domain, ",dc=")

	return newUser, nil
}

func (ug *UserGenerator) GenerateN(count int) ([]*types.User, error) {
	userList := make([]*types.User, 0)
	userMap := map[string]int{}
	for i := 0; i < count; i++ {
		tryCount := 0
		for {
			tempUser, err := ug.Generate()
			if err != nil {
				return nil, fmt.Errorf("error generating user: %v", err)
			}
			if userMap[tempUser.UID()] == 1 {
				logrus.Infof("user uid collision. generating new user")
				if tryCount > 5 {
					return nil, fmt.Errorf("too many name collisions generating users")
				}
				tryCount++
			} else {
				userMap[tempUser.UID()] = 1
				userList = append(userList, tempUser)
				break
			}
		}
	}
	return userList, nil
}
