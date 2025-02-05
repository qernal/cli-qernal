package helpers

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/qernal/cli-qernal/commands/auth"
	"github.com/qernal/cli-qernal/pkg/client"
)

func GetDefaultHost(projid string) (string, error) {
	ctx := context.Background()
	token, err := auth.GetQernalToken()
	if err != nil {
		return "", err
	}

	qc, err := client.New(ctx, nil, nil, token)
	if err != nil {
		return "", err
	}
	resp, r, err := qc.HostsAPI.ProjectsHostsList(context.Background(), projid).Execute()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ProjectsAPI.ProjectsCreate``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)

		return "", err
	}

	for _, host := range resp.Data {
		if host.ReadOnly {
			return host.Host, nil
		}
	}

	return "", errors.New("no default host on project")
}

func GetHostState(disabled bool) string {
	if disabled {
		return "Disabled"
	}
	return "Enabled"
}

func GetCertificateStatus(cert *string) string {
	if cert == nil || *cert == "" {
		return "Not configured"
	}
	return "Configured"
}

func GetReadOnlyStatus(readOnly bool) string {
	if readOnly {
		return "Yes (*.qrnl.app domain)"
	}
	return "No"
}
