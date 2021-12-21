package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type ResourceCommand struct {
}

func (c *ResourceCommand) Run(args []string) int {
	pkg := ""
	name := ""
	//dir := ""

	pwd, _ := os.Getwd()

	filepath.Join(pwd, "tmpl/resource.tmpl")

	if len(args) > 0 {
		pkg = args[0]
		name = args[1]
		//dir = args[2]
	}

	funcMap := template.FuncMap{
		"ToTitle": strings.Title,
	}

	t, err := template.New("resource.tmpl").Funcs(funcMap).ParseFiles(filepath.Join(pwd, "tmpl/resource.tmpl"))

	if err != nil {
		log.Print(err)
		return 1
	}

	f, err := os.Create(filepath.Join(pwd, "..", "internal", "service", pkg, name+".go"))
	if err != nil {
		log.Println("create file: ", err)
		return 1
	}

	config := map[string]string{
		"pkg":  pkg,
		"name": name,
	}

	err = t.Execute(f, config)
	if err != nil {
		log.Print("execute: ", err)
		return 1
	}
	f.Close()

	return 0
}

func (c *ResourceCommand) Synopsis() string {
	return "Creates a templated resource belonging to an AWS service package."
}

func (c *ResourceCommand) Help() string {
	return "halp"
}
