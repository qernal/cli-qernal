package secrets

import (
	"bytes"
	"testing"

	"github.com/qernal/cli-qernal/pkg/helpers"
	"github.com/qernal/cli-qernal/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeleteCmd(t *testing.T) {
	orgID, _, err := helpers.CreateOrg()
	if err != nil {
		t.Fatalf("failed to create organisation: %v", err)
	}

	projectID, _, err := helpers.CreateProj(orgID)
	if err != nil {
		t.Fatalf("failed to create project: %v", err)
	}

	secretName, _, err := helpers.CreateSecretEnv(projectID, helpers.RandomSecretName())
	if err != nil {
		t.Fatalf("failed to create secret: %v", err)
	}

	// Create a root command to properly handle persistent flags
	rootCmd := &cobra.Command{Use: "test"}
	rootCmd.PersistentFlags().String("project-id", "", "")
	rootCmd.PersistentFlags().String("project", "", "")

	printer := utils.NewPrinter()
	var buf bytes.Buffer
	printer.SetOut(&buf)

	cmd := NewDeleteCmd(printer)
	rootCmd.AddCommand(cmd)

	// Set the args including the project-id flag
	rootCmd.SetArgs([]string{"delete", "--name", secretName, "--project-id", projectID})

	err = rootCmd.Execute()
	require.NoError(t, err)
	assert.Contains(t, buf.String(), secretName)
}
