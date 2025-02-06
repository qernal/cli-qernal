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
	"github.com/qernal/cli-qernal/pkg/utils"
	openapi_chaos_client "github.com/qernal/openapi-chaos-go-client"
	"github.com/spf13/cobra"
)

func NewUpdateCmd(printer *utils.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update",
		Aliases: []string{"edit"},
		Example: "qernal host update --cert MY-CERT-2025 --disabled",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			cert, _ := cmd.Flags().GetString("cert")
			isDisabled, _ := cmd.Flags().GetBool("disable")
			if cert == "" && !isDisabled {
				return printer.RenderError("", errors.New("at least one of --cert or --disable must be provided"))
			}
			return nil
		},
		Short: "update the certificate or enable/disable a qernal host",
		RunE: func(cmd *cobra.Command, args []string) error {

			ctx := context.Background()
			token, err := auth.GetQernalToken()
			if err != nil {
				return charm.RenderError("unable to retreive qernal token, run qernal auth login if you haven't")

			}

			qc, err := client.New(ctx, nil, nil, token)
			if err != nil {
				return charm.RenderError("error creating qernal client", err)
			}

			hostName, _ := cmd.Flags().GetString("name")
			cert, _ := cmd.Flags().GetString("cert")
			isDisabled, _ := cmd.Flags().GetBool("disable")
			projectName, _ := cmd.Flags().GetString("project")

			project, err := qc.GetProjectByName(projectName)
			if err != nil {
				return printer.RenderError("❌", err)
			}

			ref := fmt.Sprintf("projects:%s/%s", project.Id, strings.ToUpper(cert))
			_, httpRes, err := qc.HostsAPI.ProjectsHostsUpdate(ctx, project.Id, hostName).HostBodyPatch(openapi_chaos_client.HostBodyPatch{
				Certificate: &ref,
				Disabled:    &isDisabled,
			}).Execute()
			if err != nil {
				resData, _ := client.ParseResponseData(httpRes)
				if data, ok := resData.(map[string]interface{}); ok {
					if innerData, ok := data["data"].(map[string]interface{}); ok {
						if nameErr, ok := innerData["name"].(string); ok {
							err = errors.New(nameErr)
						}
					}
				}
				printer.Logger.Debug("unable to update host, patch failed",
					slog.String("error", err.Error()),
					slog.Any("response", resData))
				return printer.RenderError("unable to update host", err)
			}

			printer.PrintResource(charm.RenderWarning("sucessfully updated host"))

			return nil

		},
	}

	cmd.Flags().StringVarP(&common.OutputFormat, "output", "o", "text", "output format (json,text)")
	_ = cmd.MarkFlagRequired("project")
	_ = cmd.MarkFlagRequired("name")
	return cmd
}
