package helpers

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/qernal/cli-qernal/commands/auth"
	"github.com/qernal/cli-qernal/pkg/client"
	"github.com/qernal/cli-qernal/pkg/utils"
	openapi_chaos_client "github.com/qernal/openapi-chaos-go-client"
)

// PaginateSecrets
func PaginateSecrets(printer *utils.Printer, ctx context.Context, qc *client.QernalAPIClient, maxResults int32, projectID string) ([]openapi_chaos_client.SecretMetaResponse, error) {

	initialResp, httpRes, err := qc.SecretsAPI.ProjectsSecretsList(ctx, projectID).Execute()
	if err != nil {
		resData, _ := client.ParseResponseData(httpRes)
		if data, ok := resData.(map[string]interface{}); ok {
			if innerData, ok := data["data"].(map[string]interface{}); ok {
				if nameErr, ok := innerData["name"].(string); ok {
					return nil, printer.RenderError("unable to list projects", errors.New(nameErr))
				}
			}
		}
		printer.Logger.Debug("unable to list projects, request failed",
			slog.String("error", err.Error()),
			slog.Any("response", resData))
		return nil, printer.RenderError("unable to list projects", err)
	}

	allProjects := initialResp.Data
	if initialResp.Meta.Results <= 20 {
		return allProjects, nil
	}

	pageSize := int32(20)
	var currentPage int32
	for currentPage < initialResp.Meta.Pages {
		if maxResults > 0 && len(allProjects) >= int(maxResults) {
			break
		}

		currentPage++
		previous := currentPage - 1
		next, httpRes, err := qc.SecretsAPI.ProjectsSecretsList(ctx, projectID).
			Page(openapi_chaos_client.OrganisationsListPageParameter{
				Size:   &pageSize,
				Before: &previous,
				After:  &currentPage,
			}).Execute()
		if err != nil {
			resData, _ := client.ParseResponseData(httpRes)
			if data, ok := resData.(map[string]interface{}); ok {
				if innerData, ok := data["data"].(map[string]interface{}); ok {
					if nameErr, ok := innerData["name"].(string); ok {
						return nil, printer.RenderError("unable to list projects", errors.New(nameErr))
					}
				}
			}
			printer.Logger.Debug("unable to list projects, request failed",
				slog.String("error", err.Error()),
				slog.Any("response", resData))
			return nil, printer.RenderError("unable to list projects", err)
		}

		allProjects = append(allProjects, next.GetData()...)
	}

	if maxResults > 0 && len(allProjects) > int(maxResults) {
		allProjects = allProjects[:maxResults]
	}

	return allProjects, nil
}

func CreateSecretEnv(projid string, secretname string) (string, string, error) {
	dek, dekRevision, err := FetchDek(projid)
	if err != nil {
		return "", "", err
	}

	encryptedSecret, err := client.EncryptLocalSecret(dek, secretname)
	if err != nil {
		return "", "", err
	}

	ctx := context.Background()
	token, err := auth.GetQernalToken()
	if err != nil {
		return "", "", err
	}

	qc, err := client.New(ctx, nil, nil, token)
	if err != nil {
		return "", "", err
	}

	secretEnvBody := *openapi_chaos_client.NewSecretBody(secretname, openapi_chaos_client.SECRETCREATETYPE_ENVIRONMENT, openapi_chaos_client.SecretCreatePayload{
		SecretEnvironment: &openapi_chaos_client.SecretEnvironment{
			EnvironmentValue: encryptedSecret,
		},
	}, fmt.Sprintf("keys/dek/%d", dekRevision))
	resp, r, err := qc.SecretsAPI.ProjectsSecretsCreate(context.Background(), projid).SecretBody(secretEnvBody).Execute()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ProjectsAPI.ProjectsSecretsCreate``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)

		return "", "", err
	}

	return resp.Name, fmt.Sprintf("projects:%s/%s@%d", projid, resp.Name, resp.Revision), nil
}
