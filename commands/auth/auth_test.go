package auth_test

import (
	"os"
	"testing"

	"github.com/qernal/cli-qernal/commands/auth"
	"github.com/stretchr/testify/assert"
)

func TestValidateToken(t *testing.T) {
	t.Run("Test Invalid token", func(t *testing.T) {
		//TODO expand test case
		token := "idjdkdddd@"
		err := auth.ValidateToken(token)
		assert.ErrorContains(t, err, "invalid token format")
	})
}

// validate that tokens from environment variables are being respected
func TestEnvTokens(t *testing.T) {
	expectedToken := "sokaodkadokad@d0kdoksl"
	os.Setenv("QERNAL_TOKEN", expectedToken)

	token, err := auth.GetQernalToken()
	if err != nil {
		t.Fatalf("test failed with error :%s", err)
	}
	assert.Equal(t, expectedToken, token)
}

//TODO test that the token in the config file is always the same after read
