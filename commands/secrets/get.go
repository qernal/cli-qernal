package secrets

import (
	"context"

	"github.com/qernal/cli-qernal/charm"
	"github.com/qernal/cli-qernal/commands/auth"
	"github.com/qernal/cli-qernal/pkg/client"
	"github.com/qernal/cli-qernal/pkg/common"
	"github.com/qernal/cli-qernal/pkg/utils"
	openapi_chaos_client "github.com/qernal/openapi-chaos-go-client"
	"github.com/spf13/cobra"
)

func NewGetCmd(printer *utils.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "get",
		Aliases: []string{"get"},
		Example: "qernal secrets get --name <org name>",
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

			secretName, _ := cmd.Flags().GetString("name")
			projectID, _ := cmd.Flags().GetString("project")

			secret, err := qc.GetSecretByName(secretName, projectID)
			if err != nil {
				return printer.RenderError("x", err)
			}

			var data interface{}
			if common.OutputFormat == "json" {
				data = secret
			} else {
				data = getSecretData(secret)
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
	cmd.Flags().StringVar(&secretName, "name", "", "name of the secret")
	cmd.Flags().StringVar(&projectID, "project", "", "ID of the project")
	_ = cmd.MarkFlagRequired("name")
	_ = cmd.MarkFlagRequired("project")
	return cmd
}

// getSecretData a map of relevant info for each secret type
func getSecretData(secret *openapi_chaos_client.SecretMetaResponse) map[string]interface{} {
	data := map[string]interface{}{}

	switch secret.Type {
	case openapi_chaos_client.SECRETMETATYPE_ENVIRONMENT:
		data["Name"] = secret.Name
		data["Type"] = secret.Type
		data["Revision"] = secret.Revision
		return data
	case openapi_chaos_client.SECRETMETATYPE_REGISTRY:
		data["Name"] = secret.Name
		data["Type"] = secret.Type
		data["Revision"] = secret.Revision
		data["Registry"] = secret.Payload.SecretMetaResponseRegistryPayload.Registry

	case openapi_chaos_client.SECRETMETATYPE_CERTIFICATE:
		data["Name"] = secret.Name
		data["Type"] = secret.Type
		data["Revision"] = secret.Revision
		data["Certificate"] = secret.Payload.SecretMetaResponseCertificatePayload.Certificate

	}

	return data
}
