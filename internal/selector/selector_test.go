package selector

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
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

// Additional tests for better coverage

func TestNewModel(t *testing.T) {
	projects := []status.ProjectDisplay{
		makePD("user/repo", "/path/repo", "dev"),
	}

	m := NewModel(projects)

	if len(m.projects) != 1 {
		t.Error("projects should be set")
	}

	if len(m.filteredProjects) != 1 {
		t.Error("filteredProjects should initially equal projects")
	}

	if m.cursor != 0 {
		t.Error("cursor should start at 0")
	}

	if m.selected != "" {
		t.Error("selected should be empty initially")
	}

	if m.quitting {
		t.Error("quitting should be false initially")
	}
}

func TestModelInit(t *testing.T) {
	projects := []status.ProjectDisplay{
		makePD("user/repo", "/path/repo", "dev"),
	}

	m := NewModel(projects)
	cmd := m.Init()

	if cmd == nil {
		t.Error("Init should return a command")
	}
}

func TestModelUpdateWithKeyMessages(t *testing.T) {
	projects := []status.ProjectDisplay{
		makePD("repo1", "/path/repo1", "dev"),
		makePD("repo2", "/path/repo2", "dev"),
		makePD("repo3", "/path/repo3", "dev"),
	}

	m := NewModel(projects)

	// Test down navigation
	keyMsg := tea.KeyMsg{Type: tea.KeyDown}
	newModel, _ := m.Update(keyMsg)
	m = newModel.(Model)
	if m.cursor != 1 {
		t.Errorf("cursor should be 1 after down, got %d", m.cursor)
	}

	// Test up navigation
	keyMsg = tea.KeyMsg{Type: tea.KeyUp}
	newModel, _ = m.Update(keyMsg)
	m = newModel.(Model)
	if m.cursor != 0 {
		t.Errorf("cursor should be 0 after up, got %d", m.cursor)
	}

	// Test cursor wrapping at bottom
	m.cursor = 2
	keyMsg = tea.KeyMsg{Type: tea.KeyDown}
	newModel, _ = m.Update(keyMsg)
	m = newModel.(Model)
	if m.cursor != 0 {
		t.Error("cursor should wrap to 0 at bottom")
	}

	// Test cursor wrapping at top
	m.cursor = 0
	keyMsg = tea.KeyMsg{Type: tea.KeyUp}
	newModel, _ = m.Update(keyMsg)
	m = newModel.(Model)
	if m.cursor != 2 {
		t.Error("cursor should wrap to last at top")
	}
}

func TestModelUpdateWithEnter(t *testing.T) {
	projects := []status.ProjectDisplay{
		makePD("repo1", "/path/repo1", "dev"),
		makePD("repo2", "/path/repo2", "dev"),
	}

	m := NewModel(projects)
	m.cursor = 1

	keyMsg := tea.KeyMsg{Type: tea.KeyEnter}
	newModel, cmd := m.Update(keyMsg)
	m = newModel.(Model)

	if m.selected == "" {
		t.Error("selected should be set after Enter")
	}

	if m.selected != "/path/repo2" {
		t.Errorf("expected selected /path/repo2, got %s", m.selected)
	}

	if cmd == nil {
		t.Error("Enter should return quit command")
	}
}

func TestModelUpdateWithQuit(t *testing.T) {
	projects := []status.ProjectDisplay{
		makePD("repo1", "/path/repo1", "dev"),
	}

	m := NewModel(projects)

	// Test Esc
	keyMsg := tea.KeyMsg{Type: tea.KeyEsc}
	newModel, cmd := m.Update(keyMsg)
	m = newModel.(Model)

	if !m.quitting {
		t.Error("quitting should be true after Esc")
	}

	if cmd == nil {
		t.Error("Esc should return quit command")
	}

	// Test Ctrl+C
	m = NewModel(projects)
	keyMsg = tea.KeyMsg{Type: tea.KeyCtrlC}
	newModel, cmd = m.Update(keyMsg)
	m = newModel.(Model)

	if !m.quitting {
		t.Error("quitting should be true after Ctrl+C")
	}

	if cmd == nil {
		t.Error("Ctrl+C should return quit command")
	}
}

func TestApplyFilterWithEmptyQuery(t *testing.T) {
	projects := []status.ProjectDisplay{
		makePD("repo1", "/path/repo1", "dev"),
		makePD("repo2", "/path/repo2", "dev"),
	}

	m := NewModel(projects)
	m.textInput.SetValue("")
	m.applyFilter()

	if len(m.filteredProjects) != len(projects) {
		t.Error("empty query should show all projects")
	}

	if m.cursor != 0 {
		t.Error("cursor should reset to 0")
	}
}

func TestApplyFilterWithQuery(t *testing.T) {
	projects := []status.ProjectDisplay{
		makePD("user1/repo1", "/path/repo1", "dev"),
		makePD("user1/repo2", "/path/repo2", "dev"),
		makePD("user2/repo3", "/path/repo3", "dev"),
	}

	m := NewModel(projects)
	m.textInput.SetValue("user1")
	m.applyFilter()

	if len(m.filteredProjects) != 2 {
		t.Errorf("expected 2 filtered projects, got %d", len(m.filteredProjects))
	}
}

func TestMatchesQueryWithSlash(t *testing.T) {
	project := makePD("user/repo", "/path/to/repo", "dev")
	m := NewModel([]status.ProjectDisplay{project})

	testCases := []struct {
		query    string
		expected bool
	}{
		{"user/repo", true},         // スラッシュあり: repoに含まれる
		{"user/re", true},           // スラッシュあり: repoに含まれる
		{"/path/to", true},          // スラッシュあり: pathに含まれる
		{"path/to/repo", true},      // スラッシュあり: pathに含まれる
		{"to/repo", true},           // スラッシュあり: pathに含まれる
		{"nonexistent/test", false}, // スラッシュあり: どこにも含まれない
		// スラッシュなしのクエリは matchesSimpleQuery が使われるため除外
	}

	for _, tc := range testCases {
		result := m.matchesQuery(project, tc.query)
		if result != tc.expected {
			t.Errorf("matchesQuery with %q: expected %v, got %v", tc.query, tc.expected, result)
		}
	}
}

func TestMatchesQueryWithoutSlash(t *testing.T) {
	// スラッシュなしのクエリは別のロジック（matchesSimpleQuery）を使うため、
	// 別のテストケースとして分離
	project := makePD("user/repo", "/path/to/repo", "dev")
	m := NewModel([]status.ProjectDisplay{project})

	testCases := []struct {
		query    string
		expected bool
	}{
		{"repo", true}, // repo名にマッチ
		{"user", true}, // owner名にマッチ
		{"re", true},   // repo名の一部にマッチ
		{"us", true},   // owner名の一部にマッチ
		{"dev", false}, // workspaceはmatchesSimpleQueryでは検索対象外
		{"nonexist", false},
	}

	for _, tc := range testCases {
		result := m.matchesQuery(project, tc.query)
		if result != tc.expected {
			t.Errorf("matchesQuery (no slash) with %q: expected %v, got %v", tc.query, tc.expected, result)
		}
	}
}

func TestMatchesSimpleQuery(t *testing.T) {
	m := NewModel([]status.ProjectDisplay{})

	testCases := []struct {
		repoLower string
		query     string
		expected  bool
	}{
		{"user/repo", "repo", true},
		{"user/repo", "user", true},
		{"user/repo", "re", true},
		{"user/repo", "us", true},
		{"user/repo", "nonexist", false},
		{"simple", "simple", true},
		{"simple", "sim", true},
	}

	for _, tc := range testCases {
		result := m.matchesSimpleQuery(tc.repoLower, tc.query)
		if result != tc.expected {
			t.Errorf("matchesSimpleQuery(%q, %q): expected %v, got %v",
				tc.repoLower, tc.query, tc.expected, result)
		}
	}
}

func TestViewWithQuitting(t *testing.T) {
	projects := []status.ProjectDisplay{
		makePD("repo1", "/path/repo1", "dev"),
	}

	m := NewModel(projects)
	m.quitting = true

	view := m.View()
	if view != "" {
		t.Error("view should be empty when quitting")
	}
}

func TestViewWithNoProjects(t *testing.T) {
	m := NewModel([]status.ProjectDisplay{})

	view := m.View()
	if view == "" {
		t.Error("view should not be empty even with no projects")
	}

	if !strings.Contains(view, "Select a project") && !strings.Contains(view, "プロジェクトを選択") {
		t.Error("view should contain title")
	}
}

func TestViewWithNoMatches(t *testing.T) {
	projects := []status.ProjectDisplay{
		makePD("repo1", "/path/repo1", "dev"),
	}

	m := NewModel(projects)
	m.textInput.SetValue("nonexistent")
	m.applyFilter()

	view := m.View()
	if !strings.Contains(view, "No matching") && !strings.Contains(view, "一致する") {
		t.Error("view should show no matches message")
	}
}

func TestRenderProjectItem(t *testing.T) {
	projects := []status.ProjectDisplay{
		makePD("user/repo", "/path/repo", "dev"),
	}

	m := NewModel(projects)

	var s strings.Builder

	// Test unselected
	m.renderProjectItem(&s, 0, projects[0])
	result := s.String()
	if result == "" {
		t.Error("rendered item should not be empty")
	}
	if !strings.Contains(result, "user/repo") {
		t.Error("rendered item should contain repo name")
	}

	// Test selected
	s.Reset()
	m.cursor = 0
	m.renderProjectItem(&s, 0, projects[0])
	result = s.String()
	if !strings.Contains(result, "❯") {
		t.Error("selected item should contain cursor indicator")
	}
}

func TestRunWithEmptyProjects(t *testing.T) {
	result, err := Run([]status.ProjectDisplay{})
	if err != nil {
		t.Errorf("Run should not error with empty projects: %v", err)
	}
	if result != "" {
		t.Error("result should be empty for empty projects")
	}
}

func TestMoveCursorUpDown(t *testing.T) {
	projects := []status.ProjectDisplay{
		makePD("repo1", "/path/repo1", "dev"),
		makePD("repo2", "/path/repo2", "dev"),
		makePD("repo3", "/path/repo3", "dev"),
	}

	m := NewModel(projects)

	// Test moveCursorDown
	m.cursor = 0
	m.moveCursorDown()
	if m.cursor != 1 {
		t.Error("moveCursorDown should increment cursor")
	}

	// Test moveCursorDown at end
	m.cursor = 2
	m.moveCursorDown()
	if m.cursor != 0 {
		t.Error("moveCursorDown should wrap at end")
	}

	// Test moveCursorUp
	m.cursor = 1
	m.moveCursorUp()
	if m.cursor != 0 {
		t.Error("moveCursorUp should decrement cursor")
	}

	// Test moveCursorUp at start
	m.cursor = 0
	m.moveCursorUp()
	if m.cursor != 2 {
		t.Error("moveCursorUp should wrap at start")
	}
}

func TestUpdateSearchInput(t *testing.T) {
	projects := []status.ProjectDisplay{
		makePD("repo1", "/path/repo1", "dev"),
		makePD("repo2", "/path/repo2", "dev"),
	}

	m := NewModel(projects)

	// Simulate typing
	keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}}
	newModel, _ := m.updateSearchInput(keyMsg)
	m = newModel.(Model)

	if !strings.Contains(m.textInput.Value(), "r") {
		t.Error("text input should contain typed character")
	}
}

func TestRenderHeader(t *testing.T) {
	projects := []status.ProjectDisplay{
		makePD("repo1", "/path/repo1", "dev"),
	}

	m := NewModel(projects)
	var s strings.Builder
	m.renderHeader(&s)

	result := s.String()
	if result == "" {
		t.Error("header should not be empty")
	}
}

func TestRenderSearchInput(t *testing.T) {
	projects := []status.ProjectDisplay{
		makePD("repo1", "/path/repo1", "dev"),
		makePD("repo2", "/path/repo2", "dev"),
	}

	m := NewModel(projects)
	var s strings.Builder
	m.renderSearchInput(&s)

	result := s.String()
	if result == "" {
		t.Error("search input should not be empty")
	}
	if !strings.Contains(result, "[2/2]") {
		t.Error("should show match count")
	}
}

func TestRenderProjectList(t *testing.T) {
	projects := []status.ProjectDisplay{
		makePD("repo1", "/path/repo1", "dev"),
	}

	m := NewModel(projects)
	var s strings.Builder
	m.renderProjectList(&s)

	result := s.String()
	if result == "" {
		t.Error("project list should not be empty")
	}
}

func TestRenderFooter(t *testing.T) {
	projects := []status.ProjectDisplay{
		makePD("repo1", "/path/repo1", "dev"),
	}

	m := NewModel(projects)
	var s strings.Builder
	m.renderFooter(&s)

	result := s.String()
	if result == "" {
		t.Error("footer should not be empty")
	}
}
