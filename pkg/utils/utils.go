package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	math_rand "math/rand"
	"os"
	"strings"

	"github.com/qernal/cli-qernal/charm"
	"github.com/qernal/cli-qernal/pkg/common"
)

func PrettyPrintJSON(data interface{}) (string, error) {
	prettyJSON, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", err
	}
	return string(prettyJSON), nil
}

func FormatOutput(data interface{}, outputType string) string {
	// Handle error type
	if err, ok := data.(error); ok {
		errMap := map[string]string{"error": err.Error()}
		if outputType == "json" {
			jsonStr, _ := json.MarshalIndent(errMap, "", "  ")
			return string(jsonStr)
		}
		return charm.PlainTextStyle.Render(fmt.Sprintf("error: %v", err))
	}

	switch outputType {
	case "json":

		prettyJSON, err := PrettyPrintJSON(data)
		if err != nil {
			if err := charm.RenderError("invalid json data"); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to render error: %v\n", err)
			}
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

	// Handle error type
	if err, ok := data.(error); ok {
		if outputType == "json" {
			return fmt.Sprintf(`{"error": "%s"}`, err.Error())
		}
	}

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

// RenderError handles API error responses, formatting them as JSON when json output is enabled.
// For general error rendering with colored output, use `charm.RenderError` instead.
// Example:
//
//	if err != nil {
//	    return printer.RenderError("failed to create resource", err) // Will output {"error": "reason"} for json
//	}
func (p *Printer) RenderError(message string, err ...error) error {
	if len(err) > 0 && err[0] != nil {
		if common.OutputFormat == "json" {
			p.PrintResource(p.FormatOutput(err[0], common.OutputFormat))
			os.Exit(1)
			return nil // Empty error to avoid duplicate output
		}
		return fmt.Errorf("%s", charm.ErrorStyle.Render(message, err[0].Error()))
	}
	return errors.New(charm.ErrorStyle.Render(message))
}

// PrintResource directly prints the given data to the output.
func (p *Printer) PrintResource(data string) {
	var out io.Writer = os.Stdout
	if p.resourceOut != nil {
		out = p.resourceOut
	}
	fmt.Fprintln(out, data)
}

// generate random strings of a given length, for testing
func GenerateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[math_rand.Intn(len(charset))]
	}
	return string(b)
}
