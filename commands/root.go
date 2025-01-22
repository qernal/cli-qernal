package commands

import (
	"fmt"
	"os"

	"github.com/qernal/cli-qernal/commands/auth"
	org "github.com/qernal/cli-qernal/commands/organisations"
	"github.com/qernal/cli-qernal/commands/projects"
	"github.com/qernal/cli-qernal/commands/secrets"
	"github.com/qernal/cli-qernal/pkg/build"
	"github.com/qernal/cli-qernal/pkg/common"
	"github.com/spf13/cobra"
)

var version bool

var RootCmd = &cobra.Command{
	Use:          "qernal",
	Short:        fmt.Sprintf("CLI for interacting with Qernal\nVersion: %s", build.Version),
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if version {
			versionCmd.Run(cmd, args)
		} else {
			cmd.Help()
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
	RootCmd.AddCommand(auth.AuthCmd)
	RootCmd.AddCommand(secrets.SecretsCmd)
	RootCmd.AddCommand(projects.ProjectsCmd)
	RootCmd.AddCommand(org.OrgCmd)

}
