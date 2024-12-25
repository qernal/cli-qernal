package org

import (
	"context"
	"log/slog"

	"github.com/qernal/cli-qernal/charm"
	"github.com/qernal/cli-qernal/commands/auth"
	"github.com/qernal/cli-qernal/pkg/client"
	"github.com/qernal/cli-qernal/pkg/common"
	"github.com/qernal/cli-qernal/pkg/utils"
	"github.com/spf13/cobra"
)

func NewOrgListCmd(printer *utils.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls", "l"},
		Short:   "list your qernal organisations",
		Example: "qernal projects ls",
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
			orgResp, httpRes, err := qc.OrganisationsAPI.OrganisationsList(ctx).Execute()
			if err != nil {
				resData, _ := client.ParseResponseData(httpRes)
				printer.Logger.Debug("unable to list projects, request failed", slog.String("error", err.Error()), slog.String("response", resData.(string)))
				charm.RenderError("unable to list projects,  request failed with:", err)
			}

			if common.OutputFormat == "json" {
				printer.PrintResource(utils.FormatOutput(orgResp.Data, common.OutputFormat))
				return nil
			}

			table := charm.RenderOrgTable(orgResp.Data)
			printer.PrintResource(table)
			//TODO: paginate results
			return nil
		},
	}
	cmd.Flags().StringVarP(&common.OutputFormat, "output", "o", "text", "output format (json,text)")

	return cmd
}
