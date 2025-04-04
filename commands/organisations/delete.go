package org

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
		Example: "qernal organisation delete --name <org name>",
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

			orgName, _ := cmd.Flags().GetString("organisation")

			org, err := qc.GetOrgByName(orgName)
			if err != nil {
				return charm.RenderError("x", err)
			}

			DeleteResp, httpRes, err := qc.OrganisationsAPI.OrganisationsDelete(ctx, org.Id).Execute()
			if err != nil {
				resData, _ := client.ParseResponseData(httpRes)
				if data, ok := resData.(map[string]interface{}); ok {
					if innerData, ok := data["data"].(map[string]interface{}); ok {
						return charm.RenderError("unable to delete organisation: ", errors.New(innerData["name"].(string)))
					}
				}
				printer.Logger.Debug("unable to delete org, request failed",
					slog.String("error", err.Error()),
					slog.Any("response", resData))
				return charm.RenderError("unable to delete organisation", err)
			}

			var data interface{}
			if common.OutputFormat == "json" {
				data = DeleteResp
			} else {
				data = map[string]interface{}{
					"sucessfully deleted organisation with name:": org.Name,
				}
			}
			printer.PrintResource(charm.RenderWarning(utils.FormatOutput(data, common.OutputFormat)))
			return nil

		},
	}
	_ = cmd.MarkFlagRequired("orgnaisation")
	return cmd
}
