package main

import (
	"flag"
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
	flagPkg  string
	flagName string
}

func (c *ResourceCommand) Run(args []string) int {
	//TODO Flag to include tag boilerplate
	//TODO Validate pkg name exists.

	c.Flags().Parse(args)

	if c.flagPkg == "" {
		log.Println("Package name must be specified: ")
		return 1
	}

	if c.flagName == "" {
		log.Println("Resource name must be specified: ")
		return 1
	}

	config := map[string]string{
		"pkg":  c.flagPkg,
		"name": c.flagName,
	}

	createResourceFiles(config)

	return 0
}

func (c *ResourceCommand) Flags() *flag.FlagSet {
	f := flag.NewFlagSet("resource", flag.ContinueOnError)

	f.StringVar(&c.flagPkg, "pkg", "", "Name of AWS Go SDK service package name which contains the API endpoints to be used by the resource.")
	f.StringVar(&c.flagName, "name", "", "Name of the Terraform Resource you wish to create.")

	return f
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
	helpText := `
Usage: awsproviderscaffold resource [options] 
Creates a templated resource belonging to an AWS service package.` + c.Flags().Usage

	return strings.TrimSpace(helpText)
}
