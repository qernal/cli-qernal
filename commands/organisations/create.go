package org

import (
	"context"
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

			org, httpRes, err := qc.OrganisationsAPI.OrganisationsCreate(ctx).OrganisationBody(openapi_chaos_client.OrganisationBody{
				Name: orgName,
			}).Execute()
			if err != nil {
				resData, _ := client.ParseResponseData(httpRes)
				printer.Logger.Debug("unable to list projects, request failed", slog.String("error", err.Error()), slog.String("response", resData.(string)))
				return charm.RenderError("unable to create project", err)

			}

			var data interface{}

			if common.OutputFormat == "json" {
				data = map[string]interface{}{
					"organisation_name": org.Name,
					"organisation_id":   org.Id,
					"user_id":           org.UserId,
				}
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
