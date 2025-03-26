package org

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	"github.com/qernal/cli-qernal/commands/auth"
	"github.com/qernal/cli-qernal/pkg/client"
	"github.com/qernal/cli-qernal/pkg/helpers"
	"github.com/qernal/cli-qernal/pkg/utils"
	openapi_chaos_client "github.com/qernal/openapi-chaos-go-client"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestCreateOrg(t *testing.T) {
	orgName := helpers.RandomSecretName()

	// Set up printer and buffer
	var buf bytes.Buffer
	printer := utils.NewPrinter()
	printer.SetOut(&buf)

	var response openapi_chaos_client.OrganisationResponse

	// Create root command for persistent flags
	rootCmd := &cobra.Command{Use: "test"}
	rootCmd.PersistentFlags().String("organisation", "", "")

	// Create and add the create command
	cmd := NewCreateCmd(printer)
	rootCmd.AddCommand(cmd)

	// Set args on root command
	rootCmd.SetArgs([]string{"create", "-o", "json", "--organisation", orgName})

	// Execute the root command
	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("unable to execute command: %v", err)
	}

	// Parse response
	err = json.Unmarshal(buf.Bytes(), &response)
	if err != nil {
		t.Fatalf("json result is not in expected format %v", err)
	}

	// Verify organization was created
	ctx := context.Background()
	token, err := auth.GetQernalToken()
	if err != nil {
		t.Fatalf("unable to obtain auth token %v", err)
	}

	qc, err := client.New(ctx, nil, nil, token)
	if err != nil {
		t.Fatalf("unable to create qernal client %v", err)
	}

	org, err := qc.GetOrgByName(orgName)
	if err != nil {
		t.Fatalf("unable to find organisation %v", err)
	}

	assert.NoError(t, err)
	assert.Equal(t, orgName, org.Name)

	t.Cleanup(func() {
		helpers.DeleteOrg(org.Id)
	})
}
