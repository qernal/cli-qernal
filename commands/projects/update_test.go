package projects

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/qernal/cli-qernal/pkg/helpers"
	"github.com/qernal/cli-qernal/pkg/utils"
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

	//set stdout to a buffer we control
	var buf bytes.Buffer
	printer := utils.NewPrinter()
	printer.SetOut(&buf)

	createCmd := NewCreateCmd(printer)
	createCmd.SetArgs([]string{"create", "--name", projectname, "--organisation", orgID, "--output", "json"})
	err = createCmd.Execute()
	if err != nil {
		t.Fatalf("unable to create to project, command failed with %v", err)
	}
	err = json.Unmarshal(buf.Bytes(), &expectedJson)
	if err != nil {
		t.Fatalf("unbale to unmarshal output, decode failed with  %v", err)
	}
	buf.Reset()

	updatecmd := NewUpdateCmd(printer)
	// "qernal projects update --project=<project ID> --org <org ID> --name <name>"
	updateArgs := []string{"update", "--project", expectedJson.ProjectID, "--organisation", expectedJson.OrganisationID, "--name", updatedPojectName, "--output", "json"}

	updatecmd.SetArgs(updateArgs)

	err = updatecmd.Execute()
	if err != nil {
		t.Fatalf("unable to update to project, command failed with %v", err)
	}

	assert.Contains(t, buf.String(), updatedPojectName)

	t.Cleanup(func() {
		helpers.DeleteOrg(expectedJson.OrganisationID)
	})

}
