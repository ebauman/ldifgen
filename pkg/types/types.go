package types

type User struct {
	Id                string
	GivenName         string
	DistinguishedName string
	CommonName        string
	Surname           string
	Manager           string
	Secretary         string
	Title             string
	Description       string
	OU                string
}

type Group struct {
	Id                 string
	CommonName         string
	OrganizationalUnit string
	DistinguishedName  string
	Members            []string
}

func (u *User) UID() string {
	return u.GivenName + "." + u.Surname
}
