package org

import (
	"bytes"
	"testing"

	"github.com/google/uuid"
	"github.com/qernal/cli-qernal/pkg/helpers"
	"github.com/qernal/cli-qernal/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestOrgUpdate(t *testing.T) {
	orgID, _, err := helpers.CreateOrg()
	if err != nil {
		t.Fatalf("failed to create org: %v", err)
	}
	orgName := uuid.NewString()

	//set stdout to a buffer we control
	var buf bytes.Buffer
	printer := utils.NewPrinter()
	printer.SetOut(&buf)

	updatecmd := NewUpdateCmd(printer)

	updateArgs := []string{"--id", orgID, "--name", orgName}

	updatecmd.SetArgs(updateArgs)

	err = updatecmd.Execute()
	if err != nil {
		t.Fatalf("unable to update to project, command failed with %v", err)
	}

	assert.Contains(t, buf.String(), orgName)

	t.Cleanup(func() {
		helpers.DeleteOrg(orgID)
	})

}
