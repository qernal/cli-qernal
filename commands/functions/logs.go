package functions

import (
	"context"
	"fmt"
	"time"

	"github.com/qernal/cli-qernal/charm"
	"github.com/qernal/cli-qernal/commands/auth"
	"github.com/qernal/cli-qernal/pkg/client"
	"github.com/qernal/cli-qernal/pkg/common"
	"github.com/qernal/cli-qernal/pkg/helpers"
	"github.com/qernal/cli-qernal/pkg/utils"
	openapi_chaos_client "github.com/qernal/openapi-chaos-go-client"
	"github.com/spf13/cobra"
)

func NewLogsCmd(printer *utils.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "logs",
		Example: "qernal func logs --project-id <project name> --function-id <function name>",
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
			watch := cmd.Flags().Changed("watch")

			// TODO: if watch and json provided, then error

			// if we're watching logs, loop forever
			if watch {
				lastHashes := [][]byte{}
				lastWatchDate := ""
				writeLog := true

				for {
					logs, err := getLogs(projectID, functionID, &qc)

					if err != nil {
						return charm.RenderError("unable to list logs,  request failed with:", err)
					}

					currentHashes := allHashes(logs.Data)

					for _, log := range logs.Data {
						currentHash := hashLog(log)
						writeLog = true

						if lastWatchDate != "" {
							logTime, err := time.Parse(time.RFC3339, *log.Log.Timestamp)
							if err != nil {
								return charm.RenderError("error parsing log timestamp", err)
							}
							lastWatchTime, err := time.Parse(time.RFC3339, lastWatchDate)
							if err != nil {
								return charm.RenderError("error parsing last watch date", err)
							}

							if logTime.Before(lastWatchTime) {
								writeLog = false
							}
						}

						// check if we've seen this log in the last batch
						// we can skip this check if we're skipping on time
						if writeLog {
							for _, hash := range lastHashes {

								if string(currentHash) == string(hash) {
									writeLog = false
									break
								}
							}
						}

						if writeLog {
							printer.PrintResource(fmt.Sprintf("%s: %s", *log.Log.Timestamp, *log.Log.Line))
							lastWatchDate = *log.Log.Timestamp
						}
					}

					// overwrite last hashes
					lastHashes = currentHashes

					// TODO: allow this to be configurable via flag
					time.Sleep(5 * time.Second)
				}
			}

			logs, err := getLogs(projectID, functionID, &qc)

			if err != nil {
				return charm.RenderError("unable to list logs,  request failed with:", err)
			}

			// show logs as json
			if common.OutputFormat == "json" {
				printer.PrintResource(utils.FormatOutput(logs, common.OutputFormat))
				return nil
			}

			// show logs (non-watch)
			printer.PrintResource(formatLogs(logs.Data))
			return nil
		},
	}

	cmd.Flags().StringVarP(&functionID, "function", "f", "", "function id")
	cmd.Flags().BoolVarP(&watch, "watch", "", false, "watch logs")

	_ = cmd.MarkFlagRequired("function")

	return cmd
}

// reverse order of logs
func reorderLogs(logResp *openapi_chaos_client.ListLogResponse) *openapi_chaos_client.ListLogResponse {
	for i, j := 0, len(logResp.Data)-1; i < j; i, j = i+1, j-1 {
		logResp.Data[i], logResp.Data[j] = logResp.Data[j], logResp.Data[i]
	}

	return logResp
}

// get logs from qernal
func getLogs(projectID string, functionID string, qc *client.QernalAPIClient) (openapi_chaos_client.ListLogResponse, error) {
	logResp, _, err := qc.LogsAPI.LogsList(context.Background()).FProject(projectID).FFunction(functionID).Execute()

	if err != nil {
		return openapi_chaos_client.ListLogResponse{}, err
	}

	logResp = reorderLogs(logResp)
	return *logResp, nil
}

// format logs for printing
func formatLogs(logs []openapi_chaos_client.Log) string {
	formattedLogs := ""
	for _, log := range logs {
		formattedLogs += fmt.Sprintf("%s: %s", *log.Log.Timestamp, *log.Log.Line) + "\n"
	}

	return formattedLogs
}

// hash a log item
func hashLog(log openapi_chaos_client.Log) []byte {
	return sdbmHash(fmt.Sprintf("%s%s%s%s%s%s%s", *log.Container, *log.Function, *log.Project, *log.Log.Type, *log.Log.Stream, *log.Log.Line, *log.Log.Timestamp))
}

// sdbm hash function
func sdbmHash(data string) []byte {
	var hash uint64
	for i := 0; i < len(data); i++ {
		hash = uint64(data[i]) + (hash << 6) + (hash << 16) - hash
	}
	return []byte(fmt.Sprintf("%x", hash))
}

// hash all log entries to build hash list
func allHashes(logs []openapi_chaos_client.Log) [][]byte {
	hashes := [][]byte{}
	for _, log := range logs {
		hashes = append(hashes, hashLog(log))
	}

	return hashes
}
