package utils

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/qernal/cli-qernal/charm"
)

func PrettyPrintJSON(data interface{}) (string, error) {
	prettyJSON, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", err
	}
	return string(prettyJSON), nil
}

func FormatOutput(data interface{}, outputType string) string {
	switch outputType {
	case "json":
		prettyJSON, err := PrettyPrintJSON(data)
		if err != nil {
			charm.RenderError("invalid json data")
			os.Exit(1)
		}
		return prettyJSON
	default:
		if mapData, ok := data.(map[string]interface{}); ok {
			formattedMap := ""
			for key, value := range mapData {
				formattedMap += fmt.Sprintf("%s: %v\n", key, value)
			}
			return charm.PlainTextStyle.Render(formattedMap)
		}
		return charm.PlainTextStyle.Render(fmt.Sprintf("%v", data))
	}
}
