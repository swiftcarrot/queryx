package action

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/swiftcarrot/queryx/generator/client/golang"
)

var generateCmd = &cobra.Command{
	Use:     "generate",
	Aliases: []string{"g"},
	Short:   "Run generator in schema",
	RunE: func(cmd *cobra.Command, args []string) error {
		sch, err := newSchema()
		if err != nil {
			return err
		}

		database := sch.Databases[0]

		for _, generator := range database.Generators {
			switch generator.Name {
			case "client-golang":
				if err := golang.Run(sch, args); err != nil {
					return err
				}
			case "client-typescript":
			// TODO: fix typescript
			// gen := &typescript.Generator{}
			// if err := gen.Generate(sch, args); err != nil {
			// 	return err
			// }
			default:
				return fmt.Errorf("only supports generator.Name: %s , %s", "client-golang", "client-typescript")
			}
		}

		return nil
	},
}
