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
	openapi_chaos_client "github.com/qernal/openapi-chaos-go-client"
	"github.com/spf13/cobra"
)

func NewUpdateCmd(printer *utils.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update ",
		Aliases: []string{"edit"},
		Example: "qernal function update --file functions.yaml",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()

			token, err := auth.GetQernalToken()
			if err != nil {
				return charm.RenderError("error creating qernal client", err)

			}

			qc, err := client.New(ctx, nil, nil, token)
			if err != nil {
				charm.RenderError("error creating qernal client", err)
			}

			functionID, _ := cmd.Flags().GetString("function")
			file, _ := cmd.Flags().GetString("file")

			qFunctions, err := helpers.ParseFunctionConfig(file, printer)
			if err != nil {
				return charm.RenderError("unable to parse function config", err)
			}

			// Get the function first to verify it exists
			qFunc, httpRes, err := qc.FunctionsAPI.FunctionsGet(ctx, functionID).Execute()
			if err != nil {
				resData, _ := client.ParseResponseData(httpRes)
				if data, ok := resData.(map[string]interface{}); ok {
					if innerData, ok := data["data"].(map[string]interface{}); ok {
						if nameErr, ok := innerData["name"].(string); ok {
							return printer.RenderError("unable to find function", errors.New(nameErr))
						}
					}
				}
				return charm.RenderError("unable to find function", err)
			}

			var matchedFunction *openapi_chaos_client.FunctionBody
			for _, function := range qFunctions {
				if function.Name == qFunc.Name {
					matchedFunction = &function
					break
				}
			}
			if matchedFunction == nil {
				return printer.RenderError("function not found in config file", fmt.Errorf("no matching function with name %s found in config", qFunc.Name))
			}
			// Update the function
			deployments := helpers.OpenAPIDeploymentsToDeployments(matchedFunction.Deployments)

			funcBody := &openapi_chaos_client.Function{
				Id:          functionID,
				ProjectId:   matchedFunction.ProjectId,
				Version:     matchedFunction.Version,
				Name:        matchedFunction.Name,
				Description: matchedFunction.Description,
				Image:       matchedFunction.Image,
				Revision:    qFunc.Revision,
				Type:        matchedFunction.Type,
				Size:        matchedFunction.Size,
				Port:        matchedFunction.Port,
				Routes:      matchedFunction.Routes,
				Scaling:     matchedFunction.Scaling,
				Deployments: deployments,
				Secrets:     matchedFunction.Secrets,
				Compliance:  matchedFunction.Compliance,
			}

			updatedFunc, httpRes, err := qc.FunctionsAPI.FunctionsUpdate(ctx, functionID).Function(*funcBody).Execute()
			if err != nil {
				resData, _ := client.ParseResponseData(httpRes)
				if data, ok := resData.(map[string]interface{}); ok {
					if innerData, ok := data["data"].(map[string]interface{}); ok {
						if nameErr, ok := innerData["name"].(string); ok {
							err = errors.New(nameErr)
						}
					}
				}
				printer.Logger.Debug("unable to update function, request failed",
					slog.String("error", err.Error()),
					slog.Any("response", resData))
				return printer.RenderError(fmt.Sprintf("unable to update function with name %s. Request failed", matchedFunction.Name), err)
			}

			printer.PrintResource(charm.SuccessStyle.Render(fmt.Sprintf("Updated function %s.\nRun qernal function ls --project-id=<project-id> to view all functions", updatedFunc.Name)))

			return nil
		},
	}
	cmd.Flags().StringVar(&functionID, "function", "", "function id")
	cmd.Flags().StringVarP(&functionFile, "file", "f", "", "path to function definition file (yaml)")
	_ = cmd.MarkFlagRequired("file")
	_ = cmd.MarkFlagRequired("function")

	return cmd
}
