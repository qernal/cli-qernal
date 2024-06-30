package charm

import (
	"errors"
	"fmt"
)

func RenderError(message string, err ...error) error {
	formattedMessage := ErrorStyle.Render(message)
	if len(err) > 0 && err[0] != nil {
		return fmt.Errorf("%s: %w", formattedMessage, err[0])
	}
	return errors.New(formattedMessage)
}

func RenderWarning(message string) string {
	formattedMessage := WarningStyle.Render(message)
	return formattedMessage
}
