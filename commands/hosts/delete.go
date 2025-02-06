package hosts

import (
	"context"
	"errors"
	"log/slog"

	"github.com/qernal/cli-qernal/charm"
	"github.com/qernal/cli-qernal/commands/auth"
	"github.com/qernal/cli-qernal/pkg/client"
	"github.com/qernal/cli-qernal/pkg/common"
	"github.com/qernal/cli-qernal/pkg/utils"
	"github.com/spf13/cobra"
)

func NewDeleteCmd(printer *utils.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete",
		Aliases: []string{"rm"},
		Example: "qernal host delete --name <host name>",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				return charm.RenderError("No arguments expected")
			}

			ctx := context.Background()
			token, err := auth.GetQernalToken()
			if err != nil {
				return charm.RenderError("unable to retreive qernal token, run qernal auth login if you haven't")

			}

			qc, err := client.New(ctx, nil, nil, token)
			if err != nil {
				return charm.RenderError("error creating qernal client", err)
			}

			projectName, _ := cmd.Flags().GetString("project")
			hostName, _ := cmd.Flags().GetString("name")

			project, err := qc.GetProjectByName(projectName)
			if err != nil {
				return charm.RenderError("x", err)
			}

			DeleteResp, httpRes, err := qc.HostsAPI.ProjectsHostsDelete(ctx, project.Id, hostName).Execute()
			if err != nil {
				resData, _ := client.ParseResponseData(httpRes)
				if data, ok := resData.(map[string]interface{}); ok {
					if innerData, ok := data["data"].(map[string]interface{}); ok {
						return charm.RenderError("unable to delete host: ", errors.New(innerData["name"].(string)))
					}
				}
				printer.Logger.Debug("unable to delete host, request failed",
					slog.String("error", err.Error()),
					slog.Any("response", resData))
				return charm.RenderError("unable to delete host", err)
			}

			var data interface{}
			if common.OutputFormat == "json" {
				data = DeleteResp
			} else {
				data = map[string]interface{}{
					"sucessfully deleted host with name": hostName,
				}
			}
			printer.PrintResource(charm.RenderWarning(utils.FormatOutput(data, common.OutputFormat)))
			return nil

		},
	}
	_ = cmd.MarkFlagRequired("name")
	_ = cmd.MarkFlagRequired("project")

	return cmd
}
