package domain

import "strings"

// RootName represents a workspace root name (dev, release, sandbox, etc.).
type RootName string

// WorkspaceType represents the classification of a workspace (e.g., sandbox, dev, release).
type WorkspaceType string

const (
	WorkspaceTypeSandbox  WorkspaceType = "sandbox"
	WorkspaceTypeDev      WorkspaceType = "dev"
	WorkspaceTypeRelease  WorkspaceType = "release"
	WorkspaceTypeUnknown  WorkspaceType = "unknown"
)

// ProjectType represents the classification of a project.
type ProjectType string

const (
	ProjectTypeSandbox     ProjectType = "sandbox"
	ProjectTypeSandboxGit  ProjectType = "sandbox+git"
	ProjectTypeDev         ProjectType = "dev"
	ProjectTypeRelease     ProjectType = "release"
	ProjectTypeDirty       ProjectType = "dirty"
	ProjectTypeExternal    ProjectType = "external"
	ProjectTypeDir         ProjectType = "dir"
)

// Root represents a workspace directory.
type Root struct {
	Name RootName
	Path string
	WorkspaceType WorkspaceType // Renamed from Zone
}

// Project represents a repository or directory in a workspace.
type Project struct {
	Name          string // Full unique name (e.g., github.com/user/repo)
	DisplayName   string // Short name for display (e.g., user/repo)
	Root          RootName
	Path          string
	WorkspaceType WorkspaceType // Renamed from Zone
	Type          ProjectType
	HasGit        bool
	Dirty         bool
	Branch        string // Lazy-loaded in TUI mode
}


// FormatDisplayName shortens a full project name for display.
// e.g., "github.com/user/repo" -> "user/repo"
func FormatDisplayName(name string) string {
	parts := strings.Split(name, "/")
	if len(parts) >= 3 {
		return strings.Join(parts[len(parts)-2:], "/")
	}
	return name
}

// DetermineWorkspaceType maps a root name to a WorkspaceType.
func DetermineWorkspaceType(rootName RootName) WorkspaceType { // Renamed from DetermineZone
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