package auth

import (
	"errors"

	"github.com/spf13/cobra"
)

var AuthCmd = &cobra.Command{
	Use:   "auth",
	Short: "Manage your auth tokens",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := cmd.Help()
		if err != nil {
			return err
		}
		return errors.New("a valid subcommand is required")
	},
}

func init() {
	AuthCmd.AddCommand(checkCmd)
	AuthCmd.AddCommand(loginCmd)

}
