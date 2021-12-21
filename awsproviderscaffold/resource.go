package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type TemplateType int64

const (
	Resource TemplateType = iota
	Test
)

type ResourceCommand struct {
}

func (c *ResourceCommand) Run(args []string) int {
	//TODO Flag to include tag boilerplate
	pkg := ""
	name := ""

	if len(args) > 0 {
		pkg = args[0]
		name = args[1]
	}

	config := map[string]string{
		"pkg":  pkg,
		"name": name,
	}

	createResourceFiles(config)

	return 0
}

func createResourceFiles(config map[string]string) error {
	execTemplate("resource", Resource, config)
	execTemplate("resource", Test, config)
	return nil
}

func execTemplate(tmpl string, tmplType TemplateType, config map[string]string) error {
	pwd, _ := os.Getwd()

	funcMap := template.FuncMap{
		"Title": strings.Title,
	}

	suffix := ""

	if tmplType == Test {
		suffix = "_test"
	}

	tmpl = fmt.Sprintf("%s%s.tmpl", tmpl, suffix)

	t, err := template.New(tmpl).Funcs(funcMap).ParseFiles(filepath.Join(pwd, fmt.Sprintf("tmpl/%s", tmpl)))

	if err != nil {
		log.Print(err)
		return err
	}

	f, err := os.Create(filepath.Join(pwd, "..", "internal", "service", config["pkg"], config["name"]+suffix+".go"))
	if err != nil {
		log.Println("create file: ", err)
		return err
	}

	err = t.Execute(f, config)
	if err != nil {
		log.Print("execute: ", err)
		return err
	}
	f.Close()

	return nil
}

func (c *ResourceCommand) Synopsis() string {
	return "Creates a templated resource belonging to an AWS service package."
}

func (c *ResourceCommand) Help() string {
	return "halp"
}
