package secrets

import (
	"errors"

	"github.com/qernal/cli-qernal/pkg/utils"
	"github.com/spf13/cobra"
)

var projectID string
var secretName string

var SecretsCmd = &cobra.Command{
	Use:     "secrets",
	Short:   "Manage your secrets",
	Aliases: []string{"secret"},
	RunE: func(cmd *cobra.Command, args []string) error {
		err := cmd.Help()
		if err != nil {
			return err
		}
		return errors.New("a valid subcommand is required")
	},
}

func init() {
	printer := utils.NewPrinter()
	SecretsCmd.AddCommand(SecretsListCmd)
	SecretsCmd.AddCommand(NewEncryptCmd(printer))
	SecretsCmd.AddCommand(NewDeleteCmd(printer))
	SecretsCmd.AddCommand(NewCreateCmd(printer))
	SecretsCmd.PersistentFlags().StringVarP(&projectID, "project", "p", "", "ID of the project")
	SecretsCmd.PersistentFlags().StringVarP(&secretName, "name", "n", "", "name of the secret")

	_ = NewCreateCmd(printer).MarkFlagRequired("name")
	_ = NewCreateCmd(printer).MarkFlagRequired("project")
	_ = NewCreateCmd(printer).MarkFlagRequired("public-key")
	_ = NewCreateCmd(printer).MarkFlagRequired("private-key")

	_ = NewDeleteCmd(printer).MarkFlagRequired("name")
	_ = NewDeleteCmd(printer).MarkFlagRequired("project")
}
