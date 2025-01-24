package projects

import (
	"context"
	"fmt"

	"github.com/qernal/cli-qernal/charm"
	"github.com/qernal/cli-qernal/commands/auth"
	"github.com/qernal/cli-qernal/pkg/client"
	"github.com/qernal/cli-qernal/pkg/common"
	"github.com/qernal/cli-qernal/pkg/helpers"
	"github.com/qernal/cli-qernal/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	pageSize int32
)

func NewProjectsListCmd(printer *utils.Printer) *cobra.Command {
	cmd := &cobra.Command{
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

			maxResults, _ := cmd.Flags().GetInt32("max")

			allProjects, err := helpers.PaginateProjects(printer, ctx, &qc, maxResults)
			if err != nil {
				return charm.RenderError("unable to list projects", err)
			}
			if maxResults > 0 && len(allProjects) > int(maxResults) {
				allProjects = allProjects[:maxResults]
			}

			if common.OutputFormat == "json" {
				fmt.Println(utils.FormatOutput(allProjects, common.OutputFormat))
				return nil
			}
			table := charm.RenderProjectTable(allProjects)
			fmt.Println(table)

			return nil
		},
	}
	return cmd
}
