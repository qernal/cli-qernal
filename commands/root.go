package commands

import (
	"os"

	"github.com/qernal/cli-qernal/commands/auth"
	"github.com/qernal/cli-qernal/commands/encrypt"
	"github.com/qernal/cli-qernal/pkg/common"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "qernal",
	Short: "CLI for interacting with qernal.com",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return cmd.Help()
		}
		return nil
	},
}

func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	RootCmd.AddCommand(auth.AuthCmd)
	RootCmd.AddCommand(encrypt.EncryptCmd)
	RootCmd.PersistentFlags().StringVarP(&common.OutputFormat, "output", "o", "text", "output format (json,text)")
}
