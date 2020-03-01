package cmd

import (
	"errors"
	"fmt"
	"github.com/ebauman/ldifgen/pkg/generators/groups"
	"github.com/ebauman/ldifgen/pkg/generators/names"
	"github.com/ebauman/ldifgen/pkg/generators/ous"
	"github.com/ebauman/ldifgen/pkg/generators/users"
	_ "github.com/ebauman/ldifgen/pkg/statik"
	"github.com/ebauman/ldifgen/pkg/types"
	"github.com/rakyll/statik/fs"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"html/template"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"time"
)

const datasetString = "path to an alternative list of %s, used in %s generation. provide list of words, separated by newlines"

func GenerateCommand() *cli.Command {
	generateFlags := []cli.Flag{
		&cli.IntFlag{
			Name:  "users",
			Value: 10,
			Usage: "number of users to generate",
		},
		&cli.IntFlag{
			Name:  "ous",
			Value: 2,
			Usage: "number of organizational units to generate",
		},
		&cli.IntFlag{
			Name:  "ou-depth",
			Value: 1,
			Usage: "depth of generated OUs. specify n>1 to create 'chains' of OUs",
		},
		&cli.IntFlag{
			Name:  "groups",
			Value: 2,
			Usage: "number of groups to generate",
		},
		&cli.StringFlag{
			Name:  "domain",
			Value: "domain.example.org",
			Usage: "domain used to generate DC components, e.g. dc=domain,dc=example,dc=org",
		},
		&cli.StringFlag{
			Name:  "user-classes",
			Value: "top,person,organizationalPerson,inetOrgPerson",
			Usage: "comma-separated list of classes for user objects",
		},
		&cli.StringFlag{
			Name:  "ou-classes",
			Value: "top,organizationalUnit",
			Usage: "comma-separated list of classes for organizational unit objects",
		},
		&cli.StringFlag{
			Name:  "group-classes",
			Value: "top,groupOfNames",
			Usage: "comma-separated list of classes for group objects",
		},
		&cli.StringFlag{
			Name: "user-change-type",
			Value: "add",
			Usage: "LDIF changetype for users",
		},
		&cli.StringFlag{
			Name: "group-change-type",
			Value: "add",
			Usage: "LDIF changetype for groups",
		},
		&cli.StringFlag{
			Name: "ou-change-type",
			Value: "add",
			Usage: "LDIF changetype for OUs",
		},
		&cli.StringFlag{
			Name:  "buzzword-dataset",
			Usage: fmt.Sprintf(datasetString, "buzzwords", "group"),
		},
		&cli.StringFlag{
			Name:  "department-dataset",
			Usage: fmt.Sprintf(datasetString, "department names", "OU"),
		},
		&cli.StringFlag{
			Name:  "first-name-dataset",
			Usage: fmt.Sprintf(datasetString, "first names", "user"),
		},
		&cli.StringFlag{
			Name:  "last-name-dataset",
			Usage: fmt.Sprintf(datasetString, "last names", "user"),
		},
		&cli.StringFlag{
			Name:  "groups-dataset",
			Usage: fmt.Sprintf(datasetString, "group names", "group"),
		},
	}

	return &cli.Command{
		Name:   "generate",
		Usage:  "generate ldif file",
		Action: generateLdif,
		Flags:  generateFlags,
	}
}

func doGenerate(gconf *types.GenerateConfig) error {
	statikFS, err := fs.New()
	if err != nil {
		logrus.Fatalf("error building statik fs: %v", err)
	}

	r, err := statikFS.Open("/template/ldif.txt")
	if err != nil {
		logrus.Fatalf("error opening ldif template: %v", err)
	}

	defer r.Close()

	data, err := ioutil.ReadAll(r)
	if err != nil {
		logrus.Fatalf("error reading ldif template: %v", err)
	}

	tmpl, err := template.New("ldif").Parse(string(data))
	if err != nil {
		logrus.Fatalf("error parsing ldif template: %v", err)
	}

	nameGen, err := names.NewNameGenerator(gconf.FirstNameDataset, gconf.LastNameDataset, gconf.DepartmentDataset, gconf.BuzzwordDataset, gconf.GroupsDataset)
	if err != nil {
		logrus.Fatalf("error creating name generator: %v", err)
	}

	ouList := ous.Generate(gconf.OUs, gconf.OUDepth, nameGen)

	userGen, err := users.New(gconf.Domain, nameGen, ouList)
	if err != nil {
		logrus.Fatalf("error creating user generator: %v", err)
	}

	userList, err := userGen.GenerateN(gconf.Users)
	if err != nil {
		logrus.Fatalf("error generating users: %v", err)
	}

	groupGen, err := groups.New(gconf.Domain, nameGen, ouList, userList)
	if err != nil {
		logrus.Fatalf("error creating group generator: %v", err)
	}

	groupList, err := groupGen.GenerateN(gconf.Groups)
	if err != nil {
		logrus.Fatalf("error generating groups: %v", err)
	}

	renderConfig := types.RenderConfig{
		Users: userList,
		Domain: gconf.Domain,
		UserChangeType: gconf.UserChangeType,
		GroupChangeType: gconf.GroupChangeType,
		OUChangeType: gconf.OUChangeType,
		UserClasses: gconf.UserClasses,
		GroupClasses: gconf.GroupClasses,
		OUClasses: gconf.OUClasses,
		OUs: ouList,
		Groups: groupList,
		Time: time.Now().Format("2006-01-02T15:04:05-0700"),
	}

	err = tmpl.Execute(os.Stdout, renderConfig)
	if err != nil {
		logrus.Fatalf("error executing template: %v", err)
	}

	return nil
}

func generateLdif(ctx *cli.Context) error {
	domainList, err := parseDomain(ctx.String("domain"))
	if err != nil {
		logrus.Fatalf("error parsing domain: %v", err)
	}

	if ctx.Int("ou-depth") < 1 {
		logrus.Fatalf("invalid ou depth (<1): %v", ctx.Int("ou-depth"))
	}

	userClassList, err := parseClassList(ctx.String("user-classes"))
	if err != nil {
		logrus.Fatalf("error parsing user classes: %v", err)
	}

	ouClassList, err := parseClassList(ctx.String("ou-classes"))
	if err != nil {
		logrus.Fatalf("error parsing ou classes: %v", err)
	}

	groupClassList, err := parseClassList(ctx.String("group-classes"))
	if err != nil {
		logrus.Fatalf("error parsing group classes: %v", err)
	}

	if ctx.String("user-change-type") == "" {
		logrus.Fatalf("invalid user change type")
	}

	if ctx.String("group-change-type") == "" {
		logrus.Fatalf("invalid group change type")
	}

	if ctx.String("ou-change-type") == "" {
		logrus.Fatalf("invalid ou change type")
	}

	if ok := checkPath(ctx.String("buzzword-dataset")); !ok {
		logrus.Fatalf("invalid buzzword dataset path: %s", ctx.String("buzzword-dataset"))
	}

	if ok := checkPath(ctx.String("department-dataset")); !ok {
		logrus.Fatalf("invalid department dataset path: %s", ctx.String("department-dataset"))
	}

	if ok := checkPath(ctx.String("first-name-dataset")); !ok {
		logrus.Fatalf("invalid first name dataset path: %s", ctx.String("first-name-dataset"))
	}

	if ok := checkPath(ctx.String("last-name-dataset")); !ok {
		logrus.Fatalf("invalid last name dataset path: %s", ctx.String("last-name-dataset"))
	}

	if ok := checkPath(ctx.String("groups-dataset")); !ok {
		logrus.Fatalf("invalid groups dataset path: %s", ctx.String("groups-dataset"))
	}

	generateConfig := &types.GenerateConfig{
		Users: ctx.Int("users"),
		Groups: ctx.Int("groups"),
		OUs: ctx.Int("ous"),
		OUDepth: ctx.Int("ou-depth"),
		UserChangeType: ctx.String("user-change-type"),
		GroupChangeType: ctx.String("group-change-type"),
		OUChangeType: ctx.String("ou-change-type"),
		Domain: *domainList,
		UserClasses: *userClassList,
		GroupClasses: *groupClassList,
		OUClasses: *ouClassList,
		BuzzwordDataset: ctx.String("buzzword-dataset"),
		DepartmentDataset: ctx.String("department-dataset"),
		FirstNameDataset: ctx.String("first-name-dataset"),
		LastNameDataset: ctx.String("last-name-dataset"),
		GroupsDataset: ctx.String("groups-dataset"),
	}

	return doGenerate(generateConfig)
}

func checkPath(path string) bool {
	if path == "" {
		return true // it's not invalid just not set
	}
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func parseDomain(domain string) (*[]string, error) {
	re := regexp.MustCompile("(?:[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?\\.)+[a-z0-9][a-z0-9-]{0,61}[a-z0-9]")
	if !re.Match([]byte(domain)) {
		return nil, errors.New(fmt.Sprintf("invalid domain %s, regex failed", domain))
	}

	domainList := strings.Split(domain, ".")
	if len(domainList) < 2 {
		return nil, errors.New(fmt.Sprintf("invalid domain %s, split resulted in < 2 segments", domain))
	}

	return &domainList, nil
}

func parseClassList(classes string) (*[]string, error) {
	classList := strings.Split(classes, ",")
	if len(classList) == 0 {
		return nil, errors.New(fmt.Sprintf("invalid class list: %v", classes))
	}
	return &classList, nil
}
