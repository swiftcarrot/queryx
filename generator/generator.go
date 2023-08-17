package generator

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/swiftcarrot/queryx/inflect"
	"github.com/swiftcarrot/queryx/schema"
)

type Generator struct {
	Adapter   string
	Templates []*template.Template
	Schema    *schema.Schema
	created   []string
}

func NewGenerator(schema *schema.Schema) *Generator {
	return &Generator{
		Schema:  schema,
		created: make([]string, 0),
	}
}

func (g *Generator) LoadTemplates(src embed.FS, adapter string) error {
	// TODO: clean templates implicitly
	g.Templates = []*template.Template{}

	if err := fs.WalkDir(src, "templates", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("error walking directory %s: %w", path, err)
		}
		if d.IsDir() {
			return nil
		}

		templateName := strings.TrimPrefix(path, "templates")
		buf, err := src.ReadFile(path)
		if err != nil {
			return fmt.Errorf("error reading file '%s': %w", path, err)
		}
		templateName = strings.TrimSuffix(templateName, "tmpl")

		ss := strings.Split(templateName, ".")
		if len(ss) > 2 {
			if stringInSlice(ss[len(ss)-2], []string{"postgresql", "mysql", "sqlite"}) {
				if ss[len(ss)-2] == adapter {
					templateName = strings.Join(ss[:len(ss)-2], ".") + "." + ss[len(ss)-1]
				} else {
					return nil
				}
			}
		}

		tmpl := template.New(templateName).Funcs(inflect.TemplateFunctions)
		tmpl, err = tmpl.Parse(string(buf))
		if err != nil {
			return fmt.Errorf("parsing template '%s': %w", path, err)
		}

		g.Templates = append(g.Templates, tmpl)

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (g *Generator) CreateFile(f string, tpl *template.Template, data interface{}) error {
	fmt.Println("Created", f)

	dir := filepath.Dir(f)
	if dir != "" {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return err
		}
	}

	wr, err := os.Create(f)
	if err != nil {
		return err
	}

	if err := tpl.Execute(wr, data); err != nil {
		return err
	}

	return nil
}

func stringInSlice(str string, list []string) bool {
	for _, v := range list {
		if v == str {
			return true
		}
	}
	return false
}

func (g *Generator) Generate() error {
	database := g.Schema.Databases[0]
	dir := database.Name

	for _, tpl := range g.Templates {
		name := tpl.Name()

		if strings.Contains(name, "[model]") {
			for _, model := range database.Models {
				n := strings.ReplaceAll(name, "[model]", inflect.Snake(model.Name))
				f := filepath.Join(dir, n)
				g.created = append(g.created, f)

				data := map[string]interface{}{
					"packageName": dir,
					"client":      database,
					"model":       model,
				}
				if err := g.CreateFile(f, tpl, data); err != nil {
					return err
				}
			}
		} else {
			f := filepath.Join(dir, name)
			g.created = append(g.created, f)
			data := map[string]interface{}{
				"packageName": dir,
				"client":      database,
			}
			if err := g.CreateFile(f, tpl, data); err != nil {
				return err
			}
		}
	}

	return nil
}

func (g *Generator) Clean() error {
	database := g.Schema.Databases[0]
	dir := database.Name

	deleted := []string{}
	files, err := readDir(dir)

	if err != nil {
		return err
	}
	for _, f := range files {
		if !stringInSlice(f, g.created) {
			deleted = append(deleted, f)
		}
	}
	for _, f := range deleted {
		os.Remove(f)
		fmt.Println("Deleted", f)
	}
	return nil
}

func readDir(dir string) ([]string, error) {
	files := []string{}
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && info.Name() == "migrations" {
			return filepath.SkipDir
		}
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil
}
