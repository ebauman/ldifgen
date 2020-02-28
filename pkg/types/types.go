package types

type User struct {
	Id          string
	GivenName   string
	CommonName  string
	Surname     string
	Manager     string
	Secretary   string
	Title       string
	Description string
	OU          string
}

type Group struct {
	Id         string
	CommonName string
	Members    []*User
}

type OrganizationalUnit struct {
	Id                 string
	CommonName         string
	OrganizationalUnit *OrganizationalUnit
}

func (u User) UID() string {
	return u.GivenName + "." + u.Surname
}

func (u User) recurseOUChain(ou *OrganizationalUnit, path string) string {
	path += "ous=" + ou.CommonName
	if ou.OrganizationalUnit == nil {
		return path
	} else {
		path += ","
		return u.recurseOUChain(ou.OrganizationalUnit, path)
	}
}
