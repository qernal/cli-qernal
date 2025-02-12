package helpers

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/qernal/cli-qernal/pkg/client"
	"github.com/qernal/cli-qernal/pkg/utils"
	openapi_chaos_client "github.com/qernal/openapi-chaos-go-client"
	"github.com/spf13/viper"
)

// ParseFunctionConfig returns a function body from the supplied config
func ParseFunctionConfig(file string) (*openapi_chaos_client.FunctionBody, error) {

	viper.SetConfigFile(file)
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return &openapi_chaos_client.FunctionBody{}, fmt.Errorf("error reading function file: %w", err)
	}

	function := openapi_chaos_client.FunctionBody{
		Name:        viper.GetString("name"),
		Description: viper.GetString("description"),
		Version:     viper.GetString("version"),
		Type:        openapi_chaos_client.FunctionType(viper.GetString("type")),
		Image:       viper.GetString("image"),
		Port:        viper.GetInt32("port"),
		ProjectId:   viper.GetString("project_id"),

		Size: openapi_chaos_client.FunctionSize{
			Cpu:    viper.GetInt32("size.cpu"),
			Memory: viper.GetInt32("size.memory"),
		},

		Scaling: openapi_chaos_client.FunctionScaling{
			Type: viper.GetString("scaling.type"),
			Low:  viper.GetInt32("scaling.low"),
			High: viper.GetInt32("scaling.high"),
		},
	}

	if viper.IsSet("deployments") {
		var deployments []openapi_chaos_client.FunctionDeploymentBody
		for _, d := range viper.Get("deployments").([]interface{}) {
			deployment := d.(map[string]interface{})
			location := deployment["location"].(map[string]interface{})
			replicas := deployment["replicas"].(map[string]interface{})

			continent := location["continent"].(string)
			city := location["city"].(string)
			country := location["country"].(string)

			deployments = append(deployments, openapi_chaos_client.FunctionDeploymentBody{
				Location: openapi_chaos_client.Location{
					Continent:  &continent,
					Country:    &country,
					City:       &city,
					ProviderId: location["provider_id"].(string),
				},
				Replicas: openapi_chaos_client.FunctionReplicas{
					Min: getInt32FromInterface(replicas["min"]),
					Max: getInt32FromInterface(replicas["max"]),
					Affinity: openapi_chaos_client.FunctionReplicasAffinity{
						Cloud:   replicas["affinity"].(map[string]interface{})["cloud"].(bool),
						Cluster: replicas["affinity"].(map[string]interface{})["cluster"].(bool),
					},
				},
			})
		}
		function.Deployments = deployments
	}

	if viper.IsSet("routes") {
		var routes []openapi_chaos_client.FunctionRoute
		for _, r := range viper.Get("routes").([]interface{}) {
			route := r.(map[string]interface{})
			methods := make([]string, len(route["methods"].([]interface{})))
			for i, m := range route["methods"].([]interface{}) {
				methods[i] = m.(string)
			}

			routes = append(routes, openapi_chaos_client.FunctionRoute{
				Path:    route["path"].(string),
				Methods: methods,
				Weight:  getInt32FromInterface(route["weight"]),
			})
		}
		function.Routes = routes
	}

	if viper.IsSet("secrets") {
		var secrets []openapi_chaos_client.FunctionEnv
		for _, s := range viper.Get("secrets").([]interface{}) {
			secret := s.(map[string]interface{})
			secrets = append(secrets, openapi_chaos_client.FunctionEnv{
				Name:      secret["name"].(string),
				Reference: secret["reference"].(string),
			})
		}
		function.Secrets = secrets
	}

	if viper.IsSet("compliance") {
		compliance := viper.GetStringSlice("compliance")
		function.Compliance = make([]openapi_chaos_client.FunctionCompliance, len(compliance))
		for i, c := range compliance {
			function.Compliance[i] = openapi_chaos_client.FunctionCompliance(c)
		}
	}
	return &function, nil
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
			return nil, printer.RenderError("unable to list projects", err)
		}

		allFunctions = append(allFunctions, next.GetData()...)
	}

	if maxResults > 0 && len(allFunctions) > int(maxResults) {
		allFunctions = allFunctions[:maxResults]
	}

	return allFunctions, nil
}

func getInt32FromInterface(v interface{}) int32 {
	switch v := v.(type) {
	case int:
		return int32(v)
	case float64:
		return int32(v)
	default:
		return 0
	}
}
