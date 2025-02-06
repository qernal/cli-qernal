package functions

import (
	"context"
	"log/slog"

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
			listResp, httpRes, err := qc.FunctionsAPI.ProjectsFunctionsList(ctx, projectID).Execute()
			if err != nil {
				resData, _ := client.ParseResponseData(httpRes)
				printer.Logger.Debug("unable to list projects, request failed", slog.String("error", err.Error()), slog.String("response", resData.(string)))
				return charm.RenderError("unable to list projects,  request failed with:", err)

			}

			if common.OutputFormat == "json" {
				printer.PrintResource(utils.FormatOutput(listResp.Data, common.OutputFormat))
				return nil
			}

			table := charm.RenderFuncTable(listResp.Data)
			printer.PrintResource(table)

			return nil
		},
	}

	return cmd
}
