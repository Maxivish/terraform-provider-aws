package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"sort"
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

	c.Parse()

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

	pkgDir := filepath.Join(pwd, "..", "internal", "service", config["pkg"])

	_ = os.MkdirAll(pkgDir, os.ModePerm)

	f, err := os.Create(filepath.Join(pkgDir, config["name"]+suffix+".go"))
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

Creates a templated resource belonging to an AWS service package.
`

	f := c.Flags()

	f.VisitAll(func(f *flag.Flag) {
		helpText += fmt.Sprintf("\n %s: %s", f.Name, f.Usage)
	})

	return helpText
}

func Entry() (a *ast.KeyValueExpr) {

	thing := &ast.KeyValueExpr{}

	thing.Key = &ast.BasicLit{
		Kind:  token.STRING,
		Value: "\"aws_acm_test\"",
	}

	thing.Value = &ast.CallExpr{
		Fun: &ast.SelectorExpr{
			X: &ast.Ident{
				Name: "acm",
			},
			Sel: &ast.Ident{
				Name: "DataSourceCertificate",
			},
		},
	}
	return thing
}

func (c *ResourceCommand) Parse() string {
	// find target map datasource/resource
	// find target pkg position
	//   if exists add

	pkg := "s3"
	resource := "aws_s3_zoo"

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "../internal/provider/provider.go", nil, parser.ParseComments)

	if err != nil {
		log.Println("Unable to parse provider.go")
	}
	var stack []ast.Node

	ast.Inspect(file, func(n ast.Node) bool {
		if n, ok := n.(*ast.Ident); ok {
			if n.Name == "DataSourcesMap" {
				dsMap := stack[len(stack)-1]

				keyValueExpr := dsMap.(*ast.KeyValueExpr)
				resourceMap := keyValueExpr.Value.(*ast.CompositeLit)

				var pkgRes = make(map[string]int)

			out:
				for i := 0; i < len(resourceMap.Elts); i++ {

					keyValueExpr = resourceMap.Elts[i].(*ast.KeyValueExpr)

					key := keyValueExpr.Key.(*ast.BasicLit)
					value := keyValueExpr.Value.(*ast.CallExpr)

					selectorExpr := value.Fun.(*ast.SelectorExpr)
					resourcePackage := selectorExpr.X.(*ast.Ident)

					if resourcePackage.Name == pkg {
						pkgRes[key.Value[1:len(key.Value)-1]] = i
					}

					if resourcePackage.Name != pkg && len(pkgRes) > 0 {
						keys := make([]string, 0, len(pkgRes))
						for k := range pkgRes {
							keys = append(keys, k)
						}

						sort.Strings(keys)

						j := 0
						for _, key := range keys {
							j++
							if key > resource || j == len(pkgRes) {

								resourceMap.Elts = insert(resourceMap.Elts, pkgRes[key], Entry())
								break out
							}
						}
					}
				}
			}
		}

		if n == nil {
			stack = stack[:len(stack)-1]
		} else {
			stack = append(stack, n)
		}

		return true
	})
	printer.Fprint(os.Stdout, fset, file)

	return ""
}

func insert(a []ast.Expr, index int, value ast.Expr) []ast.Expr {
	if len(a) == index { // nil or empty slice or after last element
		return append(a, value)
	}
	a = append(a[:index+1], a[index:]...) // index < len(a)
	a[index] = value
	return a
}
