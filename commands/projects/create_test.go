package projects

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/qernal/cli-qernal/pkg/helpers"
	"github.com/qernal/cli-qernal/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProjectCreate(t *testing.T) {
	orgID, _, err := helpers.CreateOrg()
	if err != nil {
		t.Fatalf("failed to create org: %v", err)
	}
	projectname := uuid.NewString()

	//set stdout to a buffer we control
	var buf bytes.Buffer
	printer := utils.NewPrinter()
	printer.SetOut(&buf)

	var expectedJson struct {
		ProjectName    string `json:"project_name"`
		OrganisationID string `json:"organisation_id"`
		ProjectID      string `json:"project_id"`
	}

	testCases := []struct {
		name           string
		args           []string
		expectedOutput string
		expectErr      bool
	}{
		{
			name:           "Zero Args",
			args:           []string{"create", "--organisation", orgID},
			expectedOutput: "required flag(s) not set",
			expectErr:      true,
		},
		{
			name:           "Valid Project",
			args:           []string{"create", "--name", projectname, "--organisation", orgID, "--output", "json"},
			expectedOutput: "project_id",
			expectErr:      false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cmd := NewCreateCmd(printer)
			cmd.SetArgs(tc.args)
			err := cmd.Execute()
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
