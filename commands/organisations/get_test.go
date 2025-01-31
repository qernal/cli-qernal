package org

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/qernal/cli-qernal/pkg/helpers"
	"github.com/qernal/cli-qernal/pkg/utils"
	openapi_chaos_client "github.com/qernal/openapi-chaos-go-client"
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

	cmd := NewGetCmd(printer)
	args := []string{"-o", "json", "--name", name}
	cmd.SetArgs(args)

	var response openapi_chaos_client.OrganisationResponse

	err = cmd.Execute()
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
