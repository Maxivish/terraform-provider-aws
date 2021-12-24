package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
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

func (c *ResourceCommand) Parse() string {

	pkg := "s3"
	resource := "aws_s3_zoo"

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "../internal/provider/provider.go", nil, parser.ParseComments)

	if err != nil {
		log.Println("Unable to parse provider.go")
	}
	var stack []ast.Node

	ast.Inspect(file, func(n ast.Node) bool {
		// s := ""
		// switch x := n.(type) {
		// case *ast.Ident:
		// 	s = x.Name

		// 	// Called recursively.
		// 	if s == "DataSourcesMap" {
		// 		ast.Print(fset, n)
		// 		return false
		// 	}
		// }
		// return true

		if n, ok := n.(*ast.Ident); ok {
			if n.Name == "DataSourcesMap" {

				dsMap := stack[len(stack)-1]

				var thing ast.Expr
				switch x := dsMap.(type) {
				case *ast.KeyValueExpr:
					thing = x.Value
				}

				var thng ast.CompositeLit

				switch x := thing.(type) {
				case *ast.CompositeLit:
					thng = *x
				}

				var pkgRes = make(map[string]int)

			out:
				for i := 0; i < len(thng.Elts); i++ {
					var top ast.KeyValueExpr
					switch x := thng.Elts[i].(type) {
					case *ast.KeyValueExpr:
						top = *x
					}

					var lit ast.BasicLit
					switch x := top.Key.(type) {
					case *ast.BasicLit:
						lit = *x
					}

					var call ast.CallExpr
					switch x := top.Value.(type) {
					case *ast.CallExpr:
						call = *x
					}

					var sel ast.SelectorExpr
					switch x := call.Fun.(type) {
					case *ast.SelectorExpr:
						sel = *x
					}

					var ide ast.Ident
					switch x := sel.X.(type) {
					case *ast.Ident:
						ide = *x
					}

					if ide.Name == pkg {
						pkgRes[lit.Value[1:len(lit.Value)-1]] = i
					}

					if ide.Name != pkg && len(pkgRes) > 0 {
						keys := make([]string, 0, len(pkgRes))
						for k := range pkgRes {
							keys = append(keys, k)
						}

						sort.Strings(keys)

						j := 0
						for _, key := range keys {
							//log.Println(key)
							j++
							if key > resource || j == len(pkgRes) {
								log.Printf("%d %d", j, len(pkgRes))
								log.Printf("%s %s %d", key, resource, pkgRes[key])
								break out
							}
						}
						//break out

					}

				}

				//log.Println(thng.Elts)

			}
		}

		// Manage the stack. Inspect calls a function like this:
		//   f(node)
		//   for each child {
		//      f(child) // and recursively for child's children
		//   }
		//   f(nil)
		if n == nil {
			// Done with node's children. Pop.
			stack = stack[:len(stack)-1]
		} else {
			// Push the current node for children.
			stack = append(stack, n)
		}

		return true
	})
	return ""
}
