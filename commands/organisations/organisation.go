package org

import (
	"errors"

	"github.com/qernal/cli-qernal/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	orgName string
	orgID   string
)
var OrgCmd = &cobra.Command{
	Use:     "organisation",
	Short:   "Manage your qernal organisations",
	Aliases: []string{"org", "organisations"},
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
	OrgCmd.AddCommand(NewOrgListCmd(printer))
	OrgCmd.AddCommand(NewCreateCmd(printer))
	OrgCmd.AddCommand(NewDeleteCmd(printer))
	OrgCmd.AddCommand(NewUpdateCmd(printer))
	OrgCmd.AddCommand(NewGetCmd(printer))

}
