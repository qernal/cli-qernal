package projects

import (
	"bytes"
	"testing"

	"github.com/google/uuid"
	"github.com/qernal/cli-qernal/pkg/helpers"
	"github.com/qernal/cli-qernal/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestProjectUpdate(t *testing.T) {
	orgId, _, err := helpers.CreateOrg()
	if err != nil {
		t.Fatalf("failed to create org: %v", err)
	}

	projectId, _, err := helpers.CreateProj(orgId)
	if err != nil {
		t.Fatalf("failed to create project: %v", err)
	}

	newPojectName := uuid.NewString()

	//set stdout to a buffer we control
	var buf bytes.Buffer
	printer := utils.NewPrinter()
	printer.SetOut(&buf)

	updatecmd := NewUpdateCmd(printer)
	// "qernal projects update --project=<project ID> --org <org ID> --name <name>"
	updateArgs := []string{"update", "--project", projectId, "--organisation", orgId, "--name", newPojectName, "--output", "json"}

	updatecmd.SetArgs(updateArgs)

	err = updatecmd.Execute()
	if err != nil {
		t.Fatalf("unable to update to project, command failed with %v", err)
	}

	assert.Contains(t, buf.String(), newPojectName)

	t.Cleanup(func() {
		helpers.DeleteOrg(orgId)
	})

}
