package functions

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/qernal/cli-qernal/charm"
	"github.com/qernal/cli-qernal/commands/auth"
	"github.com/qernal/cli-qernal/pkg/client"
	"github.com/qernal/cli-qernal/pkg/utils"
	openapi_chaos_client "github.com/qernal/openapi-chaos-go-client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

			filePath, _ := cmd.Flags().GetString("file")

			viper.SetConfigFile(filePath)
			viper.SetConfigType("yaml")

			if err := viper.ReadInConfig(); err != nil {
				return fmt.Errorf("error reading function file: %w", err)
			}

			var function openapi_chaos_client.FunctionBody
			if err := viper.Unmarshal(&function); err != nil {
				return fmt.Errorf("error parsing function definition: %w", err)
			}

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

				return printer.RenderError("unable to create organisation", err)
			}
			printer.PrintResource(charm.SuccessStyle.Render(qFunc.Name))

			return nil
		},
	}

	cmd.Flags().StringVarP(&functionFile, "file", "f", "", "path to function definition file (yaml)")
	_ = cmd.MarkFlagRequired("file")

	return cmd
}
