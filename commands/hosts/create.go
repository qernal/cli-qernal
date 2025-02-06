package hosts

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/qernal/cli-qernal/charm"
	"github.com/qernal/cli-qernal/commands/auth"
	"github.com/qernal/cli-qernal/pkg/client"
	"github.com/qernal/cli-qernal/pkg/common"
	"github.com/qernal/cli-qernal/pkg/helpers"
	"github.com/qernal/cli-qernal/pkg/utils"
	openapi_chaosclient "github.com/qernal/openapi-chaos-go-client"
	"github.com/spf13/cobra"
)

func NewCreateCmd(printer *utils.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create",
		Aliases: []string{"new"},
		Example: "qernal host create --name example.org --project landing-page --cert MY-CERT-2025",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return helpers.ValidateProjectFlags(cmd)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			token, err := auth.GetQernalToken()
			if err != nil {
				return charm.RenderError("unable to retrieve qernal token, run qernal auth login if you haven't")
			}

			qc, err := client.New(ctx, nil, nil, token)
			if err != nil {
				return charm.RenderError("", err)
			}

			hostName, _ := cmd.Flags().GetString("name")
			cert, _ := cmd.Flags().GetString("cert")
			isDisabled, _ := cmd.Flags().GetBool("disable")

			projectID, err := helpers.GetProjectID(cmd, &qc)
			if err != nil {
				return err
			}
			host, httpRes, err := qc.HostsAPI.ProjectsHostsCreate(ctx, projectID).HostBody(openapi_chaosclient.HostBody{
				Host:        hostName,
				Certificate: fmt.Sprintf("projects:%s/%s", projectID, strings.ToUpper(cert)),
				Disabled:    isDisabled,
			}).Execute()

			if err != nil {
				resData, _ := client.ParseResponseData(httpRes)
				if data, ok := resData.(map[string]interface{}); ok {
					if innerData, ok := data["data"].(map[string]interface{}); ok {
						if nameErr, ok := innerData["name"].(string); ok {
							return printer.RenderError("unable to create host", errors.New(nameErr))
						}
					}
				}
				printer.Logger.Debug("unable to create host, request failed",
					slog.String("error", err.Error()),
					slog.Any("response", resData))
				return printer.RenderError("unable to create host", err)
			}

			var data interface{}
			if common.OutputFormat == "json" {
				data = host
			} else {
				data = map[string]interface{}{
					"Created host at": host.Host,
					"Host ID":         host.Id,
					"Enabled":         host.Disabled,
				}
			}
			dnsRecords := map[string]string{
				"A":    publicIPV4,
				"AAAA": publicIPV6,
				"TXT":  host.TxtVerification,
			}

			printer.PrintResource(utils.FormatOutput(data, common.OutputFormat))
			printer.PrintResource(charm.RenderWarning("Please add the TXT record for host verification, then update your A and AAAA records\nwhen verification is complete. The host will not be routable until the verification has completed.\n"))
			fmt.Println(charm.RenderDNSTable(dnsRecords))
			return nil
		},
	}

	cmd.Flags().StringVarP(&hostName, "name", "n", "", "name of the host")

	_ = cmd.MarkFlagRequired("name")
	_ = cmd.MarkFlagRequired("project")
	_ = cmd.MarkFlagRequired("cert")

	return cmd
}
