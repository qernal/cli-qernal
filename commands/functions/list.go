package functions

import (
	"context"

	"github.com/qernal/cli-qernal/charm"
	"github.com/qernal/cli-qernal/commands/auth"
	"github.com/qernal/cli-qernal/pkg/client"
	"github.com/qernal/cli-qernal/pkg/common"
	"github.com/qernal/cli-qernal/pkg/helpers"
	"github.com/qernal/cli-qernal/pkg/utils"
	"github.com/spf13/cobra"
)

func NewListCmd(printer *utils.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Example: "qernal func list --project <project name> ",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return helpers.ValidateProjectFlags(cmd)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()

			token, err := auth.GetQernalToken()
			if err != nil {
				return charm.RenderError("error creating qernal client", err)
			}

			qc, err := client.New(ctx, nil, nil, token)
			if err != nil {
				return charm.RenderError("error creating qernal client", err)
			}

			projectID, err := helpers.GetProjectID(cmd, &qc)
			if err != nil {
				return err
			}
			maxResults, _ := cmd.Flags().GetInt32("max")

			functions, err := helpers.PaginateFunctions(printer, ctx, &qc, maxResults, projectID)
			if err != nil {
				return charm.RenderError("unable to list organisations", err)
			}
			if maxResults > 0 && len(functions) > int(maxResults) {
				functions = functions[:maxResults]
			}

			if common.OutputFormat == "json" {
				printer.PrintResource(utils.FormatOutput(functions, common.OutputFormat))
				return nil
			}

			table := charm.RenderFuncTable(functions)
			printer.PrintResource(table)

			return nil
		},
	}

	return cmd
}
