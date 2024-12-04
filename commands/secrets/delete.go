package secrets

import (
	"context"
	"fmt"

	"github.com/qernal/cli-qernal/charm"
	"github.com/qernal/cli-qernal/commands/auth"
	"github.com/qernal/cli-qernal/pkg/client"
	"github.com/spf13/cobra"
)

var DeleteCmd = &cobra.Command{
	Use:     "delete",
	Short:   "delete a qernal secret",
	Aliases: []string{"rm", "remove"},
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
		_, _, err = qc.SecretsAPI.ProjectsSecretsDelete(ctx, projectID, secretName).Execute()
		if err != nil {
			charm.RenderError("unable to delete secret,  request failed with:", err)
		}

		fmt.Println(charm.RenderWarning("successfully deleted project with name: " + secretName))

		return nil
	},
}

func init() {
	DeleteCmd.Flags().StringVarP(&secretName, "name", "n", "", "name of the secret")
	DeleteCmd.MarkFlagRequired("project")
	DeleteCmd.MarkFlagRequired("name")
}
