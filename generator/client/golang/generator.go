package golang

import (
	"embed"
	"os"
	"path/filepath"
	"strings"

	"github.com/swiftcarrot/queryx/generator"
	"golang.org/x/tools/imports"
)

//go:embed templates
var templates embed.FS

func Run(g *generator.Generator, args []string) error {
	database := g.Schema.Databases[0]
	dir := database.Name

	if err := g.LoadTemplates(templates, database.Adapter); err != nil {
		return err
	}

	if err := g.Generate(); err != nil {
		return err
	}

	if err := goimports(dir); err != nil {
		return err
	}

	return nil
}

// run goimports for all go files in target directory
func goimports(dir string) error {
	return filepath.Walk(dir, func(target string, info os.FileInfo, err error) error {
		if strings.HasSuffix(target, ".go") {
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
		}
		return nil
	})
}
