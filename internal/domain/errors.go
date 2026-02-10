package domain

import (
	"fmt"

	"github.com/mi8bi/ghqx/internal/i18n"
)
// Common errors used across the application.
// These provide consistent error messages and hints for users.

// Config errors
var (
	ErrConfigNotFoundAny = NewError(
		ErrCodeConfigNotFound,
		i18n.T("error.config.notFoundAny.message"),
	).WithHint(i18n.T("error.config.notFoundAny.hint"))

	ErrConfigNotFoundAt = func(path string) *GhqxError {
		return NewError(
			ErrCodeConfigNotFound,
			i18n.T("error.config.notFoundAt.message"),
		).WithHint(i18n.T("error.config.notFoundAt.hint")).
			WithInternal("path: " + path)
	}

	ErrConfigInvalidTOML = func(cause error) *GhqxError {
		return NewErrorWithCause(
			ErrCodeConfigInvalid,
			i18n.T("error.config.invalidToml.message"),
			cause,
		).WithHint(i18n.T("error.config.invalidToml.hint"))
	}

	ErrConfigNoRoots = NewError(
		ErrCodeConfigInvalid,
		i18n.T("error.config.noRoots.message"),
	).WithHint(i18n.T("error.config.noRoots.hint"))

	ErrConfigInvalidDefaultRoot = NewError(
		ErrCodeConfigInvalid,
		i18n.T("error.config.invalidDefaultRoot.message"),
	).WithHint(i18n.T("error.config.invalidDefaultRoot.hint"))
)

// Root errors
var (
	ErrRootNotFound = func(name string) *GhqxError {
		return NewError(
			ErrCodeRootNotFound,
			fmt.Sprintf(i18n.T("error.root.notFound.message"), name),
		).WithHint(i18n.T("error.root.notFound.hint"))
	}

	ErrRootDirNotExist = func(name, path string) *GhqxError {
		return NewError(
			ErrCodeRootNotFound,
			fmt.Sprintf(i18n.T("error.root.dirNotExist.message"), name),
		).WithHint(i18n.T("error.root.dirNotExist.hint")).
			WithInternal("path: " + path)
	}
)

// Project errors
var (
	ErrProjectNotFound = func(name string) *GhqxError {
		return NewError(
			ErrCodeProjectNotFound,
			fmt.Sprintf(i18n.T("error.project.notFound.message"), name),
		).WithHint(i18n.T("error.project.notFound.hint"))
	}

	ErrProjectNameInvalid = NewError(
		ErrCodeInvalidPath,
		i18n.T("error.project.nameInvalid.message"),
	).WithHint(i18n.T("error.project.nameInvalid.hint"))

	ErrArgumentRequired = NewError(
		ErrCodeInvalidPath,
		i18n.T("error.argument.required"),
	)
)

// Git errors
var (
	ErrGitDirtyRepo = NewError(
		ErrCodeDirtyRepo,
		i18n.T("error.git.dirtyRepo.message"),
	).WithHint(i18n.T("error.git.dirtyRepo.hint"))

	ErrGitTimeout = func(operation string) *GhqxError {
		return NewError(
			ErrCodeGitError,
			fmt.Sprintf(i18n.T("error.git.timeout.message"), operation),
		).WithInternal("timeout exceeded")
	}

	ErrGitCommandFailed = func(operation string, cause error) *GhqxError {
		return NewErrorWithCause(
			ErrCodeGitError,
			fmt.Sprintf(i18n.T("error.git.commandFailed.message"), operation),
			cause,
		)
	}
)



// Filesystem errors
var (
	ErrFSReadDir = func(cause error) *GhqxError {
		return NewErrorWithCause(
			ErrCodeFSError,
			i18n.T("error.fs.readDir.message"),
			cause,
		)
	}

	ErrFSCreateDir = func(cause error) *GhqxError {
		return NewErrorWithCause(
			ErrCodeFSError,
			i18n.T("error.fs.createDir.message"),
			cause,
		)
	}

	ErrFSScanRoot = func(cause error) *GhqxError {
		return NewErrorWithCause(
			ErrCodeFSError,
			i18n.T("error.fs.scanRoot.message"),
			cause,
		)
	}
)

