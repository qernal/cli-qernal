package secrets

import (
	"errors"

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
	SecretsCmd.AddCommand(SecretsListCmd)
	SecretsCmd.AddCommand(EncryptCmd)
	SecretsCmd.AddCommand(DeleteCmd)
	SecretsCmd.AddCommand(CreateCmd)
	SecretsCmd.PersistentFlags().StringVarP(&projectID, "project", "p", "", "ID of the project")
	SecretsCmd.PersistentFlags().StringVarP(&secretName, "name", "n", "", "name of the secret")

}
