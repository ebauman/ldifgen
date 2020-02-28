package names

import (
	"bufio"
	"io"
	"log"
	"math/rand"
	"os"
)

type NameGenerator struct {
	firstNames []string
	lastNames  []string
	departments []string
}

func (n NameGenerator) FirstName() string {
	index := rand.Int() % len(n.firstNames)

	return n.firstNames[index]
}

func (n NameGenerator) LastName() string {
	index := rand.Int() % len(n.lastNames)

	return n.lastNames[index]
}

func (n NameGenerator) Department() string {
	index := rand.Int() % len(n.departments)

	return n.departments[index]
}

func NewNameGenerator() (*NameGenerator, error) {
	n := &NameGenerator{}
	n.firstNames = make([]string, 0)
	n.lastNames = make([]string, 0)
	n.departments = make([]string, 0)
	err := n.initNames("datasets/first_names.txt", &n.firstNames)
	if err != nil {
		return nil, err
	}
	err = n.initNames("datasets/last_names.txt", &n.lastNames)
	if err != nil {
		return nil, err
	}
	err = n.initNames("datasets/departments.txt", &n.departments)
	if err != nil {
		return nil, err
	}

	return n, nil
}

func (n NameGenerator) initNames(file string, list *[]string) error {
	// read in the first and last name data sets
	f, err := os.Open(file)
	if err != nil {
		log.Fatalf("error opening %s: %v", file, err)
	}

	rdr := bufio.NewReader(f)

	for {
		line, _, err := rdr.ReadLine()
		if err == io.EOF {
			break
		}
		*list = append(*list, string(line))
	}

	return nil
}