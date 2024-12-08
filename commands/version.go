package commands

import (
	"fmt"
	"runtime"

	"github.com/qernal/cli-qernal/pkg/build"
	"github.com/qernal/cli-qernal/pkg/utils"

	"github.com/qernal/cli-qernal/pkg/common"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Output CLI version and build info",
	RunE: func(cmd *cobra.Command, args []string) error {

		buildInfo := struct {
			Version string `json:"version"`
			Commit  string `json:"commit"`
			Date    string `json:"date"`
			OS      string `json:"os"`
			Arch    string `json:"arch"`
		}{
			Version: build.Version,
			Commit:  build.Commit,
			Date:    build.Date,
			OS:      runtime.GOOS,
			Arch:    runtime.GOARCH,
		}

		if common.OutputFormat == "json" {
			fmt.Println(utils.FormatOutput(buildInfo, common.OutputFormat))
			return nil
		}
		fmt.Printf("Client version: %s\n", build.Version)
		fmt.Printf("Build date (client): %s\n", build.Date)
		fmt.Printf("Git commit (client): %s\n", build.Commit)
		fmt.Printf("OS/Arch (client): %s/%s\n", runtime.GOOS, runtime.GOARCH)

		return nil
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
