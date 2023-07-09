package action

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/swiftcarrot/queryx/generator"
	"github.com/swiftcarrot/queryx/generator/client/golang"
	"github.com/swiftcarrot/queryx/generator/client/python"
	"github.com/swiftcarrot/queryx/generator/client/typescript"
)

var generateCmd = &cobra.Command{
	Use:     "generate",
	Aliases: []string{"g"},
	Short:   "Run generator defined in schema",
	RunE: func(cmd *cobra.Command, args []string) error {
		sch, err := newSchema()
		if err != nil {
			return err
		}

		database := sch.Databases[0]

		g := generator.NewGenerator(sch)

		for _, generator := range database.Generators {
			switch generator.Name {
			case "client-golang":
				if err := golang.Run(g, args); err != nil {
					return err
				}
			case "client-typescript":
				if err := typescript.Run(g, args); err != nil {
					return err
				}
			case "client-python":
				if err := python.Run(g, args); err != nil {
					return err
				}
			default:
				return fmt.Errorf("only supports generator.Name: %s , %s", "client-golang", "client-typescript")
			}
		}

		if err := g.Clean(); err != nil {
			return err
		}

		return nil
	},
}
