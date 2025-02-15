package helpers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log/slog"
	"strings"

	"github.com/qernal/cli-qernal/pkg/client"
	"github.com/qernal/cli-qernal/pkg/utils"
	openapi_chaos_client "github.com/qernal/openapi-chaos-go-client"
	"sigs.k8s.io/yaml"
)

// ParseFunctionConfig reads and parses the function configurations from a file
func ParseFunctionConfig(file string, printer *utils.Printer) ([]openapi_chaos_client.FunctionBody, error) {
	// Read the file content
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("error reading function file: %w", err)
	}

	// Split the content by YAML document separator
	documents := strings.Split(string(content), "---")
	var functions []openapi_chaos_client.FunctionBody

	for _, doc := range documents {
		if strings.TrimSpace(doc) == "" {
			continue
		}

		var rawData interface{}
		if err := yaml.Unmarshal([]byte(doc), &rawData); err != nil {
			return nil, fmt.Errorf("error parsing YAML: %w", err)
		}

		jsonData, err := json.Marshal(rawData)
		if err != nil {
			return nil, fmt.Errorf("error converting to JSON: %w", err)
		}

		var config openapi_chaos_client.FunctionBody
		if err := json.Unmarshal(jsonData, &config); err != nil {
			return nil, fmt.Errorf("error parsing JSON to config: %w", err)
		}

		functions = append(functions, config)
	}

	return functions, nil
}

func PaginateFunctions(printer *utils.Printer, ctx context.Context, qc *client.QernalAPIClient, maxResults int32, projectId string) ([]openapi_chaos_client.Function, error) {

	initialResp, httpRes, err := qc.FunctionsAPI.ProjectsFunctionsList(ctx, projectId).Execute()
	if err != nil {
		resData, _ := client.ParseResponseData(httpRes)
		if data, ok := resData.(map[string]interface{}); ok {
			if innerData, ok := data["data"].(map[string]interface{}); ok {
				if nameErr, ok := innerData["name"].(string); ok {
					return nil, printer.RenderError("unable to list functions", errors.New(nameErr))
				}
			}
		}
		printer.Logger.Debug("unable to list functions, request failed",
			slog.String("error", err.Error()),
			slog.Any("response", resData))
		return nil, printer.RenderError("unable to list functions", err)
	}
	allFunctions := initialResp.GetData()
	if initialResp.Meta.Results <= 20 {
		return allFunctions, nil
	}

	pageSize := int32(20)
	var currentPage int32
	for currentPage < initialResp.Meta.Pages {
		if maxResults > 0 && len(allFunctions) >= int(maxResults) {
			break
		}

		currentPage++
		previous := currentPage - 1
		next, httpRes, err := qc.FunctionsAPI.ProjectsFunctionsList(ctx, projectId).
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
						return nil, printer.RenderError("unable to list functions", errors.New(nameErr))
					}
				}
			}
			printer.Logger.Debug("unable to list functions, request failed",
				slog.String("error", err.Error()),
				slog.Any("response", resData))
			return nil, printer.RenderError("unable to list functions", err)
		}

		allFunctions = append(allFunctions, next.GetData()...)
	}

	if maxResults > 0 && len(allFunctions) > int(maxResults) {
		allFunctions = allFunctions[:maxResults]
	}

	return allFunctions, nil
}
