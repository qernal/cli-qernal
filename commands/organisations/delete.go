package org

import (
	"context"
	"log/slog"

	"github.com/qernal/cli-qernal/charm"
	"github.com/qernal/cli-qernal/commands/auth"
	"github.com/qernal/cli-qernal/pkg/client"
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

			_, httpRes, err := qc.OrganisationsAPI.OrganisationsDelete(ctx, orgID).Execute()

			if err != nil {
				resData, _ := client.ParseResponseData(httpRes)
				printer.Logger.Debug("unable to delete org, request failed", slog.String("error", err.Error()), slog.String("response", resData.(string)))
				charm.RenderError("unable to delete organisation,request failed with:", err)
			}

			printer.PrintResource(charm.RenderWarning("sucessfully deleted project with ID: " + orgID))
			return nil

		},
	}
	cmd.Flags().StringVar(&orgID, "id", "", "ID of the organisation")
	return cmd
}
