package org

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/qernal/cli-qernal/pkg/helpers"
	"github.com/qernal/cli-qernal/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestCreateOrg(t *testing.T) {

	orgName := helpers.RandomSecretName()
	//set stdout to a buffer we control
	var buf bytes.Buffer
	printer := utils.NewPrinter()
	printer.SetOut(&buf)

	var expectedJson struct {
		ID     string `json:"organisation_id"`
		Name   string `json:"organisation_name"`
		UserID string `json:"user_id"`
	}
	cmd := NewCreateCmd(printer)
	args := []string{"-o", "json", "--name", orgName}
	cmd.SetArgs(args)

	err := cmd.Execute()
	err = json.Unmarshal(buf.Bytes(), &expectedJson)
	if err != nil {
		t.Fatalf("json result is not in expected format %v", err)

	}
	assert.NoError(t, err)

	assert.Equal(t, orgName, expectedJson.Name)

	t.Cleanup(func() {
		helpers.DeleteOrg(expectedJson.ID)

	})

}
