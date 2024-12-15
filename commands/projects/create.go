package projects

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

var (
	orgID       string
	projectName string
)

func NewCreateCmd(printer *utils.Printer) *cobra.Command {

	cmd := &cobra.Command{
		Use:     "create",
		Aliases: []string{"new"},
		Example: "qernal project create --name <project_name> --ogranisation <org ID> ",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			token, err := auth.GetQernalToken()
			if err != nil {
				return charm.RenderError("unable to retreive qernal token, run qernal auth login if you haven't")

			}
			qc, err := client.New(ctx, nil, nil, token)
			if err != nil {
				return charm.RenderError("", err)
			}

			project, _, err := qc.ProjectsAPI.ProjectsCreate(ctx).ProjectBody(openapi_chaos_client.ProjectBody{
				OrgId: orgID,
				Name:  projectName,
			}).Execute()
			if err != nil {
				return charm.RenderError("unable to create project", err)

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

			printer.PrintResource(utils.FormatOutput(data, common.OutputFormat))

			return nil
		},
	}
	cmd.Flags().StringVarP(&orgID, "organisation", "", "", "Organisation the project should be created under")
	cmd.Flags().StringVarP(&projectName, "name", "n", "", "Name of the project")
	cmd.Flags().StringVarP(&common.OutputFormat, "output", "o", "text", "output format (json,text)")
	cmd.MarkFlagRequired("organisation")
	cmd.MarkFlagRequired("name")
	return cmd
}
