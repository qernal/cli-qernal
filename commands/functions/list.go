package functions

import (
	"bytes"
	"context"
	"log/slog"

	"github.com/qernal/cli-qernal/charm"
	"github.com/qernal/cli-qernal/commands/auth"
	"github.com/qernal/cli-qernal/pkg/client"
	"github.com/qernal/cli-qernal/pkg/common"
	"github.com/qernal/cli-qernal/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	projectID string
)

func NewListCmd(printer *utils.Printer) *cobra.Command {

	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Example: "qernaal func list --project <project name> ",
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

			project, err := qc.GetProjectByName(projectID)
			if err != nil {
				return charm.RenderError("x", err)
			}
			listResp, httpRes, err := qc.FunctionsAPI.ProjectsFunctionsList(ctx, project.Id).Execute()
			if err != nil {
				resData, _ := client.ParseResponseData(httpRes)
				printer.Logger.Debug("unable to list projects, request failed", slog.String("error", err.Error()), slog.String("response", resData.(string)))
				return charm.RenderError("unable to list projects,  request failed with:", err)

			}

			if common.OutputFormat == "json" {
				printer.PrintResource(utils.FormatOutput(listResp.Data, common.OutputFormat))
				return nil
			}

			var buf bytes.Buffer
			table := charm.RenderFuncTable(&buf, listResp.Data)
			printer.PrintResourceR(&table)

			return nil
		},
	}
	cmd.Flags().StringVarP(&projectID, "project", "p", "", "name of the project")
	cmd.MarkFlagRequired("project")
	return cmd
}
