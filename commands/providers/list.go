package providers

import (
	"context"
	"errors"

	"github.com/qernal/cli-qernal/pkg/utils"
	"github.com/spf13/cobra"

	"github.com/qernal/cli-qernal/charm"
	"github.com/qernal/cli-qernal/commands/auth"
	"github.com/qernal/cli-qernal/pkg/client"
	"github.com/qernal/cli-qernal/pkg/common"
)

func NewListCmd(printer *utils.Printer) *cobra.Command {
	cmd := cobra.Command{
		Use:     "list",
		Aliases: []string{"ls", "l"},
		Short:   "list qernal providers",
		Example: "qernal providers list",
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
			providerResp, httpRes, err := qc.ProvidersAPI.ProvidersList(ctx).Execute()
			if err != nil {
				resData, _ := client.ParseResponseData(httpRes)
				if data, ok := resData.(map[string]interface{}); ok {
					if innerData, ok := data["data"].(map[string]interface{}); ok {
						if nameErr, ok := innerData["name"].(string); ok {
							return printer.RenderError("unable to list organisations", errors.New(nameErr))
						}
					}
				}
				return charm.RenderError("unable to list providers", err)
			}

			if common.OutputFormat == "json" {
				printer.PrintResource(utils.FormatOutput(providerResp, common.OutputFormat))
				return nil
			}

			table := charm.RenderProviderTable(providerResp.Data)
			printer.PrintResource(table)
			return nil
		},
	}
	return &cmd
}
