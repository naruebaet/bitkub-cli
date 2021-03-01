package util

import (
	"html/template"
	"log"
	"os"
)

// Mkdir : make dir
func Mkdir(path string, mode os.FileMode) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.MkdirAll(path, mode)
		if err != nil {
			log.Fatal(err)
		}
	}
}

// MkTemplate : make file from template file
func MkTemplate(filePath, templatePath string, templateValue interface{}) {
	// make file
	f, err := os.Create(filePath)
	if err != nil {
		log.Println("create file: ", err)
		return
	}
	defer f.Close()

	// parse template from templatePath
	t, err := template.ParseFiles(templatePath)
	if err != nil {
		log.Fatal(err)
		return
	}

	// add template value
	err = t.Execute(f, templateValue)
	if err != nil {
		log.Fatal(err)
	}
}
