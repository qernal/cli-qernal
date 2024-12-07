package projects

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"

	"github.com/qernal/cli-qernal/charm"
	utils "github.com/qernal/cli-qernal/pkg/uitls"
	"github.com/stretchr/testify/assert"
	"github.com/teris-io/shortid"
)

func TestProjectUpdate(t *testing.T) {
	orgID := os.Getenv("QERNAL_TEST_ORG")
	if orgID == "" {
		t.Fatal(charm.ErrorStyle.Render("qernal test org is not set"))

	}
	var expectedJson struct {
		ProjectName    string `json:"project_name"`
		OrganisationID string `json:"organisation_id"`
		ProjectID      string `json:"project_id"`
	}
	projectname, _ := shortid.Generate()
	updatedPojectName, _ := shortid.Generate()

	//set stdout to a buffer we control
	var buf bytes.Buffer
	printer := utils.NewPrinter()
	printer.SetOut(&buf)

	createCmd := NewCreateCmd(printer)
	createCmd.SetArgs([]string{"create", "--name", projectname, "--organisation", orgID, "--output", "json"})
	err := createCmd.Execute()
	if err != nil {
		t.Fatalf("unable to create to project, command failed with %v", err)
	}
	err = json.Unmarshal(buf.Bytes(), &expectedJson)
	if err != nil {
		t.Fatalf("unbale to unmarshal output, decode failed with  %v", err)
	}
	buf.Reset()

	updatecmd := NewupdateCmd(printer)
	// "qernal projects update --project=<project ID> --org <org ID> --name <name>"
	updateArgs := []string{"update", "--project", expectedJson.ProjectID, "--organisation", expectedJson.OrganisationID, "--name", updatedPojectName, "--output", "json"}

	updatecmd.SetArgs(updateArgs)

	err = updatecmd.Execute()
	if err != nil {
		t.Fatalf("unable to update to project, command failed with %v", err)
	}

	assert.Contains(t, buf.String(), updatedPojectName)

	t.Cleanup(func() {
		cmd := NewDeleteCmd(printer)
		cmd.SetArgs([]string{"--project", updatedPojectName})
		err := cmd.Execute()
		if err != nil {
			t.Fatalf("unable to delete project with name %s: %v", updatedPojectName, err)
		}
	})

}
