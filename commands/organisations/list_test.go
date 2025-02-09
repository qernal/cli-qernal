package org

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/qernal/cli-qernal/pkg/utils"
	openapi_chaos_client "github.com/qernal/openapi-chaos-go-client"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestListOrg(t *testing.T) {
	var buf bytes.Buffer
	printer := utils.NewPrinter()
	printer.SetOut(&buf)

	// Create root command for persistent flags
	rootCmd := &cobra.Command{Use: "test"}
	rootCmd.PersistentFlags().String("organisation", "", "")

	var expectedJson []openapi_chaos_client.OrganisationResponse

	cmd := NewOrgListCmd(printer)
	rootCmd.AddCommand(cmd)

	// Add "list" as first argument since we're using root command
	rootCmd.SetArgs([]string{"list", "-o", "json"})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("unable to execute command %v", err)
	}

	err = json.Unmarshal(buf.Bytes(), &expectedJson)
	if err != nil {
		t.Fatalf("json result is not in expected format %v", err)
	}

	assert.NoError(t, err)
	assert.True(t, len(expectedJson) > 0)
}
