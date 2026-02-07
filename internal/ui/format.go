package ui

import (
	"fmt"
	"os"

	"github.com/mi8bi/ghqx/internal/domain"
)

// FormatError formats a GhqxError for CLI output.
func FormatError(err error) string {
	if ghqxErr, ok := err.(*domain.GhqxError); ok {
		output := fmt.Sprintf("\nError: %s\n", ghqxErr.Message)

		if ghqxErr.Hint != "" {
			output += fmt.Sprintf("\nHint:\n  %s\n", ghqxErr.Hint)
		}

		return output
	}

	return fmt.Sprintf("\nError: %v\n", err)
}

// FormatDetailedError formats a GhqxError with internal details for debugging.
// Only shows internal details if GHQX_DEBUG environment variable is set.
func FormatDetailedError(err error) string {
	if ghqxErr, ok := err.(*domain.GhqxError); ok {
		output := FormatError(err)

		// Show internal details if debug mode is enabled
		if os.Getenv("GHQX_DEBUG") != "" {
			if ghqxErr.Internal != "" || ghqxErr.Cause != nil {
				output += "\nDebug Information:\n"
				if ghqxErr.Internal != "" {
					output += fmt.Sprintf("  Internal: %s\n", ghqxErr.Internal)
				}
				if ghqxErr.Cause != nil {
					output += fmt.Sprintf("  Cause: %v\n", ghqxErr.Cause)
				}
			}
		}

		return output
	}

	return FormatError(err)
}

// FormatSuccess formats a success message.
func FormatSuccess(message string) string {
	return fmt.Sprintf("✓ %s\n", message)
}

// FormatWarning formats a warning message.
func FormatWarning(message string) string {
	return fmt.Sprintf("⚠ %s\n", message)
}

// FormatInfo formats an info message.
func FormatInfo(message string) string {
	return fmt.Sprintf("• %s\n", message)
}
