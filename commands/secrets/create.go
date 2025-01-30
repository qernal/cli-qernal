package secrets

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/qernal/cli-qernal/charm"
	"github.com/qernal/cli-qernal/commands/auth"
	"github.com/qernal/cli-qernal/pkg/client"
	"github.com/qernal/cli-qernal/pkg/common"
	"github.com/qernal/cli-qernal/pkg/utils"
	openapi_chaos_client "github.com/qernal/openapi-chaos-go-client"
	"github.com/spf13/cobra"
)

var secretType string
var registry string
var publicKey string
var privateKey string

// var environmentValue string

func NewCreateCmd(printer *utils.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create",
		Aliases: []string{"new"},
		Short:   "Create a new secret",
		Long: `Create a new secret in your Qernal project.
This command supports creating registry, environment, and certificate secrets.
The secret value is read from stdin, allowing for secure input methods.`,
		Example: ` # Create a registry secret
  echo <registry-password> | qernal secret create --name MyRegistrySecret --type registry --registry-url docker.io

  # Create an environment secret
  echo <environment-value> | qernal secret create --name MyEnvSecret --type environment

  # Create a certificate secret
  qernal secret create --name MyCertSecret --type certificate --public-key /path/to/public.key --private-key /path/to/private.key`,
		RunE: func(cmd *cobra.Command, args []string) error {

			if len(args) > 0 && secretType != "certificate" {
				return charm.RenderError("No arguments expected. Please provide input through stdin.")
			}

			ctx := context.Background()

			token, err := auth.GetQernalToken()

			if err != nil {
				return charm.RenderError("unable to retrieve qernal token, run qernal auth login if you haven't")

			}
			qc, err := client.New(ctx, nil, nil, token)
			if err != nil {
				return charm.RenderError("", err)
			}

			dek, err := qc.FetchDek(ctx, projectID)
			if err != nil {
				return charm.RenderError("unable to fetch dek key", err)
			}
			switch secretType {
			case "registry":
				// mark registry flag required
				if registry == "" {
					return charm.RenderError("--registry is required to create a registry secret")
				}
				// Read from stdin
				reader := bufio.NewReader(cmd.InOrStdin())
				plaintext, err := reader.ReadString('\n')
				if err != nil {
					return charm.RenderError("Error reading input from stdin:", err)
				}
				// Remove trailing newline from input
				plaintext = strings.TrimSpace(plaintext)
				encryptedValue, err := client.EncryptLocalSecret(dek.Payload.SecretMetaResponseDek.Public, plaintext)
				if err != nil {
					return charm.RenderError("unable to  encrypt input", err)

				}
				encryptionRef := fmt.Sprintf(`keys/dek/%d`, dek.Revision)
				_, _, err = qc.SecretsAPI.ProjectsSecretsCreate(ctx, projectID).SecretBody(openapi_chaos_client.SecretBody{
					Name:       strings.ToUpper(secretName),
					Encryption: encryptionRef,
					Type:       openapi_chaos_client.SECRETCREATETYPE_REGISTRY,
					Payload: openapi_chaos_client.SecretCreatePayload{
						SecretRegistry: &openapi_chaos_client.SecretRegistry{
							Registry:      registry,
							RegistryValue: encryptedValue,
						},
					},
				}).Execute()
				if err != nil {
					return charm.RenderError("unable to  create registry secret", err)

				}
				printer.PrintResource(charm.SuccessStyle.Render("created registry secret with name " + secretName))
				return nil
			case "environment":
				// Read from stdin
				reader := bufio.NewReader(cmd.InOrStdin())
				plaintext, err := reader.ReadString('\n')
				if err != nil {
					return charm.RenderError("Error reading input from stdin:", err)
				}
				// Remove trailing newline from input
				plaintext = strings.TrimSpace(plaintext)
				encryptedValue, err := client.EncryptLocalSecret(dek.Payload.SecretMetaResponseDek.Public, plaintext)
				if err != nil {
					return charm.RenderError("unable to  encrypt input", err)
				}

				encryptionRef := fmt.Sprintf(`keys/dek/%d`, dek.Revision)
				_, _, err = qc.SecretsAPI.ProjectsSecretsCreate(ctx, projectID).SecretBody(openapi_chaos_client.SecretBody{
					Name:       strings.ToUpper(secretName),
					Encryption: encryptionRef,
					Type:       openapi_chaos_client.SECRETCREATETYPE_ENVIRONMENT,
					Payload: openapi_chaos_client.SecretCreatePayload{
						SecretEnvironment: &openapi_chaos_client.SecretEnvironment{
							EnvironmentValue: encryptedValue,
						},
					},
				}).Execute()
				if err != nil {
					return charm.RenderError("unable to create environment secret", err)
				}
				printer.PrintResource(charm.SuccessStyle.Render("created environment secret with name " + secretName))
			case "certificate":

				if publicKey == "" || privateKey == "" {
					return charm.RenderError("Both --public-key and --private-key flags must be provided", nil)
				}

				publicKeyContent, err := os.ReadFile(publicKey)
				if err != nil {
					return charm.RenderError("Unable to read public key file", err)
				}
				privateKeyContent, err := os.ReadFile(privateKey)
				if err != nil {
					return charm.RenderError("Unable to read private key file", err)
				}

				// encrypt private key
				privateKeyEncrypted, err := client.EncryptLocalSecret(dek.Payload.SecretMetaResponseDek.Public, strings.TrimSpace(string(privateKeyContent)))
				if err != nil {
					return charm.RenderError("unable to private key", err)
				}

				encryptionRef := fmt.Sprintf(`keys/dek/%d`, dek.Revision)
				_, _, err = qc.SecretsAPI.ProjectsSecretsCreate(ctx, projectID).SecretBody(openapi_chaos_client.SecretBody{
					Name:       strings.ToUpper(secretName),
					Encryption: encryptionRef,
					Type:       openapi_chaos_client.SECRETCREATETYPE_CERTIFICATE,
					Payload: openapi_chaos_client.SecretCreatePayload{
						SecretCertificate: &openapi_chaos_client.SecretCertificate{
							Certificate:      strings.TrimSpace(string(publicKeyContent)),
							CertificateValue: privateKeyEncrypted,
						},
					},
				}).Execute()
				if err != nil {
					return charm.RenderError("Unable to create certificate secret", err)
				}
				printer.PrintResource(charm.SuccessStyle.Render("Created certificate secret with name " + secretName))
				return nil
			default:
				return charm.RenderError("Invalid secret type. Must be on of 'registry', 'environment', or 'certificate'")
			}
			return nil
		},
	}

	_ = cmd.MarkFlagRequired("name")
	_ = cmd.MarkFlagRequired("project")
	cmd.Flags().StringVarP(&secretType, "type", "t", "", "type of secret to be created (registry, environment, certificate")
	cmd.Flags().StringVarP(&registry, "registry-url", "r", "", "Url to private container repository (for docker registry use docker.io)")

	cmd.Flags().StringVarP(&secretName, "name", "n", "", "name of the secret")
	cmd.Flags().StringVarP(&projectID, "project", "p", "", "ID of the project")
	cmd.Flags().StringVarP(&common.OutputFormat, "output", "o", "text", "output format (json,text)")

	cmd.Flags().StringVarP(&publicKey, "public-key", "", "", "File path to the public key for certificate type")
	cmd.Flags().StringVarP(&privateKey, "private-key", "", "", "File path to the private key for certificate type")

	return cmd
}
