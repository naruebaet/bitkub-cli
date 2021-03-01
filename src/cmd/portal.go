package cmd

import (
	"bitkub-cli/src/model"
	"bitkub-cli/src/pkg/util"
	"fmt"
	"github.com/alexflint/go-arg"
	"log"
	"os"
	"os/exec"
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
		"/config",
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
			"src/view/main.template",
			MainFile{ProjectName: projectName},
		},
		{
			"pkg/routing/fiber.go",
			"src/view/fiber.template",
			MainFile{ProjectName: projectName},
		},
		{
			"model/response.go",
			"src/view/response.template",
			MainFile{ProjectName: projectName},
		},
		{
			"config/config.yaml",
			"src/view/config.template",
			MainFile{ProjectName: projectName},
		},
		{
			"pkg/driver/mongodb.go",
			"src/view/driver/mongodb.template",
			MainFile{ProjectName: projectName},
		},
		{
			"pkg/driver/mysql.go",
			"src/view/driver/mysql.template",
			MainFile{ProjectName: projectName},
		},
		{
			"pkg/driver/postgres.go",
			"src/view/driver/postgres.template",
			MainFile{ProjectName: projectName},
		},
		{
			"pkg/driver/redis.go",
			"src/view/driver/redis.template",
			MainFile{ProjectName: projectName},
		},
	}

	for _, srcFile := range goFile {
		log.Println("file created : " + projectName + "/" + srcFile.Filename)
		util.MkTemplate(projectName+"/"+srcFile.Filename, srcFile.Template, srcFile.Data)
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
