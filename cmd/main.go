package main

import (
	"os"

	"github.com/qernal/cli-qernal/commands"
)

func main() {
	if err := commands.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
