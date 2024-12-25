package org

import (
	"bytes"
	"testing"

	"github.com/qernal/cli-qernal/pkg/helpers"
	"github.com/qernal/cli-qernal/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeleteCmd(t *testing.T) {

	orgId, _, err := helpers.CreateOrg()
	if err != nil {
		t.Fatalf("failed to create org: %v", err)
	}

	args := []string{"--id", orgId}

	printer := utils.NewPrinter()
	//set stdout to a buffer we control
	var buf bytes.Buffer
	printer.SetOut(&buf)

	cmd := NewDeleteCmd(printer)
	cmd.SetArgs(args)

	err = cmd.Execute()
	require.NoError(t, err)
	assert.Contains(t, buf.String(), orgId)

}
