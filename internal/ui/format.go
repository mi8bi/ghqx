package ui

import (
	"fmt"
	"os"

	"github.com/mi8bi/ghqx/internal/domain"
	"github.com/mi8bi/ghqx/internal/i18n"
)

// FormatError formats a GhqxError for CLI output.
func FormatError(err error) string {
	if ghqxErr, ok := err.(*domain.GhqxError); ok {
		output := fmt.Sprintf("\n%s: %s\n", i18n.T("ui.error.prefix"), ghqxErr.Message)

		if ghqxErr.Hint != "" {
			output += fmt.Sprintf("\n%s:\n  %s\n", i18n.T("ui.error.hintPrefix"), ghqxErr.Hint)
		}

		return output
	}

	return fmt.Sprintf("\n%s: %v\n", i18n.T("ui.error.prefix"), err)
}

// FormatDetailedError formats a GhqxError with internal details for debugging.
// Only shows internal details if GHQX_DEBUG environment variable is set.
func FormatDetailedError(err error) string {
	if ghqxErr, ok := err.(*domain.GhqxError); ok {
		output := FormatError(err)

		// Show internal details if debug mode is enabled
		if os.Getenv("GHQX_DEBUG") != "" {
			if ghqxErr.Internal != "" || ghqxErr.Cause != nil {
				output += fmt.Sprintf("\n%s:\n", i18n.T("ui.error.debugInfoPrefix"))
				if ghqxErr.Internal != "" {
					output += fmt.Sprintf("  %s: %s\n", i18n.T("ui.error.internalPrefix"), ghqxErr.Internal)
				}
				if ghqxErr.Cause != nil {
					output += fmt.Sprintf("  %s: %v\n", i18n.T("ui.error.causePrefix"), ghqxErr.Cause)
				}
			}
		}

		return output
	}

	return FormatError(err)
}

// FormatSuccess formats a success message.
func FormatSuccess(message string) string {
	return fmt.Sprintf("%s %s\n", i18n.T("ui.success.prefix"), message)
}

// FormatWarning formats a warning message.
func FormatWarning(message string) string {
	return fmt.Sprintf("%s %s\n", i18n.T("ui.warning.prefix"), message)
}

// FormatInfo formats an info message.
func FormatInfo(message string) string {
	return fmt.Sprintf("%s %s\n", i18n.T("ui.info.prefix"), message)
}
