package domain

import "strings"

// RootName represents a workspace root name (dev, release, sandbox, etc.).
type RootName string

// Zone represents the workspace zone classification.
type Zone string

const (
	ZoneSandbox Zone = "sandbox"
	ZoneDev     Zone = "dev"
	ZoneRelease Zone = "release"
	ZoneUnknown Zone = "unknown"
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
)

// Root represents a workspace directory.
type Root struct {
	Name RootName
	Path string
	Zone Zone
}

// Project represents a repository or directory in a workspace.
type Project struct {
	Name          string // Full unique name (e.g., github.com/user/repo)
	DisplayName   string // Short name for display (e.g., user/repo)
	Root          RootName
	Path          string
	Zone          Zone
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

// DetermineZone maps a root name to a zone.
func DetermineZone(rootName RootName) Zone {
	switch rootName {
	case "sandbox":
		return ZoneSandbox
	case "dev":
		return ZoneDev
	case "release":
		return ZoneRelease
	default:
		return ZoneUnknown
	}
}
