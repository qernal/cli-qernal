package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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

type Printer struct {
	resourceOut io.Writer
}

// NewPrinter creates a new Printer instance.
func NewPrinter() *Printer {
	return &Printer{
		resourceOut: os.Stdout, // Default output to stdout
	}
}

// SetResourceOutput sets the output for printing resources.
func (p *Printer) SetOut(out io.Writer) {
	p.resourceOut = out
}

// FormatOutput formats data based on the output type
func (p *Printer) FormatOutput(data interface{}, outputType string) string {
	var out io.Writer = os.Stdout
	if p.resourceOut != nil {
		out = p.resourceOut
	}

	switch outputType {
	case "json":
		prettyJSON, err := PrettyPrintJSON(data)
		if err != nil {
			fmt.Fprintln(out, "invalid json data")
			os.Exit(1)
		}
		return prettyJSON
	default:
		if mapData, ok := data.(map[string]interface{}); ok {
			formattedMap := new(bytes.Buffer)
			for key, value := range mapData {
				fmt.Fprintf(formattedMap, "%s: %v\n", key, value)
			}
			fmt.Fprint(out, formattedMap.String())
			return formattedMap.String()
		}
		formattedString := fmt.Sprintf("%v", data)
		return formattedString
	}
}

// PrintResource directly prints the given data to the output.
func (p *Printer) PrintResource(data string) {
	var out io.Writer = os.Stdout
	if p.resourceOut != nil {
		out = p.resourceOut
	}
	fmt.Fprintln(out, data)
}
