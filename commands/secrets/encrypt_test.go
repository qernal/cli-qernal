package secrets

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/qernal/cli-qernal/pkg/helpers"
	"github.com/qernal/cli-qernal/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEncryptCmd(t *testing.T) {
	plaintextValue := "reallyrealvalue"

	outputPrinter := utils.NewPrinter()

	orgID, _, err := helpers.CreateOrg()
	if err != nil {
		t.Fatalf("failed to create organization: %v", err)
	}
	projectID, _, err := helpers.CreateProj(orgID)
	if err != nil {
		t.Fatalf("failed to create project: %v", err)
	}

	commandArgs := []string{"--project", projectID, "--output", "json"}

	// Set stdout and stdin to controlled buffers
	var outputBuffer bytes.Buffer
	var inputBuffer bytes.Buffer
	outputPrinter.SetOut(&outputBuffer)

	_, err = inputBuffer.WriteString(plaintextValue + "\n")
	if err != nil {
		t.Fatalf("failed to write to input buffer: %v", err)
	}

	var encryptionOutput struct {
		RevisionID     int32  `json:"revision_id"`
		EncryptedValue string `json:"encrypted_value"`
	}

	encryptCmd := NewEncryptCmd(outputPrinter)
	encryptCmd.SetArgs(commandArgs)
	encryptCmd.SetIn(&inputBuffer)

	err = encryptCmd.Execute()
	require.NoError(t, err)

	err = json.Unmarshal(outputBuffer.Bytes(), &encryptionOutput)
	require.NoError(t, err)

	assert.Len(t, encryptionOutput.EncryptedValue, 84)

	t.Cleanup(func() {
		helpers.DeleteOrg(orgID)
	})
}
