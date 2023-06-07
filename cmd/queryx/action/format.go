package action

import (
	"fmt"
	"io"
	"os"

	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/spf13/cobra"
)

var formatCmd = &cobra.Command{
	Use:     "format",
	Aliases: []string{"fmt"},
	Short:   "Format current schema file",
	RunE: func(cmd *cobra.Command, args []string) error {
		f, err := os.Open(schemaFile)
		if err != nil {
			return err
		}

		inSrc, err := io.ReadAll(f)
		if err != nil {
			return err
		}

		// TODO: check syntax
		outSrc := hclwrite.Format(inSrc)

		// TODO: add overwrite option
		if err := os.WriteFile(schemaFile, outSrc, 0644); err != nil {
			return err
		}

		fmt.Println("Format", schemaFile)

		return nil
	},
}
