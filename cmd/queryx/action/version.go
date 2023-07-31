package action

import (
	"fmt"

	"github.com/spf13/cobra"
)

var Version = "test-1"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Prints current version",
	Run: func(cmd *cobra.Command, args []string) {
		RootCmd.Println(fmt.Sprintf("queryx %s", Version))
	},
}
