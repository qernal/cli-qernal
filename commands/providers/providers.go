package providers

import (
	"errors"

	"github.com/qernal/cli-qernal/pkg/utils"
	"github.com/spf13/cobra"
)

var ProvidersCmd = &cobra.Command{
	Use:          "providers",
	Short:        "View available qernal providers",
	Aliases:      []string{"prviders", "provider"},
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {

		err := cmd.Help()
		if err != nil {
			return err
		}
		return errors.New("a valid subcommand is required")
	},
}

func init() {
	printer := utils.NewPrinter()
	ProvidersCmd.AddCommand(NewListCmd(printer))
}
