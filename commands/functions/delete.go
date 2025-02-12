package functions

import (
	"context"
	"errors"
	"log/slog"

	"github.com/qernal/cli-qernal/charm"
	"github.com/qernal/cli-qernal/commands/auth"
	"github.com/qernal/cli-qernal/pkg/client"
	"github.com/qernal/cli-qernal/pkg/common"
	"github.com/qernal/cli-qernal/pkg/utils"
	"github.com/spf13/cobra"
)

func NewDeleteCmd(printer *utils.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete",
		Aliases: []string{"rm"},
		Example: "qernal function delete --function <function id>",

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

			DeleteResp, httpRes, err := qc.FunctionsAPI.FunctionsDelete(ctx, functionID).Execute()
			if err != nil {
				resData, _ := client.ParseResponseData(httpRes)
				if data, ok := resData.(map[string]interface{}); ok {
					if innerData, ok := data["data"].(map[string]interface{}); ok {
						return charm.RenderError("unable to delete host: ", errors.New(innerData["name"].(string)))
					}
				}
				printer.Logger.Debug("unable to delete host, request failed",
					slog.String("error", err.Error()),
					slog.Any("response", resData))
				return charm.RenderError("unable to delete host", err)
			}

			var data interface{}
			if common.OutputFormat == "json" {
				data = DeleteResp
			} else {
				data = map[string]interface{}{
					"sucessfully deleted function": "",
				}
			}
			printer.PrintResource(charm.RenderWarning(utils.FormatOutput(data, common.OutputFormat)))
			return nil

		},
	}
	cmd.Flags().StringVarP(&functionID, "function", "f", "", "function id")
	_ = cmd.MarkFlagRequired("function")

	return cmd
}
