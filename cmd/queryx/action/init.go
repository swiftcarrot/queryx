package action

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var sampleSchema = `database "db" {
  adapter = "sqlite"

  config "development" {
    url = "sqlite:blog_development.sqlite3"
  }

  model "Post" {
    column "title" {
      type = string
    }
    column "content" {
      type = text
    }
  }
}
`

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Creates a sample schema.hcl",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Created schema.hcl")
		file, err := os.Create("schema.hcl")
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = file.WriteString(sampleSchema)
		if err != nil {
			return err
		}
		return nil
	},
}
