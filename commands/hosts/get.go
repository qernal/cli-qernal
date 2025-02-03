package hosts

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
		Example: "qernal hosts get --name <host name>",
		Short:   "Get detailed information about a specific host",
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

			projectName, _ := cmd.Flags().GetString("project")
			hostName, _ := cmd.Flags().GetString("name")

			project, err := qc.GetProjectByName(projectName)
			if err != nil {
				return printer.RenderError("❌", err)
			}

			// check if host needs verification
			host, httpRes, err := qc.HostsAPI.ProjectsHostsGet(ctx, project.Id, hostName).Execute()
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
				data = host
			} else {
				// Format relevant host information for text output
				data = map[string]interface{}{
					"Hostname":            host.Host,
					"Verification Status": string(host.VerificationStatus),
					"Project ID":          host.ProjectId,
					"State":               getHostState(host.Disabled),
					"Certificate":         getCertificateStatus(host.Certificate),
					"Read Only":           getReadOnlyStatus(host.ReadOnly),
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
	_ = cmd.MarkFlagRequired("name")
	_ = cmd.MarkFlagRequired("project")

	return cmd
}

func getHostState(disabled bool) string {
	if disabled {
		return "Disabled"
	}
	return "Enabled"
}

func getCertificateStatus(cert *string) string {
	if cert == nil || *cert == "" {
		return "Not configured"
	}
	return "Configured"
}

func getReadOnlyStatus(readOnly bool) string {
	if readOnly {
		return "Yes (*.qrnl.app domain)"
	}
	return "No"
}
