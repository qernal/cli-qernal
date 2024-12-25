package org

import (
	"context"
	"fmt"
	"log/slog"

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
				return charm.RenderError("error creating qernal client", err)
			}

			patchResp, httpRes, err := qc.OrganisationsAPI.OrganisationsUpdate(ctx, orgID).OrganisationBody(openapi_chaos_client.OrganisationBody{
				Name: orgName,
			}).Execute()
			if err != nil {
				resData, _ := client.ParseResponseData(httpRes)
				printer.Logger.Debug("unable to delete org, request failed", slog.String("error", err.Error()), slog.String("response", resData.(string)))
				return charm.RenderError(fmt.Sprintf("unable to update org, patch failed with: %s", err))
			}
			printer.PrintResource(charm.RenderWarning("sucessfully updated organisation name to: " + patchResp.Name))
			return nil

		},
	}

	cmd.Flags().StringVar(&orgID, "id", "", "Organisation ID")
	cmd.Flags().StringVarP(&orgName, "name", "n", "", "name of the organisation to be updated")
	cmd.Flags().StringVarP(&common.OutputFormat, "output", "o", "text", "output format (json,text)")
	cmd.MarkFlagRequired("id")
	cmd.MarkFlagRequired("name")
	return cmd

}
