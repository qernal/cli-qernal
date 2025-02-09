package secrets

import (
	"context"
	"fmt"

	"github.com/qernal/cli-qernal/charm"
	"github.com/qernal/cli-qernal/commands/auth"
	"github.com/qernal/cli-qernal/pkg/client"
	"github.com/qernal/cli-qernal/pkg/common"
	"github.com/qernal/cli-qernal/pkg/helpers"
	"github.com/qernal/cli-qernal/pkg/utils"
	"github.com/spf13/cobra"
)

func NewSecretsListCmd(printer *utils.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls", "l"},
		Short:   "list your qernal project secrets",
		Example: "qernal secrets list",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return helpers.ValidateProjectFlags(cmd)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			token, err := auth.GetQernalToken()
			if err != nil {
				return charm.RenderError("unable to retreive qernal token, run qernal auth login if you haven't")
			}
			ctx := context.Background()
			qc, err := client.New(ctx, nil, nil, token)
			if err != nil {
				return charm.RenderError("", err)
			}
			maxResults, _ := cmd.Flags().GetInt32("max")
			projectID, err := helpers.GetProjectID(cmd, &qc)
			if err != nil {
				return err
			}
			secrets, err := helpers.PaginateSecrets(printer, ctx, &qc, maxResults, projectID)
			if err != nil {
				return charm.RenderError("unable to list secrets", err)
			}
			if maxResults > 0 && len(secrets) > int(maxResults) {
				secrets = secrets[:maxResults]
			}

			if common.OutputFormat == "json" {
				fmt.Println(utils.FormatOutput(secrets, common.OutputFormat))
				return nil
			}
			table := charm.RenderSecretsTable(secrets)
			fmt.Println(table)
			return nil
		},
	}
	_ = cmd.MarkFlagRequired("project")
	return cmd
}
