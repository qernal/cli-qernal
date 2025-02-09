package org

import (
	"bytes"
	"testing"

	"github.com/google/uuid"
	"github.com/qernal/cli-qernal/pkg/helpers"
	"github.com/qernal/cli-qernal/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestOrgUpdate(t *testing.T) {
	orgID, _, err := helpers.CreateOrg()
	if err != nil {
		t.Fatalf("failed to create org: %v", err)
	}

	orgName := uuid.NewString()

	var buf bytes.Buffer
	printer := utils.NewPrinter()
	printer.SetOut(&buf)

	// Create root command for persistent flags
	rootCmd := &cobra.Command{Use: "test"}
	rootCmd.PersistentFlags().String("organisation-id", "", "")
	rootCmd.PersistentFlags().String("organisation", "", "")

	updateCmd := NewUpdateCmd(printer)
	rootCmd.AddCommand(updateCmd)

	// Add "update" as first argument since we're using root command
	rootCmd.SetArgs([]string{"update", "--organisation-id", orgID, "--organisation", orgName})

	err = rootCmd.Execute()
	if err != nil {
		t.Fatalf("unable to update to project, command failed with %v", err)
	}

	assert.Contains(t, buf.String(), orgName)

	t.Cleanup(func() {
		helpers.DeleteOrg(orgID)
	})
}
