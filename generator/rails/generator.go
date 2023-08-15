package rails

import (
	"embed"
	"path"
	"strings"

	"github.com/swiftcarrot/queryx/generator"
	"github.com/swiftcarrot/queryx/inflect"
	"github.com/swiftcarrot/queryx/schema"
)

//go:embed templates
var templates embed.FS

func Run(g *generator.Generator, generatorConfig *schema.Generator, args []string) error {
	database := g.Schema.Databases[0]

	if err := g.LoadTemplates(templates, database.Adapter); err != nil {
		return err
	}

	if err := generate(g); err != nil {
		return err
	}

	return nil
}

func generate(g *generator.Generator) error {
	database := g.Schema.Databases[0]
	dir := "./"

	for _, tpl := range g.Templates {
		name := tpl.Name()

		if strings.Contains(name, "[model]") {
			for _, model := range database.Models {
				n := strings.ReplaceAll(name, "[model]", inflect.Snake(model.Name))
				f := path.Join(dir, n)

				data := map[string]interface{}{
					"client": database,
					"model":  model,
				}
				if err := g.CreateFile(f, tpl, data); err != nil {
					return err
				}
			}
		} else {
			f := path.Join(dir, name)
			data := map[string]interface{}{
				"client": database,
			}
			if err := g.CreateFile(f, tpl, data); err != nil {
				return err
			}
		}
	}

	return nil
}
