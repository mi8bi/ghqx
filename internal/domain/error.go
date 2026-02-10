package domain

import "fmt"

// ErrorCode represents a machine-readable error code.
type ErrorCode string

const (
	ErrCodeUnknown          ErrorCode = "UNKNOWN"
	ErrCodeConfigNotFound   ErrorCode = "CONFIG_NOT_FOUND"
	ErrCodeConfigInvalid    ErrorCode = "CONFIG_INVALID"
	ErrCodeRootNotFound     ErrorCode = "ROOT_NOT_FOUND"
	ErrCodeProjectNotFound  ErrorCode = "PROJECT_NOT_FOUND"
	ErrCodeDirtyRepo        ErrorCode = "DIRTY_REPO"
	ErrCodeInvalidPath      ErrorCode = "INVALID_PATH"
	ErrCodeGitError         ErrorCode = "GIT_ERROR"
	ErrCodeFSError          ErrorCode = "FS_ERROR"
)

// GhqxError represents a domain error with user-friendly output.
// This separates user-facing messages from internal error details.
type GhqxError struct {
	Code    ErrorCode
	Message string // User-facing message
	Hint    string // Actionable advice for user
	Cause   error  // Internal error (not shown to user by default)
	Internal string // Internal debug information
}

// Error implements the error interface.
// Returns user-facing message for normal output.
func (e *GhqxError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Unwrap returns the underlying cause.
func (e *GhqxError) Unwrap() error {
	return e.Cause
}

// DetailedError returns the full error with internal details.
// Use this for logging or debugging.
func (e *GhqxError) DetailedError() string {
	msg := e.Error()
	if e.Internal != "" {
		msg += fmt.Sprintf(" [internal: %s]", e.Internal)
	}
	if e.Cause != nil {
		msg += fmt.Sprintf(" [cause: %v]", e.Cause)
	}
	return msg
}

// IsUserError returns true if this is a user-correctable error.
func (e *GhqxError) IsUserError() bool {
	return e.Hint != ""
}

// NewError creates a new GhqxError.
func NewError(code ErrorCode, message string) *GhqxError {
	return &GhqxError{
		Code:    code,
		Message: message,
	}
}

// NewErrorWithCause creates a new GhqxError with an underlying cause.
func NewErrorWithCause(code ErrorCode, message string, cause error) *GhqxError {
	return &GhqxError{
		Code:    code,
		Message: message,
		Cause:   cause,
	}
}

// WithHint adds a hint to the error.
func (e *GhqxError) WithHint(hint string) *GhqxError {
	e.Hint = hint
	return e
}

// WithInternal adds internal debug information.
func (e *GhqxError) WithInternal(internal string) *GhqxError {
	e.Internal = internal
	return e
}
