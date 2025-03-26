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
)

func TestProjectUpdate(t *testing.T) {
	orgID, _, err := helpers.CreateOrg()
	if err != nil {
		t.Fatalf("failed to create org: %v", err)
	}

	projectname := uuid.NewString()
	var expectedJson struct {
		ProjectName    string `json:"project_name"`
		OrganisationID string `json:"organisation_id"`
		ProjectID      string `json:"project_id"`
	}

	updatedPojectName := uuid.NewString()

	var buf bytes.Buffer
	printer := utils.NewPrinter()
	printer.SetOut(&buf)

	// Create root command for persistent flags
	rootCmd := &cobra.Command{Use: "test"}
	rootCmd.PersistentFlags().String("organisation-id", "", "")

	createCmd := NewCreateCmd(printer)
	rootCmd.AddCommand(createCmd)
	rootCmd.SetArgs([]string{"create", "--name", projectname, "--organisation-id", orgID, "--output", "json"})

	err = rootCmd.Execute()
	if err != nil {
		t.Fatalf("unable to create to project, command failed with %v", err)
	}

	err = json.Unmarshal(buf.Bytes(), &expectedJson)
	if err != nil {
		t.Fatalf("unable to unmarshal output, decode failed with %v", err)
	}

	buf.Reset()
	rootCmd = &cobra.Command{Use: "test"}
	rootCmd.PersistentFlags().String("organisation-id", "", "")

	updateCmd := NewUpdateCmd(printer)
	rootCmd.AddCommand(updateCmd)
	rootCmd.SetArgs([]string{
		"update",
		"--project", expectedJson.ProjectID,
		"--organisation-id", expectedJson.OrganisationID,
		"--name", updatedPojectName,
		"--output", "json",
	})

	err = rootCmd.Execute()
	if err != nil {
		t.Fatalf("unable to update to project, command failed with %v", err)
	}

	assert.Contains(t, buf.String(), updatedPojectName)

	t.Cleanup(func() {
		helpers.DeleteOrg(expectedJson.OrganisationID)
	})
}
