package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"

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

type Printer struct {
	resourceOut io.Writer
	//mostly for debug level logs, for rendering errors, see charm package
	Logger *slog.Logger
}

// NewPrinter creates a new Printer instance
func NewPrinter() *Printer {
	// Initialize a LevelVar with the default level set to Info
	lvl := new(slog.LevelVar)
	lvl.Set(slog.LevelInfo)

	// Check if the LOG_LEVEL environment variable is set to "debug"
	envLogLevel := os.Getenv("LOG_LEVEL")
	if strings.ToLower(envLogLevel) == "debug" {
		lvl.Set(slog.LevelDebug)
	}

	// Initialize the logger with the specified log level
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: lvl,
	}))

	return &Printer{
		resourceOut: os.Stdout,
		Logger:      logger,
	}
}

// SetResourceOutput sets the output for printing resources.
func (p *Printer) SetOut(out io.Writer) {
	p.resourceOut = out
}

// FormatOutput formats data based on the output type
func (p *Printer) FormatOutput(data interface{}, outputType string) string {
	switch outputType {
	case "json":
		prettyJSON, err := PrettyPrintJSON(data)
		if err != nil {
			return "invalid json data"
		}
		return prettyJSON
	default:
		if mapData, ok := data.(map[string]interface{}); ok {
			formattedMap := new(bytes.Buffer)
			for key, value := range mapData {
				fmt.Fprintf(formattedMap, "%s: %v\n", key, value)
			}
			return formattedMap.String()
		}
		return fmt.Sprintf("%v", data)
	}
}

// PrintResource directly prints the given data to the output.
func (p *Printer) PrintResource(data string) {
	var out io.Writer = os.Stdout
	if p.resourceOut != nil {
		out = p.resourceOut
	}

	// Since we're keeping the string parameter, we'll check if it's empty
	if len(strings.TrimSpace(data)) == 0 {
		fmt.Fprintln(out, "no data found, cannot print nil value")
		return
	}

	fmt.Fprintln(out, data)
}

func (p *Printer) PrintResourceR(data *string) {
	var out io.Writer = os.Stdout
	if p.resourceOut != nil {
		out = p.resourceOut
	}

	if data == nil {
		fmt.Fprintln(out, "no data found, cannot print nil value")
		return
	}

	fmt.Fprintln(out, *data)
}
