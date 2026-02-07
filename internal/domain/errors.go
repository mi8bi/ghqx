package domain

// Common errors used across the application.
// These provide consistent error messages and hints for users.

// Config errors
var (
	ErrConfigNotFoundAny = NewError(
		ErrCodeConfigNotFound,
		"No configuration file found",
	).WithHint("Run 'ghqx config init' to create a config file")

	ErrConfigNotFoundAt = func(path string) *GhqxError {
		return NewError(
			ErrCodeConfigNotFound,
			"Config file not found at specified path",
		).WithHint("Check the path provided with --config flag").
			WithInternal("path: " + path)
	}

	ErrConfigInvalidTOML = func(cause error) *GhqxError {
		return NewErrorWithCause(
			ErrCodeConfigInvalid,
			"Failed to parse config file",
			cause,
		).WithHint("Check the TOML syntax in your config file")
	}

	ErrConfigNoRoots = NewError(
		ErrCodeConfigInvalid,
		"No roots defined in configuration",
	).WithHint("Add at least one root in the [roots] section")

	ErrConfigInvalidDefaultRoot = NewError(
		ErrCodeConfigInvalid,
		"Default root does not exist in roots",
	).WithHint("Set default.root to one of the defined roots")
)

// Root errors
var (
	ErrRootNotFound = func(name string) *GhqxError {
		return NewError(
			ErrCodeRootNotFound,
			"Root not found: "+name,
		).WithHint("Check your config.toml for available roots")
	}

	ErrRootDirNotExist = func(name, path string) *GhqxError {
		return NewError(
			ErrCodeRootNotFound,
			"Root directory does not exist: "+name,
		).WithHint("Create the directory or update config.toml").
			WithInternal("path: " + path)
	}
)

// Project errors
var (
	ErrProjectNotFound = func(name string) *GhqxError {
		return NewError(
			ErrCodeProjectNotFound,
			"Project not found: "+name,
		).WithHint("Use 'ghqx status' to see all available projects")
	}

	ErrProjectNameInvalid = NewError(
		ErrCodeInvalidPath,
		"Invalid project name",
	).WithHint("Project name contains forbidden characters")
)

// Git errors
var (
	ErrGitDirtyRepo = NewError(
		ErrCodeDirtyRepo,
		"Repository has uncommitted changes",
	).WithHint("Commit or stash changes, or use --force")

	ErrGitTimeout = func(operation string) *GhqxError {
		return NewError(
			ErrCodeGitError,
			"Git operation timed out: "+operation,
		).WithInternal("timeout exceeded")
	}

	ErrGitCommandFailed = func(operation string, cause error) *GhqxError {
		return NewErrorWithCause(
			ErrCodeGitError,
			"Git operation failed: "+operation,
			cause,
		)
	}

	ErrGitWorktreeList = func(cause error) *GhqxError {
		return NewErrorWithCause(
			ErrCodeWorktreeError,
			"Failed to list git worktrees",
			cause,
		).WithHint("Ensure git worktree is properly configured")
	}
)

// Promote errors
var (
	ErrPromoteDestExists = func(path string) *GhqxError {
		return NewError(
			ErrCodePromoteFailed,
			"Destination already exists",
		).WithHint("Choose a different name or remove the existing directory").
			WithInternal("path: " + path)
	}

	ErrPromoteSourceNotFound = func(name string) *GhqxError {
		return NewError(
			ErrCodeProjectNotFound,
			"Project not found in source root: "+name,
		).WithHint("Check project name and source root")
	}

	ErrPromoteMoveFailed = func(cause error) *GhqxError {
		return NewErrorWithCause(
			ErrCodePromoteFailed,
			"Failed to move project",
			cause,
		)
	}
)

// Undo errors
var (
	ErrUndoDisabled = NewError(
		ErrCodeUndoNotAvailable,
		"History is disabled",
	).WithHint("Enable history in config.toml")

	ErrUndoNoHistory = NewError(
		ErrCodeUndoNotAvailable,
		"No promote history found",
	).WithHint("Nothing to undo")

	ErrUndoDestMissing = NewError(
		ErrCodeUndoNotAvailable,
		"Promoted project no longer exists at destination",
	).WithHint("Cannot undo: project may have been moved or deleted")

	ErrUndoSourceOccupied = NewError(
		ErrCodeUndoNotAvailable,
		"Source location is occupied",
	).WithHint("Cannot undo: original location is no longer empty")
)

// Filesystem errors
var (
	ErrFSReadDir = func(cause error) *GhqxError {
		return NewErrorWithCause(
			ErrCodeFSError,
			"Failed to read directory",
			cause,
		)
	}

	ErrFSCreateDir = func(cause error) *GhqxError {
		return NewErrorWithCause(
			ErrCodeFSError,
			"Failed to create directory",
			cause,
		)
	}

	ErrFSScanRoot = func(cause error) *GhqxError {
		return NewErrorWithCause(
			ErrCodeFSError,
			"Failed to scan root directory",
			cause,
		)
	}
)
