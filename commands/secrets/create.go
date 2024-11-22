package secrets

import (
	"bufio"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/qernal/cli-qernal/charm"
	"github.com/qernal/cli-qernal/commands/auth"
	"github.com/qernal/cli-qernal/pkg/client"
	openapi_chaos_client "github.com/qernal/openapi-chaos-go-client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var secretType string
var registry string
var publicKey string
var privateKey string

// var environmentValue string

var CreateCmd = &cobra.Command{
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
		if len(args) > 0 {
			return charm.RenderError("No arguments expected. Please provide input through stdin.")
		}
		// Read from stdin
		reader := bufio.NewReader(os.Stdin)
		plaintext, err := reader.ReadString('\n')
		if err != nil {
			return charm.RenderError("Error reading input from stdin:", err)
		}
		// Remove trailing newline from input
		plaintext = strings.TrimSpace(plaintext)
		ctx := context.Background()

		token, err := auth.GetQernalToken()

		if err != nil {
			return charm.RenderError("unable to retreive qernal token, run qernal auth login if you haven't")

		}
		qc, err := client.New(ctx, token)
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
			encryptedValue, err := client.EncryptLocalSecret(dek.Payload.SecretMetaResponseDek.Public, plaintext)
			if err != nil {
				charm.RenderError("unable to  encrypt input", err)

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
				charm.RenderError("unable to  create registry secret", err)

			}
			fmt.Println(charm.SuccessStyle.Render("created registry secret with name " + secretName))
			return nil
		case "environment":
			encryptedValue, err := client.EncryptLocalSecret(dek.Payload.SecretMetaResponseDek.Public, plaintext)
			if err != nil {
				charm.RenderError("unable to  encrypt input", err)

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
				charm.RenderError("unable to create envrionment secret", err)

			}
		case "certificate":
			if publicKey == "" || privateKey == "" {
				return charm.RenderError("Both --public-key and --private-key are required for certificate type")
			}

			publicKeyContent, err := ioutil.ReadFile(publicKey)
			if err != nil {
				return charm.RenderError("Unable to read public key file", err)
			}

			privateKeyContent, err := ioutil.ReadFile(privateKey)
			if err != nil {
				return charm.RenderError("Unable to read private key file", err)
			}

			encryptionRef := fmt.Sprintf(`keys/dek/%d`, dek.Revision)
			_, _, err = qc.SecretsAPI.ProjectsSecretsCreate(ctx, projectID).SecretBody(openapi_chaos_client.SecretBody{
				Name:       strings.ToUpper(secretName),
				Encryption: encryptionRef,
				Type:       openapi_chaos_client.SECRETCREATETYPE_CERTIFICATE,
				Payload: openapi_chaos_client.SecretCreatePayload{
					SecretCertificate: &openapi_chaos_client.SecretCertificate{
						Certificate:      strings.TrimSpace(string(publicKeyContent)),
						CertificateValue: strings.TrimSpace(string(privateKeyContent)),
					},
				},
			}).Execute()
			if err != nil {
				return charm.RenderError("Unable to create certificate secret", err)
			}
			fmt.Println(charm.SuccessStyle.Render("Created certificate secret with name " + secretName))
			return nil
		default:
			return charm.RenderError("Invalid secret type. Must be on of 'registry', 'environment', or 'certificate'")
		}
		return nil
	},
}

func init() {
	CreateCmd.MarkFlagRequired("name")
	CreateCmd.MarkFlagRequired("project")
	CreateCmd.Flags().StringVarP(&secretType, "type", "t", "", "type of secret to be created (registry, environment, certificate")
	CreateCmd.Flags().StringVarP(&registry, "registry-url", "r", "", "Url to private container repository (for docker registry use docker.io)")

	CreateCmd.Flags().StringVar(&publicKey, "public-key", "", "File path to the public key for certificate type")
	CreateCmd.Flags().StringVar(&privateKey, "private-key", "", "File path to the private key for certificate type")

	viper.BindPFlag("public-key", CreateCmd.Flags().Lookup("public-key"))
	viper.BindPFlag("private-key", CreateCmd.Flags().Lookup("private-key"))
}
