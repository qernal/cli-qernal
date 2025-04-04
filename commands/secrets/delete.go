package secrets

import (
	"context"

	"github.com/qernal/cli-qernal/charm"
	"github.com/qernal/cli-qernal/commands/auth"
	"github.com/qernal/cli-qernal/pkg/client"
	"github.com/qernal/cli-qernal/pkg/helpers"
	"github.com/qernal/cli-qernal/pkg/utils"
	"github.com/spf13/cobra"
)

func NewDeleteCmd(printer *utils.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete",
		Short:   "delete a qernal secret",
		Aliases: []string{"rm", "remove"},
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

			projectID, err := helpers.GetProjectID(cmd, &qc)
			if err != nil {
				return err
			}

			_, _, err = qc.SecretsAPI.ProjectsSecretsDelete(ctx, projectID, secretName).Execute()
			if err != nil {
				return charm.RenderError("unable to delete secret,  request failed with:", err)
			}

			printer.PrintResource(charm.RenderWarning("sucessfully deleted project with name: " + secretName))

			return nil
		},
	}
	cmd.Flags().StringVarP(&secretName, "name", "n", "", "name of the secret")
	_ = cmd.MarkFlagRequired("project")
	_ = cmd.MarkFlagRequired("name")
	return cmd
}
