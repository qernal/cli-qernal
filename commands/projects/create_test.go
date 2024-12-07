package projects

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"

	"github.com/qernal/cli-qernal/charm"
	utils "github.com/qernal/cli-qernal/pkg/uitls"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/teris-io/shortid"
)

func TestProjectCreate(t *testing.T) {
	orgID := os.Getenv("QERNAL_TEST_ORG")
	if orgID == "" {
		t.Fatal(charm.ErrorStyle.Render("qernal test org is not set"))

	}
	projectname, _ := shortid.Generate()

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
				json.Unmarshal(buf.Bytes(), &expectedJson)

			}
		})
	}
	t.Cleanup(func() {
		cmd := NewDeleteCmd(printer)
		cmd.SetArgs([]string{"--project", projectname})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("unable to delete project with name %s: %v", projectname, err)
		}
	})
}
