package commands

import (
	"fmt"
	"runtime"

	"github.com/qernal/cli-qernal/pkg/build"
	"github.com/spf13/cobra"
)

var (
	verbose bool
)
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Output CLI version and build info",
	Run: func(cmd *cobra.Command, args []string) {
		switch {
		case verbose:
			fmt.Printf("Client version: %s\n", build.Version)
			fmt.Printf("Build date (client): %s\n", build.Date)
			fmt.Printf("Git commit (client): %s\n", build.Commit)
			fmt.Printf("OS/Arch (client): %s/%s\n", runtime.GOOS, runtime.GOARCH)
		default:
			fmt.Printf("Qernal CLI %s\n", build.Version)
		}
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
	versionCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Use verbose output to see full information")
}
