package names

import (
	"bufio"
	"fmt"
	"github.com/rakyll/statik/fs"
	"io"
	"log"
	"math/rand"
	"os"
	"regexp"
)

type NameGenerator struct {
	firstNames  []string
	lastNames   []string
	departments []string
	buzzwords   []string
	groups      []string
}

func (n *NameGenerator) FirstName() string {
	index := rand.Intn(len(n.firstNames))

	return n.firstNames[index]
}

func (n *NameGenerator) LastName() string {
	index := rand.Intn(len(n.lastNames))

	return n.lastNames[index]
}

func (n *NameGenerator) Department() string {
	index := rand.Intn(len(n.departments))

	return n.departments[index]
}

func (n *NameGenerator) Group() string {
	buzzwordIndex := rand.Intn(len(n.buzzwords))
	groupIndex := rand.Intn(len(n.groups))

	return n.buzzwords[buzzwordIndex] + " " + n.groups[groupIndex]
}

func NewNameGenerator(firstNamePath string, lastNamePath string, departmentPath string, buzzwordPath string, groupPath string) (*NameGenerator, error) {
	n := &NameGenerator{}
	n.firstNames = make([]string, 0)
	n.lastNames = make([]string, 0)
	n.departments = make([]string, 0)
	n.buzzwords = make([]string, 0)
	n.groups = make([]string, 0)

	if firstNamePath == "" {
		firstNamePath = "statik://datasets/first_names.txt"
	}

	if lastNamePath == "" {
		lastNamePath = "statik://datasets/last_names.txt"
	}

	if departmentPath == "" {
		departmentPath = "statik://datasets/departments.txt"
	}

	if buzzwordPath == "" {
		buzzwordPath = "statik://datasets/buzzwords.txt"
	}

	if groupPath == "" {
		groupPath = "statik://datasets/groups.txt"
	}


	err := load(firstNamePath, &n.firstNames)
	if err != nil {
		return nil, err
	}
	err = load(lastNamePath, &n.lastNames)
	if err != nil {
		return nil, err
	}
	err = load(departmentPath, &n.departments)
	if err != nil {
		return nil, err
	}
	err = load(buzzwordPath, &n.buzzwords)
	if err != nil {
		return nil, err
	}
	err = load(groupPath, &n.groups)
	if err != nil {
		return nil, err
	}

	return n, nil
}

func load(file string, list *[]string) error {
	if file[0:9] == "statik://" {
		return loadStatik(file[8:], list) // 9 is not a mistake, we're stealing a / from statik://
	}

	if file[0:7] == "file://" {
		return loadFile(file[7:], list)
	}

	return fmt.Errorf("invalid file: %s", file)
}

func loadStatik(file string, list *[]string) error {
	statikFS, err := fs.New()
	if err != nil {
		return fmt.Errorf("error creating statik filesystem for file %s: %v", file, err)
	}

	f, err := statikFS.Open(file)
	if err != nil {
		return fmt.Errorf("error opening statik file %s: %v", file, err)
	}

	return readFile(f, list)
}

func readFile(file io.Reader, list *[]string) error {
	rdr := bufio.NewReader(file)
	re := regexp.MustCompile(`[\w ]+`)
	for {
		line, _, err := rdr.ReadLine()
		if err == io.EOF {
			break
		}
		if string(line) == "" {
			continue // get rid of empty lines if they occur
		}
		// get rid of non-words
		if len(re.Find(line)) < len(line) {
			// this means that the regex matched less than the total string
			continue
		}
		*list = append(*list, string(line))
	}

	return nil
}

func loadFile(file string, list *[]string) error {
	f, err := os.Open(file)
	if err != nil {
		log.Fatalf("error opening %s: %v", file, err)
	}

	return readFile(f, list)
}
