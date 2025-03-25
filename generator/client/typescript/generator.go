package typescript

import (
	"embed"

	"github.com/swiftcarrot/queryx/generator"
	"github.com/swiftcarrot/queryx/schema"
)

//go:embed templates
var templates embed.FS

func Run(g *generator.Generator, generatorConfig *schema.Generator, args []string) error {
	database := g.Schema.Databases[0]

	if err := g.LoadTemplates(templates, database.Adapter); err != nil {
		return err
	}

	if err := g.Generate(transform, ""); err != nil {
		return err
	}

	return nil
}

func transform(b []byte) []byte {
	return b
}
