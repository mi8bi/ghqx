package domain

import "strings"

// RootName represents a workspace root identifier (dev, release, sandbox, etc.).
type RootName string

// WorkspaceType categorizes the workspace based on its purpose and usage patterns.
type WorkspaceType string

const (
	// WorkspaceTypeSandbox represents an experimental/temporary workspace
	WorkspaceTypeSandbox WorkspaceType = "sandbox"
	// WorkspaceTypeDev represents an active development workspace
	WorkspaceTypeDev WorkspaceType = "dev"
	// WorkspaceTypeRelease represents a production/stable release workspace
	WorkspaceTypeRelease WorkspaceType = "release"
	// WorkspaceTypeUnknown represents an unrecognized workspace type
	WorkspaceTypeUnknown WorkspaceType = "unknown"
)

// ProjectType categorizes the project based on its location and git status.
type ProjectType string

const (
	// ProjectTypeSandbox represents a non-git project in the sandbox workspace
	ProjectTypeSandbox ProjectType = "sandbox"
	// ProjectTypeSandboxGit represents a git repository in the sandbox workspace
	ProjectTypeSandboxGit ProjectType = "sandbox+git"
	// ProjectTypeDev represents a git repository in the dev workspace
	ProjectTypeDev ProjectType = "dev"
	// ProjectTypeRelease represents a git repository in the release workspace
	ProjectTypeRelease ProjectType = "release"
	// ProjectTypeDirty represents a project with uncommitted changes
	ProjectTypeDirty ProjectType = "dirty"
	// ProjectTypeExternal represents a project outside managed workspaces
	ProjectTypeExternal ProjectType = "external"
	// ProjectTypeDir represents a regular directory (not a git repository)
	ProjectTypeDir ProjectType = "dir"
)

// Root represents a workspace directory configuration.
type Root struct {
	// Name is the identifier for this root (e.g., "dev", "sandbox")
	Name RootName
	// Path is the absolute filesystem path to this workspace root
	Path string
	// WorkspaceType indicates the purpose of this workspace
	WorkspaceType WorkspaceType
}

// Project represents a repository or directory within a workspace.
type Project struct {
	// Name is the fully qualified name (e.g., "github.com/user/repo")
	Name string
	// DisplayName is a shorter version for UI display (e.g., "user/repo")
	DisplayName string
	// Root identifies which workspace root this project belongs to
	Root RootName
	// Path is the absolute filesystem path to this project
	Path string
	// WorkspaceType indicates which type of workspace this project is in
	WorkspaceType WorkspaceType
	// Type categorizes the project based on location and git status
	Type ProjectType
	// HasGit indicates whether this is a git repository
	HasGit bool
	// Dirty indicates whether the repository has uncommitted changes
	Dirty bool
	// Branch is the current git branch name (lazy-loaded in TUI mode)
	Branch string
}

// FormatDisplayName shortens a fully qualified project name for display purposes.
// Removes domain/org prefix and keeps only the last two path components.
// Example: "github.com/user/repo" -> "user/repo"
func FormatDisplayName(name string) string {
	parts := strings.Split(name, "/")
	if len(parts) >= 3 {
		// Return last two components (organization/repo)
		return strings.Join(parts[len(parts)-2:], "/")
	}
	// Fallback for short names
	return name
}

// DetermineWorkspaceType maps a root name to its corresponding WorkspaceType.
// Returns WorkspaceTypeUnknown for unrecognized root names.
func DetermineWorkspaceType(rootName RootName) WorkspaceType {
	switch rootName {
	case "sandbox":
		return WorkspaceTypeSandbox
	case "dev":
		return WorkspaceTypeDev
	case "release":
		return WorkspaceTypeRelease
	default:
		return WorkspaceTypeUnknown
	}
}