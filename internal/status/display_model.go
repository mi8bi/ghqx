package status

import (
	"github.com/mi8bi/ghqx/internal/domain"
	"github.com/mi8bi/ghqx/internal/i18n"
)

// ProjectDisplay represents a project for display purposes.
// It contains formatted strings suitable for CLI/TUI output.
type ProjectDisplay struct {
	Repo       string         // Short repository name (e.g., user/repo)
	Workspace  string         // Workspace name (e.g., sandbox, dev, release)
	GitManaged string         // Git management status (e.g., "管理" / "Managed")
	Status     string         // Repository status (e.g., "clean" / "dirty")
	FullPath   string         // Full filesystem path to the project
	RawProject domain.Project // Original project data for detailed operations
}

// NewProjectDisplay creates a ProjectDisplay from a domain.Project.
// It applies i18n translations and formatting for display.
func NewProjectDisplay(p domain.Project) ProjectDisplay {
	return ProjectDisplay{
		Repo:       p.DisplayName,
		Workspace:  string(p.WorkspaceType),
		GitManaged: formatGitManaged(p.HasGit),
		Status:     formatStatus(p.HasGit, p.Dirty),
		FullPath:   p.Path,
		RawProject: p,
	}
}

// formatGitManaged returns a localized string for git management status.
func formatGitManaged(hasGit bool) string {
	if hasGit {
		return i18n.T("status.git.managed")
	}
	return i18n.T("status.git.unmanaged")
}

// formatStatus returns a localized string for repository status.
func formatStatus(hasGit, isDirty bool) string {
	if !hasGit {
		return "-" // Not applicable for non-git directories
	}

	if isDirty {
		return i18n.T("status.repo.dirty")
	}
	return i18n.T("status.repo.clean")
}
