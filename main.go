package main

import (
	"github.com/ebauman/ldifgen/pkg/generators/names"
	"github.com/ebauman/ldifgen/pkg/generators/ous"
	users2 "github.com/ebauman/ldifgen/pkg/generators/users"
	"github.com/ebauman/ldifgen/pkg/types"
	"html/template"
	"log"
	"os"
)

func main() {
	tmpl, err := template.ParseFiles("template/users.txt")
	if err != nil {
		log.Fatalf("%v", err)
	}


	nameGen, err := names.NewNameGenerator()
	if err != nil {
		log.Fatalf("error creating name generator: %v", err)
	}

	ouGen, err := ous.New(5, nameGen)
	if err != nil {
		log.Fatalf("error creating ous generator: %v", err)
	}

	ouList := make([][]string, 0)
	for i := 0; i < 10; i++ {
		tempOU := ouGen.Generate()
		ouList = append(ouList, tempOU)
	}

	userGen, err := users2.New(nameGen, ouList)
	if err != nil {
		log.Fatalf("error creating user generator: %v", err)
	}

	users := make([]*types.User, 0)
	for i := 0; i < 10000; i++ {
		tempUser, err := userGen.Generate()
		if err != nil {
			log.Printf("error generating user: %v", err)
		}
		users = append(users, tempUser)
	}

	config := types.RenderConfig{Users: users, Domain: []string{"testing", "rancher", "com"}, UserClasses: []string{"top", "person", "organizationalPerson", "inetOrgPerson"}}
	err = tmpl.Execute(os.Stdout, config)
	if err != nil {
		log.Fatalf("%v", err)
	}
}
