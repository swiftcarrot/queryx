package typescript

import (
	"embed"

	"github.com/swiftcarrot/queryx/generator"
)

//go:embed templates
var templates embed.FS

func Run(g *generator.Generator, args []string) error {
	database := g.Schema.Databases[0]

	if err := g.LoadTemplates(templates, database.Adapter); err != nil {
		return err
	}

	if err := g.Generate(); err != nil {
		return err
	}

	return nil
}
