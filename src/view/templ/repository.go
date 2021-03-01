package templ

const Repository = `
package repository

type {{.RepoNameUpperCamel}}Repository interface {
	Todo(val string) string
}

func New{{.RepoNameUpperCamel}}Repository() ({{.RepoNameUpperCamel}}Repository, error) {
	rp := {{.RepoNameLowerCamel}}Repository{}
	return &rp, nil
}

type {{.RepoNameLowerCamel}}Repository struct{}

func (r *{{.RepoNameLowerCamel}}Repository) Todo(val string) string {
	return val 
}
`

