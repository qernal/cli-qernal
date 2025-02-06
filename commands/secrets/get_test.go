package secrets

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

func TestGetSecret(t *testing.T) {
	orgID, _, err := helpers.CreateOrg()
	if err != nil {
		t.Fatalf("unable to create org: %v", err)
	}
	projId, _, err := helpers.CreateProj(orgID)
	if err != nil {
		t.Fatalf("unable to create project: %v", err)
	}
	secretName, _, err := helpers.CreateSecretEnv(projId, helpers.RandomSecretName())
	if err != nil {
		t.Fatalf("unable to create test secret : %v", err)
	}

	// Create root command for persistent flags
	rootCmd := &cobra.Command{Use: "test"}
	rootCmd.PersistentFlags().String("project-id", "", "")
	rootCmd.PersistentFlags().String("project", "", "")

	var buf bytes.Buffer
	printer := utils.NewPrinter()
	printer.SetOut(&buf)

	cmd := NewGetCmd(printer)
	rootCmd.AddCommand(cmd)

	args := []string{"get", "-o", "json", "--project-id", projId, "--name", secretName}
	rootCmd.SetArgs(args)

	var response openapi_chaos_client.SecretMetaResponse
	err = rootCmd.Execute()
	if err != nil {
		t.Fatalf("unable to execute command: %v", err)
	}

	err = json.Unmarshal(buf.Bytes(), &response)
	if err != nil {
		t.Fatalf("json result is not in expected format %v", err)
	}

	assert.NoError(t, err)
	assert.Equal(t, secretName, response.Name)

	t.Cleanup(func() {
		helpers.DeleteOrg(orgID)
	})
}
