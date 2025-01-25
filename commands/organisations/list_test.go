package org

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/qernal/cli-qernal/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestListOrg(t *testing.T) {

	//set stdout to a buffer we control
	var buf bytes.Buffer
	printer := utils.NewPrinter()
	printer.SetOut(&buf)

	var expectedJson []struct {
		Date struct {
			CreatedAt string `json:"created_at"`
			UpdatedAt string `json:"updated_at"`
		} `json:"date"`
		ID     string `json:"id"`
		Name   string `json:"name"`
		UserID string `json:"user_id"`
	}
	cmd := NewOrgListCmd(printer)
	cmd.SetArgs([]string{"-o", "json"})

	cmd.Execute()
	err := json.Unmarshal(buf.Bytes(), &expectedJson)
	if err != nil {
		t.Fatalf("json result is not in expected format %v", err)

	}

	assert.NoError(t, err)

	assert.True(t, len(expectedJson) > 0)
}
