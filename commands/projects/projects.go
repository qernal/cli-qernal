package projects

import (
	"github.com/qernal/cli-qernal/charm"
	"github.com/qernal/cli-qernal/pkg/utils"
	"github.com/spf13/cobra"
)

var ProjectsCmd = &cobra.Command{
	Use:     "projects",
	Short:   "Manage your projects",
	Aliases: []string{"project"},
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
	ProjectsCmd.AddCommand(ProjectsListCmd)
	ProjectsCmd.AddCommand(NewCreateCmd(printer))
	ProjectsCmd.AddCommand(NewDeleteCmd(printer))
	ProjectsCmd.AddCommand(NewUpdateCmd(printer))
}
