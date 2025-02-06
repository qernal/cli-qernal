package secrets

import (
	"bytes"
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/qernal/cli-qernal/pkg/helpers"
	"github.com/qernal/cli-qernal/pkg/utils"
	"github.com/spf13/cobra"
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

	// Create root command for persistent flags
	rootCmd := &cobra.Command{Use: "test"}
	rootCmd.PersistentFlags().String("project-id", "", "")
	rootCmd.PersistentFlags().String("project", "", "")

	commandArgs := []string{"encrypt", "--project-id", projectID, "--output", "json"}

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
	rootCmd.AddCommand(encryptCmd)
	rootCmd.SetArgs(commandArgs)
	encryptCmd.SetIn(&inputBuffer)

	err = rootCmd.Execute()
	require.NoError(t, err)

	err = json.Unmarshal(outputBuffer.Bytes(), &encryptionOutput)
	require.NoError(t, err)
	assert.Len(t, encryptionOutput.EncryptedValue, 84)

	t.Cleanup(func() {
		helpers.DeleteOrg(orgID)
	})
}

func TestEncryptCmdWithFile(t *testing.T) {
	filePath := "/tmp/" + utils.GenerateRandomString(6) + ".txt"
	randomStrings := make([]string, 10)
	for i := range randomStrings {
		randomStrings[i] = utils.GenerateRandomString(60)
	}
	fileContent := []byte(strings.Join(randomStrings, "\n"))
	err := os.WriteFile(filePath, fileContent, 0644)
	require.NoError(t, err)

	outputPrinter := utils.NewPrinter()
	orgID, _, err := helpers.CreateOrg()
	if err != nil {
		t.Fatalf("failed to create organization: %v", err)
	}
	projectID, _, err := helpers.CreateProj(orgID)
	if err != nil {
		t.Fatalf("failed to create project: %v", err)
	}

	// Create root command for persistent flags
	rootCmd := &cobra.Command{Use: "test"}
	rootCmd.PersistentFlags().String("project-id", "", "")
	rootCmd.PersistentFlags().String("project", "", "")

	commandArgs := []string{"encrypt", "--project-id", projectID, "--output", "json"}

	var outputBuffer bytes.Buffer
	outputPrinter.SetOut(&outputBuffer)

	var encryptionOutput struct {
		RevisionID     int32  `json:"revision_id"`
		EncryptedValue string `json:"encrypted_value"`
	}

	inputBuffer, err := os.ReadFile(filePath)
	require.NoError(t, err)

	encryptCmd := NewEncryptCmd(outputPrinter)
	rootCmd.AddCommand(encryptCmd)
	rootCmd.SetArgs(commandArgs)
	encryptCmd.SetIn(bytes.NewReader(inputBuffer))

	err = rootCmd.Execute()
	require.NoError(t, err)

	err = json.Unmarshal(outputBuffer.Bytes(), &encryptionOutput)
	require.NoError(t, err)
	println(len(encryptionOutput.EncryptedValue))
	assert.Greater(t, len(encryptionOutput.EncryptedValue), 200)

	t.Cleanup(func() {
		helpers.DeleteOrg(orgID)
		os.Remove(filePath)
	})
}
