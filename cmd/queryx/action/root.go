package action

import (
	"os"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/spf13/cobra"
	"github.com/swiftcarrot/queryx/schema"
)

var schemaFile string

var RootCmd = &cobra.Command{
	Use: "queryx",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func newSchema() (*schema.Schema, error) {
	content, err := os.ReadFile(schemaFile)
	if err != nil {
		return nil, err
	}

	file, diagnostics := hclsyntax.ParseConfig(content, schemaFile, hcl.Pos{Line: 1, Column: 1, Byte: 0})
	if diagnostics != nil && diagnostics.HasErrors() {
		return nil, diagnostics.Errs()[0]
	}
	sch, err := schema.Parse(file.Body)
	if err != nil {
		return nil, err
	}

	return sch, nil
}

func init() {
	RootCmd.PersistentFlags().StringVar(&schemaFile, "schema", "schema.hcl", "queryx schema file")

	RootCmd.AddCommand(versionCmd)
	RootCmd.AddCommand(formatCmd)
	RootCmd.AddCommand(generateCmd)

	RootCmd.AddCommand(dbCreateCmd)
	RootCmd.AddCommand(dbDropCmd)
	RootCmd.AddCommand(dbMigrateCmd)
	RootCmd.AddCommand(dbMigrateGenerateCmd)
	RootCmd.AddCommand(dbRollbackCmd)
	RootCmd.AddCommand(dbMigrateStatusCmd)
	RootCmd.AddCommand(dbVersionCmd)
}
