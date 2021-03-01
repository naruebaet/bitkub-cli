package model

type Cli struct {
	Init   string   `arg:"-i,--init" help:"init repository pattern project with the bitkub cli naja"`
	Create []string `arg:"-c,--create" help:"create a repository, service, model"`
}

type FileWithTemplate struct {
	Filename string
	Template string
	Data     interface{}
}
