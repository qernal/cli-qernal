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
		Example: "qernal organisation create --name <organisation_name>",
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
						if nameErr, ok := innerData["name"].(string); ok {
							return printer.RenderError("unable to create organisation", errors.New(nameErr))
						}
					}
				}
				printer.Logger.Debug("unable to list organisations, request failed",
					slog.String("error", err.Error()),
					slog.Any("response", resData))

				return printer.RenderError("unable to create organisation", err)
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
	_ = cmd.MarkFlagRequired("name")
	return cmd
}