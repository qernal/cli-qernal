package projects

import (
	"context"
	"fmt"

	"github.com/qernal/cli-qernal/charm"
	"github.com/qernal/cli-qernal/commands/auth"
	"github.com/qernal/cli-qernal/pkg/client"
	"github.com/qernal/cli-qernal/pkg/common"
	"github.com/qernal/cli-qernal/pkg/utils"
	"github.com/spf13/cobra"
)

var ProjectsListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls", "l"},
	Short:   "list your qernal projects",
	Example: "qernal projects list",
	RunE: func(cmd *cobra.Command, args []string) error {

		token, err := auth.GetQernalToken()
		if err != nil {
			return charm.RenderError("unable to retrieive qernal token, run qernal auth login if you haven't")
		}
		ctx := context.Background()
		qc, err := client.New(ctx, nil, nil, token)
		if err != nil {
			return charm.RenderError("", err)
		}
		projectsResp, _, err := qc.ProjectsAPI.ProjectsList(ctx).Execute()
		if err != nil {
			if err := charm.RenderError("unable to list projects, request failed with:", err); err != nil {
				return fmt.Errorf("failed to render error: %w", err)
			}
			return err
		}

		if common.OutputFormat == "json" {
			fmt.Println(utils.FormatOutput(projectsResp.Data, common.OutputFormat))
			return nil
		}
		table := charm.RenderProjectTable(projectsResp.Data)
		fmt.Println(table)

		return nil
	},
}
