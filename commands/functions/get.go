package functions

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

func NewGetCmd(printer *utils.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "get",
		Aliases: []string{"get"},
		Example: "qernal function get --function <function ID>",
		Short:   "Get detailed information about a function ",
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

			functionID, _ := cmd.Flags().GetString("function")
			qFunc, httpRes, err := qc.FunctionsAPI.FunctionsGet(ctx, functionID).Execute()
			if err != nil {
				resData, _ := client.ParseResponseData(httpRes)
				if data, ok := resData.(map[string]interface{}); ok {
					if innerData, ok := data["data"].(map[string]interface{}); ok {
						if nameErr, ok := innerData["name"].(string); ok {
							return printer.RenderError("unable to find host", errors.New(nameErr))
						}
					}
				}
				return charm.RenderError("unable to find host", err)
			}

			var data interface{}
			if common.OutputFormat == "json" {
				data = qFunc
			} else {
				// TODO: persist order of this render

				data = map[string]interface{}{
					"Name":        qFunc.Name,
					"Project ID":  qFunc.ProjectId,
					"Function ID": qFunc.Id,
					"Secrets":     len(qFunc.Secrets),
				}
			}

			response := charm.SuccessStyle.Render(utils.FormatOutput(data, common.OutputFormat))
			printer.PrintResource(response)
			return nil

		},
	}
	cmd.Flags().StringVarP(&functionID, "function", "f", "", "function id")
	_ = cmd.MarkFlagRequired("function")

	return cmd
}
