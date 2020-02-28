package users

import (
	"errors"
	"github.com/ebauman/ldifgen/pkg/generators/names"
	"github.com/ebauman/ldifgen/pkg/types"
	"math/rand"
)

type UserGenerator struct {
	NameGenerator *names.NameGenerator
	OUList [][]string
}

func New(nameGenerator *names.NameGenerator, ouList [][]string) (*UserGenerator, error) {
	if nameGenerator == nil {
		return nil, errors.New("undefined name generator")
	}

	return &UserGenerator{NameGenerator: nameGenerator, OUList:ouList}, nil
}

func (ug UserGenerator) Generate() (*types.User, error) {
	newUser := &types.User{}
	newUser.GivenName = ug.NameGenerator.FirstName()
	newUser.Surname = ug.NameGenerator.LastName()
	newUser.CommonName = newUser.GivenName + " " + newUser.Surname
	newUser.Description = "This is the description for " + newUser.CommonName + " " + newUser.Surname

	// pick a random OU to put someone in
	topLevelRand := rand.Intn(len(ug.OUList))
	depthRand := rand.Intn(len(ug.OUList[topLevelRand]))
	if depthRand == 0 {
		depthRand = 1
	}

	ouString := ""
	for i := depthRand; i > 0; i-- {
		ouString += "ou=" + ug.OUList[topLevelRand][i]
		if i-1 > 0 {
			ouString += ","
		}
	}

	newUser.OU = ouString

	return newUser, nil
}