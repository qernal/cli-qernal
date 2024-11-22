package secrets

import (
	"context"
	"fmt"

	"github.com/qernal/cli-qernal/charm"
	"github.com/qernal/cli-qernal/commands/auth"
	"github.com/qernal/cli-qernal/pkg/client"
	"github.com/qernal/cli-qernal/pkg/common"
	utils "github.com/qernal/cli-qernal/pkg/uitls"
	"github.com/spf13/cobra"
)

var SecretsListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls", "l"},
	Short:   "list your qernal project secrets",
	Example: "qernal secrets list",
	RunE: func(cmd *cobra.Command, args []string) error {
		token, err := auth.GetQernalToken()
		if err != nil {
			return charm.RenderError("unable to retreive qernal token, run qernal auth login if you haven't")
		}
		ctx := context.Background()
		qc, err := client.New(ctx, token)
		if err != nil {
			return charm.RenderError("", err)
		}
		secretsResp, _, err := qc.SecretsAPI.ProjectsSecretsList(ctx, projectID).Execute()
		if err != nil {
			charm.RenderError("unable to list secrets,  request failed with:", err)
		}

		if secretsResp == nil {
			fmt.Println(charm.RenderWarning("no secrets found in project"))
			return nil
		}

		if common.OutputFormat == "json" {
			fmt.Println(utils.FormatOutput(secretsResp.Data, common.OutputFormat))
			return nil
		}
		table := charm.RenderSecretsTable(secretsResp.Data)
		fmt.Println(table)
		return nil
	},
}

func init() {
	SecretsListCmd.MarkFlagRequired("project")
}
