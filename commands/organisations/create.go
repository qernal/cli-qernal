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
	openapi_chaos_client "github.com/qernal/openapi-chaos-go-client"
	"github.com/spf13/cobra"
)

func NewCreateCmd(printer *utils.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create",
		Aliases: []string{"new"},
		Example: "qernal organisation create --name <project_name>",
		RunE: func(cmd *cobra.Command, args []string) error {

			ctx := context.Background()
			token, err := auth.GetQernalToken()
			if err != nil {
				return charm.RenderError("unable to retreive qernal token, run qernal auth login if you haven't")

			}
			qc, err := client.New(ctx, nil, nil, token)
			if err != nil {
				return charm.RenderError("", err)
			}

			orgName, _ := cmd.Flags().GetString("name")

			org, httpRes, err := qc.OrganisationsAPI.OrganisationsCreate(ctx).OrganisationBody(openapi_chaos_client.OrganisationBody{
				Name: orgName,
			}).Execute()
			if err != nil {
				resData, _ := client.ParseResponseData(httpRes)
				if data, ok := resData.(map[string]interface{}); ok {
					if innerData, ok := data["data"].(map[string]interface{}); ok {
						return charm.RenderError("unable to create project: ", errors.New(innerData["name"].(string)))
					}
				}
				printer.Logger.Debug("unable to list projects, request failed",
					slog.String("error", err.Error()),
					slog.Any("response", resData))
				return charm.RenderError("unable to create project", err)
			}
			var data interface{}

			if common.OutputFormat == "json" {
				data = org
			} else {
				data = map[string]interface{}{
					"Created organisation with ID": org.Id,
				}
			}
			printer.PrintResource(utils.FormatOutput(data, common.OutputFormat))

			return nil
		},
	}
	cmd.Flags().StringVarP(&orgName, "name", "n", "", "name of the organisation")
	cmd.Flags().StringVarP(&common.OutputFormat, "output", "o", "text", "output format (json,text)")
	cmd.MarkFlagRequired("name")

	return cmd
}
