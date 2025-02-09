package org

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
	_, name, err := helpers.CreateOrg()
	if err != nil {
		t.Fatalf("failed to create org: %v", err)
	}

	// Create root command for persistent flags
	rootCmd := &cobra.Command{Use: "test"}
	rootCmd.PersistentFlags().String("organisation", "", "")

	printer := utils.NewPrinter()
	var buf bytes.Buffer
	printer.SetOut(&buf)

	cmd := NewDeleteCmd(printer)
	rootCmd.AddCommand(cmd)

	// Add "delete" as first argument since we're using root command
	rootCmd.SetArgs([]string{"delete", "--organisation", name})

	err = rootCmd.Execute()
	require.NoError(t, err)

	t.Log(buf.String())
	assert.Contains(t, buf.String(), "deleted")
}
