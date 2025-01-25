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
		Date   struct {
			CreatedAt string `json:"created_at"`
			UpdatedAt string `json:"updated_at"`
		} `json:"date"`
	}
	cmd := NewCreateCmd(printer)
	args := []string{"-o", "json", "--name", orgName}
	cmd.SetArgs(args)

	cmd.Execute()

	// check if json is as expected
	err := json.Unmarshal(buf.Bytes(), &expectedJson)
	if err != nil {
		t.Fatalf("json result is not in expected format %v", err)
	}

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
