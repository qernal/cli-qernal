package org

import (
	"context"

	"github.com/qernal/cli-qernal/charm"
	"github.com/qernal/cli-qernal/commands/auth"
	"github.com/qernal/cli-qernal/pkg/client"
	"github.com/qernal/cli-qernal/pkg/common"
	"github.com/qernal/cli-qernal/pkg/utils"
	"github.com/spf13/cobra"
)

func NewGetCmd(printer *utils.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "get",
		Aliases: []string{"get"},
		Example: "qernal organisation get --name <org name>",
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

			orgName, _ := cmd.Flags().GetString("name")

			org, err := qc.GetOrgByName(orgName)
			if err != nil {
				return printer.RenderError("x", err)
			}

			var data interface{}
			if common.OutputFormat == "json" {
				data = org
			} else {
				data = map[string]interface{}{
					"Name":    org.Name,
					"User ID": org.UserId,
					"Org ID":  org.Id,
				}
			}

			if common.OutputFormat == "json" {
				printer.PrintResource(utils.FormatOutput(data, common.OutputFormat))
				return nil
			}

			response := charm.SuccessStyle.Render(utils.FormatOutput(data, common.OutputFormat))
			printer.PrintResource(response)
			return nil

		},
	}
	cmd.Flags().StringVar(&orgName, "name", "", "name of the organisation")
	_ = cmd.MarkFlagRequired("name")
	return cmd
}
