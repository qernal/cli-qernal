package org

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/qernal/cli-qernal/pkg/helpers"
	"github.com/qernal/cli-qernal/pkg/utils"
	openapi_chaos_client "github.com/qernal/openapi-chaos-go-client"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestGetOrg(t *testing.T) {
	ID, name, err := helpers.CreateOrg()
	if err != nil {
		t.Fatalf("unable to create org: %v", err)
	}

	var buf bytes.Buffer
	printer := utils.NewPrinter()
	printer.SetOut(&buf)

	// Create root command for persistent flags
	rootCmd := &cobra.Command{Use: "test"}
	rootCmd.PersistentFlags().String("organisation", "", "")

	cmd := NewGetCmd(printer)
	rootCmd.AddCommand(cmd)

	// Add "get" as first argument since we're using root command
	rootCmd.SetArgs([]string{"get", "-o", "json", "--organisation", name})

	var response openapi_chaos_client.OrganisationResponse
	err = rootCmd.Execute()
	if err != nil {
		t.Fatalf("unable to execute command: %v", err)
	}

	// check if json is as expected
	err = json.Unmarshal(buf.Bytes(), &response)
	if err != nil {
		t.Fatalf("json result is not in expected format %v", err)
	}

	assert.NoError(t, err)
	assert.Equal(t, ID, response.Id)

	t.Cleanup(func() {
		helpers.DeleteOrg(ID)
	})
}
