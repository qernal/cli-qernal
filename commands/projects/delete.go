package projects

import (
	"context"
	"fmt"

	"github.com/qernal/cli-qernal/charm"
	"github.com/qernal/cli-qernal/commands/auth"
	"github.com/qernal/cli-qernal/pkg/client"
	"github.com/qernal/cli-qernal/pkg/utils"
	"github.com/spf13/cobra"
)

func NewDeleteCmd(printer *utils.Printer) *cobra.Command {

	cmd := &cobra.Command{
		Use:     "delete",
		Aliases: []string{"rm"},
		Example: "qernal projects delete --project <project ID>",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				return charm.RenderError("No arguments expected")
			}

			ctx := context.Background()
			token, err := auth.GetQernalToken()
			if err != nil {
				return charm.RenderError("unable to retreive qernal token, run qernal auth login if you haven't")

			}

			qc, err := client.New(ctx, nil, nil, token)
			if err != nil {
				return charm.RenderError("error creating qernal client", err)
			}

			projects, httpRes, err := qc.ProjectsAPI.ProjectsList(ctx).FName(projectId).Execute()

			if err != nil {
				resData, _ := client.ParseResponseData(httpRes)
				return charm.RenderError("Could not retrieve project, unexpected error: " + err.Error() + " with" + fmt.Sprintf(", detail: %v", resData))
			}

			if len(projects.Data) <= 0 {
				return charm.RenderError("unable to find project with name " + projectId)

			}

			project := projects.Data[0]
			_, _, err = qc.ProjectsAPI.ProjectsDelete(ctx, project.Id).Execute()
			if err != nil {
				return charm.RenderError("error deleting qernal project", err)
			}

			printer.PrintResource(charm.RenderWarning("sucessfully deleted project with ID: " + projectId))
			return nil

		},
	}
	cmd.Flags().StringVarP(&projectId, "project", "p", "", "name of the project")
	_ = cmd.MarkFlagRequired("project")
	return cmd
}
