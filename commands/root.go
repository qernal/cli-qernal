package commands

import (
	"fmt"
	"os"

	"github.com/qernal/cli-qernal/commands/auth"
	"github.com/qernal/cli-qernal/commands/functions"
	"github.com/qernal/cli-qernal/commands/hosts"

	org "github.com/qernal/cli-qernal/commands/organisations"
	"github.com/qernal/cli-qernal/commands/projects"
	"github.com/qernal/cli-qernal/commands/providers"
	"github.com/qernal/cli-qernal/commands/secrets"
	"github.com/qernal/cli-qernal/pkg/build"
	"github.com/qernal/cli-qernal/pkg/common"
	"github.com/spf13/cobra"
)

var (
	version    bool
	maxResults int32
	projectID  string
	project    string
	orgID      string
	orgName    string
)
var RootCmd = &cobra.Command{
	Use:          "qernal",
	Short:        fmt.Sprintf("CLI for interacting with Qernal\nVersion: %s", build.Version),
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if version {
			versionCmd.Run(cmd, args)
		} else {
			if err := cmd.Help(); err != nil {
				return fmt.Errorf("failed to display help: %w", err)
			}
		}
		return nil
	},
}

func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	RootCmd.Flags().BoolVarP(&version, "version", "v", false, "Print the version of the CLI")
	RootCmd.PersistentFlags().StringVarP(&common.OutputFormat, "output", "o", "text", "output format (json,text)")
	RootCmd.PersistentFlags().Int32Var(&maxResults, "max", 0, "Maximum number of results to return, defaults to all")
	RootCmd.PersistentFlags().StringVar(&projectID, "project-id", "", "ID of the project")
	RootCmd.PersistentFlags().StringVar(&project, "project", "", "name of the project")
	RootCmd.PersistentFlags().StringVar(&orgID, "organisation-id", "", "Organisation ID")
	RootCmd.PersistentFlags().StringVar(&orgName, "organisation", "", "name of the organisation")
	RootCmd.AddCommand(auth.AuthCmd)
	RootCmd.AddCommand(secrets.SecretsCmd)
	RootCmd.AddCommand(projects.ProjectsCmd)
	RootCmd.AddCommand(functions.FunctionCmd)
	RootCmd.AddCommand(org.OrgCmd)
	RootCmd.AddCommand(hosts.HostCmd)
	RootCmd.AddCommand(providers.ProvidersCmd)
}
