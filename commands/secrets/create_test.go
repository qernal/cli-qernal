package secrets

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"bytes"

	"github.com/google/uuid"
	"github.com/qernal/cli-qernal/pkg/helpers"
	"github.com/qernal/cli-qernal/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSceretCreate(t *testing.T) {
	orgID, _, err := helpers.CreateOrg()
	if err != nil {
		t.Fatalf("failed to create org: %v", err)
	}
	projectID, _, err := helpers.CreateProj(orgID)
	if err != nil {
		t.Fatalf("failed to create org: %v", err)
	}

	secretName := uuid.NewString()
	secretValue := strings.ToUpper(uuid.NewString())
	printer := utils.NewPrinter()
	//set stdout to a buffer we control
	var buf bytes.Buffer
	printer.SetOut(&buf)

	testCases := []struct {
		name           string
		args           []string
		expectedOutput string
		expectErr      bool
	}{

		{
			name:           "Invalid Secret Type",
			args:           []string{"create", "--name", secretName, orgID, "--type", "invalidtype"},
			expectedOutput: "Invalid secret type. Must be one of ('registry', 'environment', or 'certificate')",
			expectErr:      true,
		},
		{
			args:           []string{"--name", secretName, "--project", projectID, "--type", "environment"},
			name:           "Valid Environment Secret",
			expectedOutput: "created environment secret with name",
			expectErr:      false,
		},
		{
			args:           []string{"--name", secretName, "--project", projectID, "--type", "registry", "--registry-url", "docker.io"},
			name:           "Valid Registry Secret",
			expectedOutput: "created environment secret with name",
			expectErr:      false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cmd := NewCreateCmd(printer)

			/// randomize secret name on each run to avoid a 409
			tc.args[1] = uuid.NewString()
			cmd.SetArgs(tc.args)

			// input buffer for stdin
			var inputBuf bytes.Buffer
			inputBuf.WriteString(secretValue + "\n")
			cmd.SetIn(&inputBuf)
			err := cmd.Execute()
			if tc.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Contains(t, buf.String(), tc.expectedOutput)

			}
		})
	}
	t.Cleanup(func() {
		helpers.DeleteOrg(orgID)
		helpers.DeleteProj(projectID)
	})
}

func TestCertCreate(t *testing.T) {
	orgID, _, err := helpers.CreateOrg()
	if err != nil {
		t.Fatalf("failed to create org: %v", err)
	}
	projectID, _, err := helpers.CreateProj(orgID)
	if err != nil {
		t.Fatalf("failed to create org: %v", err)
	}

	pubKey, privKey, err := helpers.GenerateSelfSignedCert()
	require.NoError(t, err, "failed to generate self-signed cert")

	tempDir := t.TempDir()
	certFilePath := filepath.Join(tempDir, "cert.pem")
	keyFilePath := filepath.Join(tempDir, "key.pem")

	// Write the certificate and private key to files
	err = os.WriteFile(certFilePath, pubKey, 0600)
	require.NoError(t, err, "failed to write certificate to file")

	err = os.WriteFile(keyFilePath, privKey, 0600)
	require.NoError(t, err, "failed to write private key to file")

	secretName := uuid.NewString()
	printer := utils.NewPrinter()
	//set stdout to a buffer we control
	var buf bytes.Buffer
	printer.SetOut(&buf)

	testCases := []struct {
		name           string
		args           []string
		expectedOutput string
		expectErr      bool
	}{

		{
			name:           "Valid Certificate Secret",
			args:           []string{"--name", secretName, "--type", "certificate", "--public-key", certFilePath, "--private-key", keyFilePath, "--project", projectID},
			expectedOutput: "Created certificate secret with name ",
			expectErr:      false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cmd := NewCreateCmd(printer)
			cmd.SetArgs(tc.args)
			err := cmd.Execute()
			if tc.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Contains(t, buf.String(), tc.expectedOutput)

			}
		})
	}

	t.Cleanup(func() {
		helpers.DeleteOrg(orgID)
		helpers.DeleteProj(projectID)
		os.Remove(certFilePath)
		os.Remove(keyFilePath)
	})

}
