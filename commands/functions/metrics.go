package functions

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/NimbleMarkets/ntcharts/linechart/timeserieslinechart"
	"github.com/charmbracelet/lipgloss"
	"github.com/qernal/cli-qernal/charm"
	"github.com/qernal/cli-qernal/commands/auth"
	"github.com/qernal/cli-qernal/pkg/client"
	"github.com/qernal/cli-qernal/pkg/common"
	"github.com/qernal/cli-qernal/pkg/helpers"
	"github.com/qernal/cli-qernal/pkg/utils"
	openapi_chaos_client "github.com/qernal/openapi-chaos-go-client"
	"github.com/spf13/cobra"
)

func NewMetricsCmd(printer *utils.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "metrics",
		Example: "qernal func metrics --project-id <project name> --function-id <function name>",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return helpers.ValidateProjectFlags(cmd)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()

			token, err := auth.GetQernalToken()
			if err != nil {
				return charm.RenderError("error creating qernal client", err)
			}

			qc, err := client.New(ctx, nil, nil, token)
			if err != nil {
				return charm.RenderError("error creating qernal client", err)
			}

			projectID, err := helpers.GetProjectID(cmd, &qc)
			if err != nil {
				return err
			}

			functionID, _ := cmd.Flags().GetString("function")
			// TODO: add watch flag for graph refreshes
			currentTime := time.Now().Format(time.RFC3339)
			pastTime := time.Now().Add(-15 * time.Minute).Format(time.RFC3339)

			// show http requests
			metricResp, httpRes, err := qc.MetricsAPI.MetricsAggregationsList(context.Background(), "httprequests").
				FProject(projectID).
				FFunction(functionID).
				FHistogramInterval(60).
				FTimestamps(openapi_chaos_client.LogsListFTimestampsParameter{
					After:  &pastTime,
					Before: &currentTime,
				}).
				Execute()
			if err != nil {
				resData, _ := client.ParseResponseData(httpRes)
				if data, ok := resData.(map[string]interface{}); ok {
					if innerData, ok := data["data"].(map[string]interface{}); ok {
						if nameErr, ok := innerData["name"].(string); ok {

							return printer.RenderError("unable to find function", errors.New(nameErr))
						}
					}
				}
				printer.Logger.Debug("Metrics collection failed ",
					slog.String("error", err.Error()),
					slog.Any("response", httpRes))
				return charm.RenderError("unable to find function", err)
			}

			if len(metricResp.MetricHttpAggregation.HttpCodes.Buckets) <= 0 {
				return errors.New(charm.RenderWarning("Function metrics are currently unavailable, make a few requests and try again"))
			}

			for _, r := range metricResp.MetricHttpAggregation.HttpCodes.Buckets {
				fmt.Println(*r.Key)
				HTTPGraph(*r.Histogram)
			}

			// show resource stats
			metricResp, httpRes, err = qc.MetricsAPI.MetricsAggregationsList(context.Background(), "resourcestats").
				FProject(projectID).
				FFunction(functionID).
				FHistogramInterval(60).
				FTimestamps(openapi_chaos_client.LogsListFTimestampsParameter{
					After:  &pastTime,
					Before: &currentTime,
				}).
				Execute()
			if err != nil {
				resData, _ := client.ParseResponseData(httpRes)
				if data, ok := resData.(map[string]interface{}); ok {
					if innerData, ok := data["data"].(map[string]interface{}); ok {
						if nameErr, ok := innerData["name"].(string); ok {
							printer.Logger.Debug("MEtrics collection failed ",
								slog.String("error", err.Error()),
								slog.Any("response", nameErr))
							return printer.RenderError("unable to find function", errors.New(nameErr))
						}
					}
				}
				printer.Logger.Debug("MEtrics collection failed ",
					slog.String("error", err.Error()),
					slog.Any("response", httpRes))
				return charm.RenderError("unable to find function", err)
			}

			// TODO: format header
			fmt.Println("Resource Stats")

			networkData := map[string]openapi_chaos_client.MetricResourceAggregationResourcesBucketsInner{}
			memoryData := map[string]openapi_chaos_client.MetricResourceAggregationResourcesBucketsInner{}

			res := metricResp.MetricResourceAggregation.Resources.Buckets
			for _, r := range res {
				if *r.Key == "cpu-usage" {
					CPUGraph(r)
				}

				if *r.Key == "network-tx" {
					networkData["tx"] = r
				}

				if *r.Key == "network-rx" {
					networkData["rx"] = r
				}

				if *r.Key == "memory-usage" {
					memoryData["usage"] = r
				}

				if *r.Key == "memory-available" {
					memoryData["capacity"] = r
				}
			}

			NetworkGraph(networkData["tx"], networkData["rx"])
			MemoryGraph(memoryData["usage"], memoryData["capacity"])

			// TODO: if watch and json provided, then error

			if err != nil {
				return charm.RenderError("unable to retrieve metrics, request failed with:", err)
			}

			// TODO: show both metric requests as json
			if common.OutputFormat == "json" {
				return nil
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&functionID, "function", "f", "", "function id")
	cmd.Flags().BoolVarP(&watch, "watch", "", false, "watch logs")

	_ = cmd.MarkFlagRequired("function")

	return cmd
}

func CPUGraph(res openapi_chaos_client.MetricResourceAggregationResourcesBucketsInner) {
	tslc := timeserieslinechart.New(41, 10)
	tslc.XLabelFormatter = timeserieslinechart.HourTimeLabelFormatter()

	for _, v := range res.Histogram.Buckets {
		date, err := time.Parse(time.RFC3339, *v.KeyAsString)
		if err != nil {
			fmt.Println("Error parsing date:", err)
			continue
		}

		tslc.Push(timeserieslinechart.TimePoint{
			Time:  date,
			Value: float64(*v.Gauge.Avg.Get() / 1000),
		})
	}

	tslc.DrawBraille()

	fmt.Println(tslc.View())
}

// network-tx
// network-rx
func NetworkGraph(tx openapi_chaos_client.MetricResourceAggregationResourcesBucketsInner, rx openapi_chaos_client.MetricResourceAggregationResourcesBucketsInner) {
	tslc := timeserieslinechart.New(41, 10)
	tslc.XLabelFormatter = timeserieslinechart.HourTimeLabelFormatter()

	// tx bucket
	for _, t := range tx.Histogram.Buckets {
		date, err := time.Parse(time.RFC3339, *t.KeyAsString)
		if err != nil {
			fmt.Println("Error parsing date:", err)
			continue
		}

		tslc.Push(timeserieslinechart.TimePoint{
			Time:  date,
			Value: float64(t.Counter.GetMax() / 1024 / 1024),
		})
	}

	tslc.SetStyle(
		lipgloss.NewStyle().
			Foreground(lipgloss.Color("10")), // green
	)

	// rx bucket
	for _, r := range rx.Histogram.Buckets {
		date, err := time.Parse(time.RFC3339, *r.KeyAsString)
		if err != nil {
			fmt.Println("Error parsing date:", err)
			continue
		}

		tslc.PushDataSet("rx", timeserieslinechart.TimePoint{
			Time:  date,
			Value: float64(r.Counter.GetMax() / 1024 / 1024),
		})
	}

	tslc.SetDataSetStyle("rx",
		lipgloss.NewStyle().
			Foreground(lipgloss.Color("21")), // green
	)

	// chart
	tslc.DrawBrailleAll()
	fmt.Println(tslc.View())
}

// memory-usage
// memory-available
func MemoryGraph(usage openapi_chaos_client.MetricResourceAggregationResourcesBucketsInner, capacity openapi_chaos_client.MetricResourceAggregationResourcesBucketsInner) {
	tslc := timeserieslinechart.New(41, 10)
	tslc.XLabelFormatter = timeserieslinechart.HourTimeLabelFormatter()

	// available bucket
	for _, t := range capacity.Histogram.Buckets {
		date, err := time.Parse(time.RFC3339, *t.KeyAsString)
		if err != nil {
			fmt.Println("Error parsing date:", err)
			continue
		}

		tslc.Push(timeserieslinechart.TimePoint{
			Time:  date,
			Value: float64(t.Counter.GetAvg() / 1024 / 1024),
		})
	}

	tslc.SetStyle(
		lipgloss.NewStyle().
			Foreground(lipgloss.Color("10")), // green
	)

	// usage bucket
	for _, r := range usage.Histogram.Buckets {
		date, err := time.Parse(time.RFC3339, *r.KeyAsString)
		if err != nil {
			fmt.Println("Error parsing date:", err)
			continue
		}

		tslc.PushDataSet("usage", timeserieslinechart.TimePoint{
			Time:  date,
			Value: float64(r.Counter.GetAvg() / 1024 / 1024),
		})
	}

	tslc.SetDataSetStyle("usage",
		lipgloss.NewStyle().
			Foreground(lipgloss.Color("21")), // green
	)

	// chart
	tslc.DrawBrailleAll()
	fmt.Println(tslc.View())
}

// http requests
func HTTPGraph(res openapi_chaos_client.MetricHttpAggregationHttpCodesBucketsInnerHistogram) {
	tslc := timeserieslinechart.New(41, 10)
	tslc.XLabelFormatter = timeserieslinechart.HourTimeLabelFormatter()

	for _, v := range res.Buckets {
		date, err := time.Parse(time.RFC3339, *v.KeyAsString)
		if err != nil {
			fmt.Println("Error parsing date:", err)
			continue
		}

		tslc.Push(timeserieslinechart.TimePoint{
			Time:  date,
			Value: float64(v.Gauge.GetMax()),
		})
	}

	tslc.DrawBraille()
	fmt.Println(tslc.View())
}
