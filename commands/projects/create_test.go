package projects

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/qernal/cli-qernal/pkg/helpers"
	"github.com/qernal/cli-qernal/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProjectCreate(t *testing.T) {
	orgID, _, err := helpers.CreateOrg()
	if err != nil {
		t.Fatalf("failed to create org: %v", err)
	}

	projectname := uuid.NewString()

	var buf bytes.Buffer
	printer := utils.NewPrinter()
	printer.SetOut(&buf)

	var expectedJson struct {
		ProjectName    string `json:"project_name"`
		OrganisationID string `json:"organisation_id"`
		ProjectID      string `json:"project_id"`
	}

	// Create root command for persistent flags
	rootCmd := &cobra.Command{Use: "test"}
	rootCmd.PersistentFlags().String("project-id", "", "")
	rootCmd.PersistentFlags().String("project", "", "")
	rootCmd.PersistentFlags().String("organisation-id", "", "")

	testCases := []struct {
		name           string
		args           []string
		expectedOutput string
		expectErr      bool
	}{
		{
			name:           "Valid Project",
			args:           []string{"create", "--name", projectname, "--organisation-id", orgID, "--output", "json"},
			expectedOutput: "project_id",
			expectErr:      false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cmd := NewCreateCmd(printer)
			rootCmd.AddCommand(cmd)
			rootCmd.SetArgs(tc.args)

			err := rootCmd.Execute()
			if tc.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Contains(t, buf.String(), tc.expectedOutput)
				err := json.Unmarshal(buf.Bytes(), &expectedJson)
				require.NoError(t, err, "Failed to unmarshal JSON response")
			}
		})
	}

	t.Cleanup(func() {
		helpers.DeleteOrg(expectedJson.OrganisationID)
	})
}
