package typescript

import (
	"embed"

	"github.com/swiftcarrot/queryx/generator"
	"github.com/swiftcarrot/queryx/schema"
)

//go:embed templates
var templates embed.FS

func Run(schema *schema.Schema, args []string) error {
	g := &generator.Generator{}
	database := schema.Databases[0]
	// dbName := database.Name // TODO: add output config

	if err := g.LoadTemplates(templates, database.Adapter); err != nil {
		return err
	}

	if err := g.Generate(schema); err != nil {
		return err
	}

	return nil
}
