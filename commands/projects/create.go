package projects

import (
	"context"
	"fmt"

	"github.com/qernal/cli-qernal/charm"
	"github.com/qernal/cli-qernal/commands/auth"
	"github.com/qernal/cli-qernal/pkg/client"
	"github.com/qernal/cli-qernal/pkg/common"
	utils "github.com/qernal/cli-qernal/pkg/uitls"
	openapi_chaos_client "github.com/qernal/openapi-chaos-go-client"
	"github.com/spf13/cobra"
)

var orgID string

var CreateCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"new"},
	Example: "qernal project create <project_name>",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return charm.RenderError("Arguments expected please provide a project name")
		}

		projectName := args[0]
		ctx := context.Background()
		token, err := auth.GetQernalToken()
		if err != nil {
			return charm.RenderError("unable to retreive qernal token, run qernal auth login if you haven't")

		}
		qc, err := client.New(ctx, token)
		if err != nil {
			return charm.RenderError("", err)
		}

		project, _, err := qc.ProjectsAPI.ProjectsCreate(ctx).ProjectBody(openapi_chaos_client.ProjectBody{
			OrgId: orgID,
			Name:  projectName,
		}).Execute()
		if err != nil {
			charm.RenderError("unable to create project", err)

		}
		var data interface{}

		if common.OutputFormat == "json" {
			data = map[string]interface{}{
				"project_name":    project.Name,
				"organisation_id": project.OrgId,
				"project_id":      project.Id,
			}
		} else {
			data = map[string]interface{}{
				"Created project with ID": project.Id,
			}
		}

		fmt.Println(utils.FormatOutput(data, common.OutputFormat))

		return nil
	},
}

func init() {
	CreateCmd.Flags().StringVarP(&orgID, "organisation", "", "", "Organisation ID")
	CreateCmd.MarkFlagRequired("organisation")
}
