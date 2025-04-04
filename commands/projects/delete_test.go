package projects

import (
	"bytes"
	"testing"

	"github.com/qernal/cli-qernal/pkg/helpers"
	"github.com/qernal/cli-qernal/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeleteCmd(t *testing.T) {
	orgID, _, err := helpers.CreateOrg()
	if err != nil {
		t.Fatalf("failed to create org: %v", err)
	}
	_, projectName, err := helpers.CreateProj(orgID)
	if err != nil {
		t.Fatalf("failed to create org: %v", err)
	}

	args := []string{"--project", projectName}

	printer := utils.NewPrinter()
	//set stdout to a buffer we control
	var buf bytes.Buffer
	printer.SetOut(&buf)

	cmd := NewDeleteCmd(printer)
	cmd.SetArgs(args)

	err = cmd.Execute()
	require.NoError(t, err)
	assert.Contains(t, buf.String(), projectName)

	t.Cleanup(func() {
		helpers.DeleteOrg(orgID)
	})

}
