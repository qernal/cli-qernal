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

// CreateOrg returns the ID and name of the created org
func CreateOrg() (string, string, error) {
	token, err := auth.GetQernalToken()
	if err != nil {
		return "", "", err
	}
	ctx := context.Background()
	client, err := client.New(ctx, nil, nil, token)
	if err != nil {
		return "", "", err
	}

	organisationBody := *openapi_chaos_client.NewOrganisationBody(uuid.NewString())
	resp, r, err := client.OrganisationsAPI.OrganisationsCreate(context.Background()).OrganisationBody(organisationBody).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `OrganisationsAPI.OrganisationsCreate``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)

		return "", "", err
	}

	return resp.Id, resp.Name, nil
}

func DeleteOrg(orgid string) {
	token, err := auth.GetQernalToken()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", charm.RenderError("obtaining token failed with:", err).Error())
	}

	ctx := context.Background()
	client, err := client.New(ctx, nil, nil, token)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", charm.RenderError("unable to create qernal client", err).Error())
	}
	_, r, err := client.OrganisationsAPI.OrganisationsDelete(context.Background(), orgid).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", charm.RenderError("unable to create qernal client", err).Error())
		fmt.Fprintf(os.Stderr, "Error when calling `OrganisationsAPI.OrganisationsDelete`: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}

func PaginateOrganisations(printer *utils.Printer, ctx context.Context, qc *client.QernalAPIClient, maxResults int32) ([]openapi_chaos_client.OrganisationResponse, error) {

	initialResp, httpRes, err := qc.OrganisationsAPI.OrganisationsList(ctx).Execute()
	if err != nil {
		resData, _ := client.ParseResponseData(httpRes)
		if data, ok := resData.(map[string]interface{}); ok {
			if innerData, ok := data["data"].(map[string]interface{}); ok {
				if nameErr, ok := innerData["name"].(string); ok {
					return nil, printer.RenderError("unable to list organisations", errors.New(nameErr))
				}
			}
		}
		printer.Logger.Debug("unable to list organisations, request failed",
			slog.String("error", err.Error()),
			slog.Any("response", resData))
		return nil, printer.RenderError("unable to list organisations", err)
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
		next, httpRes, err := qc.OrganisationsAPI.OrganisationsList(ctx).
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
