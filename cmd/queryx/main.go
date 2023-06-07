package main

import (
	"fmt"
	"os"

	"github.com/swiftcarrot/queryx/cmd/queryx/action"
)

func main() {
	if err := action.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
