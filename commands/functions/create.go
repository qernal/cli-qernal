package functions

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/qernal/cli-qernal/charm"
	"github.com/qernal/cli-qernal/commands/auth"
	"github.com/qernal/cli-qernal/pkg/client"
	"github.com/qernal/cli-qernal/pkg/helpers"
	"github.com/qernal/cli-qernal/pkg/utils"
	"github.com/spf13/cobra"
)

func NewCreateCmd(printer *utils.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create",
		Aliases: []string{"new"},
		Example: "qernal functions create -f function.yaml",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()

			token, err := auth.GetQernalToken()
			if err != nil {
				return charm.RenderError("error creating qernal client", err)
			}

			qc, err := client.New(ctx, nil, nil, token)
			if err != nil {
				return charm.RenderError("error creating qernal client", err)
			}

			file, _ := cmd.Flags().GetString("file")

			qFunctions, err := helpers.ParseFunctionConfig(file, printer)
			if err != nil {
				return charm.RenderError("unable to parse function config", err)
			}

			for _, function := range qFunctions {
				qFunc, httpRes, err := qc.FunctionsAPI.FunctionsCreate(ctx).FunctionBody(function).Execute()
				if err != nil {
					resData, _ := client.ParseResponseData(httpRes)
					if data, ok := resData.(map[string]interface{}); ok {
						if innerData, ok := data["data"].(map[string]interface{}); ok {
							if nameErr, ok := innerData["name"].(string); ok {
								err = errors.New(nameErr)
							}
						}
					}
					printer.Logger.Debug("unable to create function, request failed",
						slog.String("error", err.Error()),
						slog.Any("response", resData))
					return printer.RenderError(fmt.Sprintf("unable to create function with name %s. Request failed with", qFunc.Name), err)
				}
				printer.PrintResource(charm.SuccessStyle.Render(fmt.Sprintf("created function %s", function.Name)))
			}

			printer.PrintResource(charm.SuccessStyle.Render("Success.\nrun qernal function ls --project-id=<project-id> to view all functions"))

			return nil
		},
	}

	cmd.Flags().StringVarP(&functionFile, "file", "f", "", "path to function definition file (yaml)")
	_ = cmd.MarkFlagRequired("file")

	return cmd
}
