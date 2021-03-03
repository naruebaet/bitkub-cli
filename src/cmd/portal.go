package cmd

import (
	"bitkub-cli/src/model"
	"bitkub-cli/src/pkg/util"
	"bitkub-cli/src/view/templ"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/alexflint/go-arg"
	"github.com/iancoleman/strcase"
)

func Init() {
	// declare bitkub cli
	args := model.Cli{}
	arg.MustParse(&args)
	process(args)
}

func process(cli model.Cli) {
	if cli.Init != "" {
		// do init project
		initProject(cli.Init)
	}

	if len(cli.Create) == 2 {
		switch cli.Create[0] {
		case "repository":
			createRepository(cli.Create[1])
		case "model":
			createModel(cli.Create[1])
		case "service":
			createService(cli.Create[1])
		}
	}
}

type MainFile struct {
	ProjectName string
}

func initProject(projectName string) {

	// create project folder
	util.Mkdir(projectName, os.ModePerm)

	// go version
	if out, err := exec.Command("go", "version").Output(); err != nil {
		log.Fatal("please check your go version it installed")
	} else {
		fmt.Print(string(out))
	}

	// swagger version
	if out, err := exec.Command("swag", "--version").Output(); err != nil {
		log.Fatal("please check your swag version it installed")
	} else {
		fmt.Print(string(out))
	}

	// go module init.
	cmd := exec.Command("go", "mod", "init", projectName)
	cmd.Dir = projectName
	if err := cmd.Run(); err != nil {
		log.Fatal("an error from go mod init "+projectName, err)
	}

	// just path list of project : singular style
	folder := []string{
		"/configs",
		"/model",
		"/repository",
		"/pkg/routing",
		"/pkg/driver",
		"/service",
		"/docs",
	}

	// loop make the project path
	for _, p := range folder {
		log.Println("folder created : " + projectName + p)
		util.Mkdir(projectName+p, os.ModePerm)
	}

	// create file in project
	goFile := []model.FileWithTemplate{
		{
			"main.go",
			templ.Main,
			MainFile{ProjectName: projectName},
		},
		{
			"pkg/routing/fiber.go",
			templ.Fiber,
			MainFile{ProjectName: projectName},
		},
		{
			"model/response.go",
			templ.Response,
			MainFile{ProjectName: projectName},
		},
		{
			"configs/config.yaml",
			templ.Config,
			MainFile{ProjectName: projectName},
		},
		{
			"pkg/driver/mongodb.go",
			templ.Mongodb,
			MainFile{ProjectName: projectName},
		},
		{
			"pkg/driver/mysql.go",
			templ.Mysql,
			MainFile{ProjectName: projectName},
		},
		{
			"pkg/driver/postgres.go",
			templ.Postgres,
			MainFile{ProjectName: projectName},
		},
		{
			"pkg/driver/redis.go",
			templ.Redis,
			MainFile{ProjectName: projectName},
		},
	}

	for _, srcFile := range goFile {
		log.Println("file created : " + projectName + "/" + srcFile.Filename)
		util.MkTemplateStr(projectName+"/"+srcFile.Filename, srcFile.Template, srcFile.Data)
	}

	// swagger init
	swag := exec.Command("swag", "init")
	swag.Dir = projectName
	if out, err := swag.Output(); err != nil {
		log.Fatal("an error from swag init "+projectName, err)
	} else {
		fmt.Print(string(out))
	}
}

type RepoFile struct {
	RepoName           string
	RepoNameLowerCamel string
	RepoNameUpperCamel string
}

func createRepository(repoName string) {
	if repoName != "" {
		log.Println("file created : ")
		util.MkTemplateStr("repository/"+repoName+".go", templ.Repository, RepoFile{
			RepoName:           repoName,
			RepoNameLowerCamel: strcase.ToLowerCamel(repoName),
			RepoNameUpperCamel: strcase.ToCamel(repoName),
		})
	}
}

func createModel(modelName string) {}

func createService(serviceName string) {}
