package auth

import (
	"context"
	"fmt"

	"github.com/qernal/cli-qernal/charm"
	"github.com/qernal/cli-qernal/pkg/client"
	"github.com/spf13/cobra"
)

var token string

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check if the current authentication token is valid.",
	Long:  "Verify the validity of your authentication token. By default, the command checks the token stored in your configuration file. Use the -t flag to specify a token directly from the command line for validation.",
	RunE: func(cmd *cobra.Command, args []string) error {
		var tokenToUse string
		var err error

		if token != "" {
			tokenToUse = token
		} else {
			// Fallback to token from configuration
			tokenToUse, err = GetQernalToken()
			if err != nil {
				return charm.RenderError("unable to fetch token", err)
			}
		}

		// should fail if token is invalid
		ctx := context.Background()
		_, err = client.New(ctx, tokenToUse)
		if err != nil {
			return charm.RenderError("❌ invalid token, auth check failed with", err)
		}

		// TODO: Display last token refresh
		fmt.Println(charm.SuccessStyle.Render("Token is valid ✅"))
		return nil
	},
}

func init() {
	checkCmd.Flags().StringVarP(&token, "token", "t", "", "Authentication token for validation.")
}
