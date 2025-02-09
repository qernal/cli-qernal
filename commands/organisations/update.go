package org

import (
	"context"
	"errors"

	"github.com/qernal/cli-qernal/charm"
	"github.com/qernal/cli-qernal/commands/auth"
	"github.com/qernal/cli-qernal/pkg/client"
	"github.com/qernal/cli-qernal/pkg/common"
	"github.com/qernal/cli-qernal/pkg/utils"
	openapi_chaos_client "github.com/qernal/openapi-chaos-go-client"
	"github.com/spf13/cobra"
)

func NewUpdateCmd(printer *utils.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update",
		Aliases: []string{"edit"},
		Example: "qernal organisation update -id  <org ID> --name <name>",
		Short:   "edit qernal organisation name",
		RunE: func(cmd *cobra.Command, args []string) error {

			ctx := context.Background()
			token, err := auth.GetQernalToken()
			if err != nil {
				return charm.RenderError("unable to retreive qernal token, run qernal auth login if you haven't")

			}

			qc, err := client.New(ctx, nil, nil, token)
			if err != nil {
				return printer.RenderError("error creating qernal client", err)
			}
			orgName, _ := cmd.Flags().GetString("organisation")
			orgID, _ := cmd.Flags().GetString("organisation-id")

			patchResp, httpRes, err := qc.OrganisationsAPI.OrganisationsUpdate(ctx, orgID).OrganisationBody(openapi_chaos_client.OrganisationBody{
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
			}

			var data interface{}
			if common.OutputFormat == "json" {
				data = patchResp
			} else {
				data = map[string]interface{}{
					"sucessfully updated organisation name to:": patchResp.Name,
				}
			}
			printer.PrintResource(charm.RenderWarning(utils.FormatOutput(data, common.OutputFormat)))
			return nil

		},
	}
	cmd.Flags().StringVarP(&common.OutputFormat, "output", "o", "text", "output format (json,text)")
	_ = cmd.MarkFlagRequired("organisation")
	_ = cmd.MarkFlagRequired("organisation-id")
	return cmd

}
