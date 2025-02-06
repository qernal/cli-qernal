package hosts

import (
	"context"
	"errors"

	"github.com/qernal/cli-qernal/charm"
	"github.com/qernal/cli-qernal/commands/auth"
	"github.com/qernal/cli-qernal/pkg/client"
	"github.com/qernal/cli-qernal/pkg/common"
	"github.com/qernal/cli-qernal/pkg/utils"
	"github.com/spf13/cobra"
)

func NewListCmd(printer *utils.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls", "l"},
		Short:   "list your qernal hosts",
		Example: "qernal host ls --project <your project name>",
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

			projectName, _ := cmd.Flags().GetString("project")

			project, err := qc.GetProjectByName(projectName)
			if err != nil {
				return printer.RenderError("❌", err)
			}

			hostResp, httpRes, err := qc.HostsAPI.ProjectsHostsList(ctx, project.Id).Execute()
			if err != nil {
				resData, _ := client.ParseResponseData(httpRes)
				if data, ok := resData.(map[string]interface{}); ok {
					if innerData, ok := data["data"].(map[string]interface{}); ok {
						if nameErr, ok := innerData["name"].(string); ok {
							return printer.RenderError("unable to list hosts", errors.New(nameErr))
						}
					}
				}
			}

			if common.OutputFormat == "json" {
				printer.PrintResource(utils.FormatOutput(hostResp.Data, common.OutputFormat))
				return nil
			}

			table := charm.RenderHostTable(hostResp.Data)
			printer.PrintResource(table)
			return nil
		},
	}
	cmd.Flags().StringVarP(&common.OutputFormat, "output", "o", "text", "output format (json,text)")
	_ = cmd.MarkFlagRequired("project")
	return cmd
}
