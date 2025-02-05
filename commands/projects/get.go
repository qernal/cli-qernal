package projects

import (
	"context"

	"github.com/qernal/cli-qernal/charm"
	"github.com/qernal/cli-qernal/commands/auth"
	"github.com/qernal/cli-qernal/pkg/client"
	"github.com/qernal/cli-qernal/pkg/common"
	"github.com/qernal/cli-qernal/pkg/utils"
	"github.com/spf13/cobra"
)

func NewGetCmd(printer *utils.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "get",
		Aliases: []string{"get"},
		Example: "qernal project get --name <org name>",
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

			name, _ := cmd.Flags().GetString("name")

			project, err := qc.GetProjectByName(name)
			if err != nil {
				return charm.RenderError("", err)
			}

			var data interface{}

			if common.OutputFormat == "json" {
				data = project
			} else {
				data = map[string]interface{}{
					"Name":       project.Name,
					"Project ID": project.Id,
					"Org ID":     project.OrgId,
				}
			}

			if common.OutputFormat == "json" {
				printer.PrintResource(utils.FormatOutput(data, common.OutputFormat))
				return nil
			}

			response := charm.SuccessStyle.Render(utils.FormatOutput(data, common.OutputFormat))
			printer.PrintResource(response)
			return nil

		},
	}
	cmd.Flags().StringVarP(&name, "name", "n", "", "name of the project")
	cmd.Flags().StringVarP(&common.OutputFormat, "output", "o", "text", "output format (json,text)")
	_ = cmd.MarkFlagRequired("name")
	return cmd
}
