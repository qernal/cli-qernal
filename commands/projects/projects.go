package projects

import (
	"github.com/qernal/cli-qernal/charm"
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
	ProjectsCmd.AddCommand(ProjectsListCmd)
	ProjectsCmd.AddCommand(CreateCmd)
	ProjectsCmd.AddCommand(DeleteCmd)
	ProjectsCmd.AddCommand(updateCmd)
}
