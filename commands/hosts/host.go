package hosts

import (
	"errors"

	"github.com/qernal/cli-qernal/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	hostName    string
	projectName string
	certRef     string
	diabled     bool
	publicIPV4  = "45.133.240.10"
	publicIPV6  = "2a13:2b00:1::1"
)
var HostCmd = &cobra.Command{
	Use:     "hosts",
	Short:   "Manage your qernal hosts",
	Aliases: []string{"host", "host"},
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
	HostCmd.AddCommand(NewListCmd(printer))
	HostCmd.AddCommand(NewCreateCmd(printer))
	HostCmd.AddCommand(NewDeleteCmd(printer))
	HostCmd.AddCommand(NewVerifyCmd(printer))

	HostCmd.AddCommand(NewUpdateCmd(printer))
	HostCmd.AddCommand(NewGetCmd(printer))
	HostCmd.PersistentFlags().StringVarP(&projectName, "project", "p", "", "project to associate this host with")
	HostCmd.PersistentFlags().StringVarP(&hostName, "name", "n", "", "name of the host")
	HostCmd.PersistentFlags().StringVar(&certRef, "cert", "", "name of the secret storing the TLS certificate - the secret must be of type 'certificate'")
	HostCmd.PersistentFlags().BoolVarP(&diabled, "disable", "e", false, "hosts are routable by default, setting this to false will disable this host")

}
