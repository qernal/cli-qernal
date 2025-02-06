package hosts

import (
	"context"
	"errors"
	"fmt"

	"github.com/qernal/cli-qernal/charm"
	"github.com/qernal/cli-qernal/commands/auth"
	"github.com/qernal/cli-qernal/pkg/client"
	"github.com/qernal/cli-qernal/pkg/common"
	"github.com/qernal/cli-qernal/pkg/helpers"
	"github.com/qernal/cli-qernal/pkg/utils"
	openapi_chaos_client "github.com/qernal/openapi-chaos-go-client"
	"github.com/spf13/cobra"
)

func NewVerifyCmd(printer *utils.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "verify",
		Aliases: []string{"verify"},
		Short:   "schdeule a host for (re)verification",
		Example: "qernal hosts verify example.org --project landing-page",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return helpers.ValidateProjectFlags(cmd)
		},
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

			projectName, _ := cmd.Flags().GetString("project")
			hostName, _ := cmd.Flags().GetString("name")

			projectID, err := helpers.GetProjectID(cmd, &qc)
			if err != nil {
				return err
			}

			// check if host needs verification
			host, httpRes, err := qc.HostsAPI.ProjectsHostsGet(ctx, projectID, hostName).Execute()
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
			if host.VerificationStatus != openapi_chaos_client.HOSTVERIFICATIONSTATUS_FAILED {
				var message string

				if host.VerificationStatus == openapi_chaos_client.HOSTVERIFICATIONSTATUS_PENDING {
					message = fmt.Sprintf("❌ Host '%s' is currently undergoing verification. Please wait for the current verification process to complete before attempting to reverify.", host.Host)
				} else {
					message = fmt.Sprintf("❌ Cannot reverify host '%s' - current status is '%s'. Reverification is only needed for hosts with failed verification status.",
						host.Host,
						string(host.VerificationStatus))
				}

				return printer.RenderError("", errors.New(message))
			}

			hostResp, httpRes, err := qc.HostsAPI.ProjectsHostsVerifyCreate(ctx, projectID, hostName).Execute()
			if err != nil {
				resData, _ := client.ParseResponseData(httpRes)
				if data, ok := resData.(map[string]interface{}); ok {
					if innerData, ok := data["data"].(map[string]interface{}); ok {
						if nameErr, ok := innerData["name"].(string); ok {
							return printer.RenderError("unable to verify hosts", errors.New(nameErr))
						}
					}
				}
				return charm.RenderError("unable to verfiy host", err)
			}

			var data interface{}

			if common.OutputFormat == "json" {
				printer.PrintResource(utils.FormatOutput(hostResp, common.OutputFormat))
				return nil
			} else {
				data = map[string]interface{}{
					"message": "Host verification scheduled successfully",
					"details": fmt.Sprintf("Verifying host '%s' in project '%s'", hostName, projectName),
					"note":    "DNS propagation times are dependent upon your provider. Use 'qernal hosts list' to check verification status.",
				}
			}
			printer.PrintResource(utils.FormatOutput(data, common.OutputFormat))

			return nil
		},
	}
	cmd.Flags().StringVarP(&common.OutputFormat, "output", "o", "text", "output format (json,text)")
	_ = cmd.MarkFlagRequired("project")
	_ = cmd.MarkFlagRequired("name")
	return cmd
}
