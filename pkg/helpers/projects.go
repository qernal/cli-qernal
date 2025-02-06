package helpers

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/google/uuid"
	"github.com/qernal/cli-qernal/charm"
	"github.com/qernal/cli-qernal/commands/auth"
	"github.com/qernal/cli-qernal/pkg/client"
	"github.com/qernal/cli-qernal/pkg/utils"
	openapi_chaos_client "github.com/qernal/openapi-chaos-go-client"
)

// CreateProj returns the ID and name of the created project
func CreateProj(orgid string) (string, string, error) {
	token, err := auth.GetQernalToken()
	if err != nil {
		return "", "", err
	}
	ctx := context.Background()
	client, err := client.New(ctx, nil, nil, token)
	if err != nil {
		return "", "", err
	}

	projectBody := *openapi_chaos_client.NewProjectBody(orgid, uuid.NewString())
	resp, r, err := client.ProjectsAPI.ProjectsCreate(context.Background()).ProjectBody(projectBody).Execute()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ProjectsAPI.ProjectsCreate``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)

		return "", "", err
	}

	return resp.Id, resp.Name, nil
}

func DeleteProj(projid string) {
	token, err := auth.GetQernalToken()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", charm.RenderError("obtaining token failed with:", err).Error())
	}

	ctx := context.Background()
	client, err := client.New(ctx, nil, nil, token)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", charm.RenderError("unable to create qernal client", err).Error())
	}

	_, r, err := client.ProjectsAPI.ProjectsDelete(context.Background(), projid).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ProjectsAPI.ProjectsDelete``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}

func PaginateProjects(printer *utils.Printer, ctx context.Context, qc *client.QernalAPIClient, maxResults int32) ([]openapi_chaos_client.ProjectResponse, error) {

	initialResp, httpRes, err := qc.ProjectsAPI.ProjectsList(ctx).Execute()
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
	allProjects := initialResp.GetData()
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
		next, httpRes, err := qc.ProjectsAPI.ProjectsList(ctx).
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
