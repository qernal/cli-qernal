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

func NewDeleteCmd(printer *utils.Printer) *cobra.Command {
	var functionID string
	var functionFile string
	var projectID string

	cmd := &cobra.Command{
		Use:     "delete",
		Aliases: []string{"rm"},
		Example: "qernal function delete --function <function id>\nqernal function delete --file function.yaml --project-id <project-id>",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				return errors.New("no arguments expected")
			}

			// Check if either function ID or file is provided
			hasFunction, _ := cmd.Flags().GetString("function")
			hasFile, _ := cmd.Flags().GetString("file")
			hasProject, _ := cmd.Flags().GetString("project-id")

			if hasFunction == "" && hasFile == "" {
				return errors.New("either --function or --file must be specified")
			}

			if hasFile != "" && hasProject == "" {
				return errors.New("when using --file, --project-id must also be specified")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			token, err := auth.GetQernalToken()
			if err != nil {
				return charm.RenderError("unable to retrieve qernal token, run qernal auth login if you haven't", err)
			}

			qc, err := client.New(ctx, nil, nil, token)
			if err != nil {
				return charm.RenderError("error creating qernal client", err)
			}

			// Check if we're deleting by file or by function ID
			file, _ := cmd.Flags().GetString("file")
			if file != "" {
				// Get functions from the file
				qFunctions, err := helpers.ParseFunctionConfig(file, printer)
				if err != nil {
					return charm.RenderError("unable to parse function config", err)
				}

				// Get list of functions from the project
				functions, httpRes, err := qc.FunctionsAPI.ProjectsFunctionsList(ctx, projectID).Execute()
				if err != nil {
					resData, _ := client.ParseResponseData(httpRes)
					printer.Logger.Debug("unable to list functions, request failed",
						slog.String("error", err.Error()),
						slog.Any("response", resData))
					return charm.RenderError("unable to list functions", err)
				}

				// Build a map of function names to IDs for quick lookup
				functionMap := make(map[string]string)
				for _, fn := range functions.Data {
					functionMap[fn.Name] = fn.Id
				}

				// Delete each function that matches
				for _, function := range qFunctions {
					if id, exists := functionMap[function.Name]; exists {
						_, httpRes, err := qc.FunctionsAPI.FunctionsDelete(ctx, id).Execute()
						if err != nil {
							resData, _ := client.ParseResponseData(httpRes)
							printer.Logger.Debug("unable to delete function, request failed",
								slog.String("error", err.Error()),
								slog.Any("response", resData),
								slog.String("function_name", function.Name))
							printer.PrintResource(charm.WarningStyle.Render(
								fmt.Sprintf("failed to delete function %s: %s", function.Name, err.Error())))
							continue
						}
						printer.PrintResource(charm.SuccessStyle.Render(
							fmt.Sprintf("deleted function %s", function.Name)))
					} else {
						printer.PrintResource(charm.WarningStyle.Render(
							fmt.Sprintf("function %s not found in project", function.Name)))
					}
				}
			} else {
				// Delete by function ID
				_, httpRes, err := qc.FunctionsAPI.FunctionsDelete(ctx, functionID).Execute()
				if err != nil {
					resData, _ := client.ParseResponseData(httpRes)
					if data, ok := resData.(map[string]interface{}); ok {
						if innerData, ok := data["data"].(map[string]interface{}); ok {
							if nameErr, ok := innerData["name"].(string); ok {
								return charm.RenderError("unable to delete function: ", errors.New(nameErr))
							}
						}
					}
					printer.Logger.Debug("unable to delete function, request failed",
						slog.String("error", err.Error()),
						slog.Any("response", resData))
					return charm.RenderError("unable to delete function", err)
				}
				printer.PrintResource(charm.SuccessStyle.Render(
					fmt.Sprintf("deleted function %s", functionID)))
			}

			printer.PrintResource(charm.SuccessStyle.Render("Success.\nrun qernal function ls --project-id=<project-id> to view all functions"))
			return nil
		},
	}

	cmd.Flags().StringVar(&functionID, "function", "", "function id")
	cmd.Flags().StringVar(&functionFile, "file", "", "path to function definition file (yaml)")
	cmd.Flags().StringVar(&projectID, "project-id", "", "project id (required when using --file)")

	return cmd
}
