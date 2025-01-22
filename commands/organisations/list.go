package org

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

func NewOrgListCmd(printer *utils.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls", "l"},
		Short:   "list your qernal organisations",
		Example: "qernal organisations ls",
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
				if data, ok := resData.(map[string]interface{}); ok {
					if innerData, ok := data["data"].(map[string]interface{}); ok {
						if nameErr, ok := innerData["name"].(string); ok {
							return printer.RenderError("unable to create organisation", errors.New(nameErr))
						}
					}
				}
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
