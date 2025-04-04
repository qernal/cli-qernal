package secrets

import (
	"bufio"
	"context"
	"io"
	"strings"

	"github.com/qernal/cli-qernal/charm"
	"github.com/qernal/cli-qernal/commands/auth"
	"github.com/qernal/cli-qernal/pkg/client"
	"github.com/qernal/cli-qernal/pkg/common"
	"github.com/qernal/cli-qernal/pkg/helpers"
	utils "github.com/qernal/cli-qernal/pkg/utils"
	"github.com/spf13/cobra"
)

func NewEncryptCmd(printer *utils.Printer) *cobra.Command {

	cmd := &cobra.Command{
		Use:     "encrypt",
		Short:   "encrypt a secret from stdin",
		Example: "qernal encrypt <plaintext>",
		Aliases: []string{
			"e",
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return helpers.ValidateProjectFlags(cmd)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				return charm.RenderError("No arguments expected. Please provide input through stdin.")
			}

			// Read from stdin
			reader := bufio.NewReader(cmd.InOrStdin())
			var sb strings.Builder
			for {
				line, err := reader.ReadString('\n')
				if err != nil {
					if err == bufio.ErrBufferFull {
						sb.WriteString(line)
						continue
					} else if err == io.EOF {
						sb.WriteString(line)
						break
					} else {
						return charm.RenderError("Error reading input from stdin:", err)
					}
				}
				sb.WriteString(line)
			}
			plaintext := sb.String()

			// Remove trailing newline from input
			plaintext = strings.TrimSpace(plaintext)
			ctx := context.Background()

			token, err := auth.GetQernalToken()

			if err != nil {
				return charm.RenderError("unable to retreive qernal token, run qernal auth login if you haven't")

			}

			qc, err := client.New(ctx, nil, nil, token)

			if err != nil {
				return charm.RenderError("", err)
			}

			projectID, err := helpers.GetProjectID(cmd, &qc)
			if err != nil {
				return err
			}

			dek, err := qc.FetchDek(ctx, projectID)
			if err != nil {
				return charm.RenderError("unable to fetch dek key", err)
			}

			encryptedVal, err := client.EncryptLocalSecret(dek.Payload.SecretMetaResponseDek.Public, plaintext)
			if err != nil {
				return charm.RenderError("unable to  encrypt input", err)

			}

			var data interface{}
			if common.OutputFormat == "json" {
				data = map[string]interface{}{
					"encrypted_value": encryptedVal,
					"revision_id":     dek.Revision,
				}
			} else {
				data = map[string]interface{}{
					"Encrypted Value": encryptedVal,
					"Revision ID":     dek.Revision,
				}
			}

			printer.PrintResource(utils.FormatOutput(data, common.OutputFormat))

			return nil
		},
	}
	cmd.Flags().StringVarP(&common.OutputFormat, "output", "o", "text", "output format (json,text)")
	_ = cmd.MarkFlagRequired("project")
	return cmd
}
