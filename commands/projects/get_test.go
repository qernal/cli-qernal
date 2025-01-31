package projects

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/qernal/cli-qernal/pkg/helpers"
	"github.com/qernal/cli-qernal/pkg/utils"
	openapi_chaos_client "github.com/qernal/openapi-chaos-go-client"
	"github.com/stretchr/testify/assert"
)

func TestGetProj(t *testing.T) {
	orgID, _, err := helpers.CreateOrg()
	if err != nil {
		t.Fatalf("unable to create org: %v", err)
	}

	projId, projName, err := helpers.CreateProj(orgID)
	if err != nil {
		t.Fatalf("unable to create project: %v", err)
	}

	var buf bytes.Buffer
	printer := utils.NewPrinter()
	printer.SetOut(&buf)

	cmd := NewGetCmd(printer)
	args := []string{"-o", "json", "--name", projName}
	cmd.SetArgs(args)

	var response openapi_chaos_client.ProjectResponse

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

	assert.Equal(t, projId, response.Id)
	assert.Equal(t, projName, response.Name)

	t.Cleanup(func() {
		helpers.DeleteOrg(orgID)
	})
}
