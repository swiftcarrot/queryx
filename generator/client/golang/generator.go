package golang

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/swiftcarrot/queryx/generator"
	"github.com/swiftcarrot/queryx/schema"
	"golang.org/x/sync/errgroup"
	"golang.org/x/tools/imports"
)

//go:embed templates
var templates embed.FS

func Run(generatorConfig *schema.Generator, schema *schema.Schema, args []string) error {
	g := &generator.Generator{}
	database := schema.Databases[0]
	dbName := database.Name

	if err := g.LoadTemplates(templates, database.Adapter); err != nil {
		return err
	}

	// remove test files from templates by default
	if !generatorConfig.Test {
		templates := []*template.Template{}
		for _, tpl := range g.Templates {
			name := tpl.Name()
			if !strings.HasSuffix(name, "_test.go") {
				templates = append(templates, tpl)
			}
		}
		g.Templates = templates
	}

	// remove unused types in templates
	templates := []*template.Template{}
	typs := typesFromSchema(schema)
	for _, tpl := range g.Templates {
		name := tpl.Name()
		name = strings.TrimPrefix(name, "/queryx/")
		name = strings.TrimSuffix(name, ".go")
		name = strings.TrimSuffix(name, "_column")
		if b, ok := typs[name]; !ok || b {
			templates = append(templates, tpl)
		}
	}
	g.Templates = templates

	if err := g.Generate(schema); err != nil {
		return err
	}

	fmt.Println("Running goimports...")
	if err := goimports(dbName); err != nil {
		return err
	}

	return nil
}

func typesFromSchema(sch *schema.Schema) map[string]bool {
	m := map[string]bool{}
	typs := []string{"string", "text", "boolean",
		"date", "time", "datetime", "float",
		"integer", "bigint", "json", "uuid"}
	for _, typ := range typs {
		m[typ] = false
	}

	for _, database := range sch.Databases {
		for _, model := range database.Models {
			for _, column := range model.Columns {
				typ := column.Type
				if typ == "jsonb" {
					typ = "json"
				}
				m[typ] = true
			}
		}
	}

	return m
}

// run goimports for all go files in target directory
func goimports(dir string) error {
	g := new(errgroup.Group)
	g.SetLimit(20)

	targets := []string{}
	err := filepath.Walk(dir, func(target string, info os.FileInfo, err error) error {
		if strings.HasSuffix(target, ".go") {
			targets = append(targets, target)
		}
		return nil
	})
	if err != nil {
		return err
	}

	for _, target := range targets {
		func(target string) {
			g.Go(func() error {
				content, err := os.ReadFile(target)
				if err != nil {
					return err
				}
				src, err := imports.Process(target, content, nil)
				if err != nil {
					return err
				}
				if err := os.WriteFile(target, src, 0644); err != nil {
					return err
				}

				return nil
			})
		}(target)
	}

	if err := g.Wait(); err != nil {
		return err
	}

	return nil
}
