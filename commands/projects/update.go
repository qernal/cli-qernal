package projects

import (
	"context"
	"fmt"

	"github.com/qernal/cli-qernal/charm"
	"github.com/qernal/cli-qernal/commands/auth"
	"github.com/qernal/cli-qernal/pkg/client"
	openapi_chaos_client "github.com/qernal/openapi-chaos-go-client"
	"github.com/spf13/cobra"
)

var name string
var projectId string
var updateCmd = &cobra.Command{
	Use:     "update",
	Aliases: []string{"edit"},
	Example: "qernal projects update --project=<project ID> --org <org ID> --name <name>",
	Short:   "edit qernal project name",
	Long:    "qernal projects update --project=<project ID>  if --name is not supplied cli will prompt for a new project name",
	RunE: func(cmd *cobra.Command, args []string) error {

		ctx := context.Background()
		token, err := auth.GetQernalToken()
		if err != nil {
			return charm.RenderError("unable to retreive qernal token, run qernal auth login if you haven't")

		}

		qc, err := client.New(ctx, token)
		if err != nil {
			return charm.RenderError("error creating qernal client", err)
		}

		patchResp, _, err := qc.ProjectsAPI.ProjectsUpdate(ctx, projectId).ProjectBodyPatch(openapi_chaos_client.ProjectBodyPatch{
			OrgId: &orgID,
			Name:  &name,
		}).Execute()
		if err != nil {
			return charm.RenderError(fmt.Sprintf("unable to update project, patch failed with: %s", err))
		}

		fmt.Println(charm.RenderWarning("sucessfully updated project name to: " + patchResp.Name))

		return nil

	},
}

func init() {
	updateCmd.Flags().StringVarP(&projectId, "project", "p", "", "Project ID")
	updateCmd.Flags().StringVarP(&name, "name", "n", "", "name of the project to be updated")
	updateCmd.Flags().StringVarP(&orgID, "organisation", "", "", "organisation of the project to be updated")

	updateCmd.MarkFlagRequired("project")
	updateCmd.MarkFlagRequired("name")
}
