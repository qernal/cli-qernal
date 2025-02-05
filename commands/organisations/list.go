package org

import (
	"context"

	"github.com/qernal/cli-qernal/charm"
	"github.com/qernal/cli-qernal/commands/auth"
	"github.com/qernal/cli-qernal/pkg/client"
	"github.com/qernal/cli-qernal/pkg/common"
	"github.com/qernal/cli-qernal/pkg/helpers"
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

			maxResults, _ := cmd.Flags().GetInt32("max")

			orgs, err := helpers.PaginateOrganisations(printer, ctx, &qc, maxResults)
			if err != nil {
				return charm.RenderError("unable to list organisations", err)
			}
			if maxResults > 0 && len(orgs) > int(maxResults) {
				orgs = orgs[:maxResults]
			}

			if common.OutputFormat == "json" {
				printer.PrintResource(utils.FormatOutput(orgs, common.OutputFormat))
				return nil
			}

			table := charm.RenderOrgTable(orgs)
			printer.PrintResource(table)
			return nil
		},
	}
	cmd.Flags().StringVarP(&common.OutputFormat, "output", "o", "text", "output format (json,text)")

	return cmd
}