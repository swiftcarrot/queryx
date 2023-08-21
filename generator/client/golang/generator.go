package golang

import (
	"embed"
	"go/format"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/swiftcarrot/queryx/generator"
	"github.com/swiftcarrot/queryx/schema"
	"golang.org/x/mod/modfile"
)

//go:embed templates
var templates embed.FS

func Run(g *generator.Generator, generatorConfig *schema.Generator, args []string) error {
	schema := g.Schema
	database := schema.Databases[0]

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

	// TODO: wrap this in a function
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	roots := findModuleRoot(cwd)
	f := filepath.Join(roots, "go.mod")
	content, err := os.ReadFile(f)
	if err != nil {
		return err
	}
	mf, err := modfile.Parse("go.mod", content, nil)
	if err != nil {
		return err
	}
	rel, err := filepath.Rel(roots, cwd)
	if err != nil {
		return err
	}
	goModPath := path.Join(mf.Module.Mod.Path, rel)

	if err := g.Generate(transform, goModPath); err != nil {
		return err
	}

	return nil
}

func transform(b []byte) []byte {
	b, err := format.Source(b)
	if err != nil {
		panic(err)
	}
	return b
}

func typesFromSchema(sch *schema.Schema) map[string]bool {
	m := map[string]bool{}
	// TODO: move to schema package
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

// findModuleRoot finds the module root by looking for go.mod file in the current directory and its parents.
// This function is copied from https://github.com/golang/go/blob/master/src/cmd/go/internal/modload/init.go
func findModuleRoot(dir string) (roots string) {
	if dir == "" {
		// TODO: add go mod init in docs
		panic("dir not set") // TODO: improve this error message
	}
	dir = filepath.Clean(dir)

	for {
		if fi, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil && !fi.IsDir() {
			return dir
		}
		d := filepath.Dir(dir)
		if d == dir {
			break
		}
		dir = d
	}
	return ""
}
