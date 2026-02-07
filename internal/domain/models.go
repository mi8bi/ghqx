package domain

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
	Name          string
	Root          RootName
	Path          string
	Zone          Zone
	Type          ProjectType
	HasGit        bool
	Dirty         bool
	Branch        string // Lazy-loaded in TUI mode
	WorktreeCount int    // Number of git worktrees (0 if not a git repo)
}

// Worktree represents a git worktree.
type Worktree struct {
	Path   string
	Branch string
	Bare   bool
	Locked bool
}

// PromoteRecord represents a single promote operation for undo history.
type PromoteRecord struct {
	Timestamp   int64
	ProjectName string
	FromRoot    RootName
	FromPath    string
	ToRoot      RootName
	ToPath      string
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
