package selector

import (
	"testing"

	"github.com/mi8bi/ghqx/internal/status"
)

func makePD(repo, fullpath, workspace string) status.ProjectDisplay {
	return status.ProjectDisplay{Repo: repo, FullPath: fullpath, Workspace: workspace}
}

func TestSelectorFilteringAndCursor(t *testing.T) {
	projects := []status.ProjectDisplay{
		makePD("mi8bi/ghqx", "/p/mi8bi/ghqx", "sandbox"),
		makePD("other/repo", "/p/other/repo", "dev"),
	}

	m := NewModel(projects)

	// Test simple matching
	if !m.matchesSimpleQuery("mi8bi/ghqx", "ghqx") {
		t.Fatalf("expected simple query to match repo name")
	}

	// Test complex query
	if !m.matchesQuery(projects[0], "mi8bi/gh") {
		t.Fatalf("expected complex query to match")
	}

	// Directly set text input value and apply filter
	m.textInput.SetValue("ghqx")
	m.applyFilter()
	if len(m.filteredProjects) != 1 {
		t.Fatalf("expected 1 filtered project, got %d", len(m.filteredProjects))
	}

	// Cursor movement
	m.moveCursorDown()
	if m.cursor != 0 {
		t.Fatalf("cursor should wrap to 0 when only one item")
	}
	m.moveCursorUp()
	if m.cursor != 0 {
		t.Fatalf("cursor should remain 0 when only one item")
	}

	// View should render without panic
	_ = m.View()
}
