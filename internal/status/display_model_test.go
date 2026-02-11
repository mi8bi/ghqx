package status

import (
	"testing"

	"github.com/mi8bi/ghqx/internal/domain"
	"github.com/mi8bi/ghqx/internal/i18n"
)

func TestFormatGitManagedAndStatus(t *testing.T) {
	if mg := formatGitManaged(true); mg != i18n.T("status.git.managed") {
		t.Fatalf("formatGitManaged(true) = %q", mg)
	}
	if mg := formatGitManaged(false); mg != i18n.T("status.git.unmanaged") {
		t.Fatalf("formatGitManaged(false) = %q", mg)
	}

	if s := formatStatus(false, false); s != "-" {
		t.Fatalf("formatStatus(non-git) should be '-' but got %q", s)
	}
	if s := formatStatus(true, true); s != i18n.T("status.repo.dirty") {
		t.Fatalf("formatStatus(git,dirty) = %q", s)
	}
	if s := formatStatus(true, false); s != i18n.T("status.repo.clean") {
		t.Fatalf("formatStatus(git,clean) = %q", s)
	}
}

func TestNewProjectDisplay(t *testing.T) {
	p := domain.Project{
		DisplayName:   "user/repo",
		WorkspaceType: domain.WorkspaceTypeDev,
		HasGit:        true,
		Dirty:         false,
		Path:          "/path/to/repo",
	}

	d := NewProjectDisplay(p)
	if d.Repo != p.DisplayName {
		t.Fatalf("Repo mismatch: %q", d.Repo)
	}
	if d.GitManaged != i18n.T("status.git.managed") {
		t.Fatalf("GitManaged mismatch: %q", d.GitManaged)
	}
}
