package hosts

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/qernal/cli-qernal/charm"
	"github.com/qernal/cli-qernal/commands/auth"
	"github.com/qernal/cli-qernal/pkg/client"
	"github.com/qernal/cli-qernal/pkg/common"
	"github.com/qernal/cli-qernal/pkg/helpers"
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
				routeable := ""
				if host.VerificationStatus != "completed" && !host.Disabled {
					routeable = " (unroutable, not verified)"
				}

				certName := "None"
				if host.Certificate != nil && *host.Certificate != "" {
					certRefParts := strings.Split(*host.Certificate, "/")
					certName = certRefParts[len(certRefParts)-1]
				}

				// TODO: persist order of this render
				data = map[string]interface{}{
					"Hostname":                host.Host,
					"Project ID":              host.ProjectId,
					"State":                   fmt.Sprintf("%s%s", helpers.GetHostState(host.Disabled), routeable),
					"Certificate":             certName,
					"Read Only":               helpers.GetReadOnlyStatus(host.ReadOnly),
					"Verification TXT Record": host.TxtVerification,
					"Verification Status":     string(host.VerificationStatus),
					"-------------------":     "",
					"A Record":                publicIPV4,
					"AAAA Record":             publicIPV6,
				}
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
