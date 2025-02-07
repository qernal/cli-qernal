package functions

import (
	"github.com/qernal/cli-qernal/charm"
	"github.com/qernal/cli-qernal/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	functionFile string
)
var FunctionCmd = &cobra.Command{
	Use:     "functions",
	Short:   "Manage your projects",
	Aliases: []string{"func", "fn", "function"},
	RunE: func(cmd *cobra.Command, args []string) error {
		err := cmd.Help()
		if err != nil {
			return err
		}
		return charm.RenderError("a valid subcommand is required")
	},
}

func init() {
	printer := utils.NewPrinter()
	FunctionCmd.AddCommand(NewListCmd(printer))
	FunctionCmd.AddCommand(NewCreateCmd(printer))

}
