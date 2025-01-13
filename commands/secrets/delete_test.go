package secrets

import (
	"bytes"
	"strings"
	"testing"

	"github.com/google/uuid"
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
	projectID, _, err := helpers.CreateProj(orgID)
	if err != nil {
		t.Fatalf("failed to create org: %v", err)
	}

	secretName, _, err := helpers.CreateSecretEnv(projectID, strings.ToUpper(uuid.NewString()))
	if err != nil {
		t.Fatalf("failed to create org: %v", err)
	}

	args := []string{"--name", secretName, "--project", projectID}

	printer := utils.NewPrinter()
	//set stdout to a buffer we control
	var buf bytes.Buffer
	printer.SetOut(&buf)

	cmd := NewDeleteCmd(printer)
	cmd.SetArgs(args)

	err = cmd.Execute()
	require.NoError(t, err)
	assert.Contains(t, buf.String(), secretName)

}
