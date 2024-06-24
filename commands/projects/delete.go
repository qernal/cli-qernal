package projects

import (
	"context"
	"fmt"

	"github.com/qernal/cli-qernal/charm"
	"github.com/qernal/cli-qernal/commands/auth"
	"github.com/qernal/cli-qernal/pkg/client"
	"github.com/spf13/cobra"
)

var projectID string
var DeleteCmd = &cobra.Command{
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

		qc, err := client.New(ctx, token)
		if err != nil {
			return charm.RenderError("error creating qernal client", err)
		}

		_, _, err = qc.ProjectsAPI.ProjectsDelete(ctx, projectID).Execute()
		if err != nil {
			return charm.RenderError("error creating qernal client", err)
		}

		fmt.Println(charm.RenderWarning("sucessfully deleted project with ID: " + projectID))

		return nil

	},
}

func init() {
	DeleteCmd.Flags().StringVarP(&projectID, "project", "p", "", "Project ID")
	DeleteCmd.MarkFlagRequired("project")
}
