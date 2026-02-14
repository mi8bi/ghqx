package tui

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mi8bi/ghqx/internal/app"
	"github.com/mi8bi/ghqx/internal/config"
	"github.com/mi8bi/ghqx/internal/domain"
	"github.com/mi8bi/ghqx/internal/status"
)

func TestNewStatusModel(t *testing.T) {
	cfg := &config.Config{
		Roots:   map[string]string{"sandbox": "/tmp"},
		Default: config.DefaultConfig{Root: "sandbox"},
	}
	appInstance := app.New(cfg)

	model := NewStatusModel(appInstance)

	if model.app == nil {
		t.Fatal("app should not be nil")
	}

	if len(model.projects) != 0 {
		t.Error("projects should be empty initially")
	}

	if model.cursor != 0 {
		t.Error("cursor should be 0 initially")
	}

	if model.viewState != ViewStateLoading {
		t.Error("viewState should be ViewStateLoading initially")
	}
}

func TestStatusModelInit(t *testing.T) {
	tmp, err := os.MkdirTemp("", "ghqx-tui-init")
	if err != nil {
		t.Fatalf("tempdir: %v", err)
	}
	defer os.RemoveAll(tmp)

	cfg := &config.Config{
		Roots:   map[string]string{"sandbox": tmp},
		Default: config.DefaultConfig{Root: "sandbox"},
	}
	appInstance := app.New(cfg)

	model := NewStatusModel(appInstance)
	cmd := model.Init()

	if cmd == nil {
		t.Fatal("Init should return a command")
	}
}

func TestStatusModelUpdateWithKeys(t *testing.T) {
	tmp, err := os.MkdirTemp("", "ghqx-tui-update")
	if err != nil {
		t.Fatalf("tempdir: %v", err)
	}
	defer os.RemoveAll(tmp)

	// Create test project
	repo := filepath.Join(tmp, "github.com", "user", "repo")
	if err := os.MkdirAll(repo, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	cfg := &config.Config{
		Roots:   map[string]string{"sandbox": tmp},
		Default: config.DefaultConfig{Root: "sandbox"},
	}
	appInstance := app.New(cfg)

	model := NewStatusModel(appInstance)

	// Load projects first
	model.viewState = ViewStateList
	model.projects = []ProjectRow{
		NewProjectRow(status.ProjectDisplay{
			Repo:      "user/repo",
			Workspace: "sandbox",
			FullPath:  repo,
		}),
	}

	// Test quit
	keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}
	newModel, cmd := model.Update(keyMsg)
	if cmd == nil {
		t.Error("quit should return tea.Quit command")
	}
	model = newModel.(StatusModel)

	// Reset for next test
	model = NewStatusModel(appInstance)
	model.viewState = ViewStateList
	model.projects = []ProjectRow{
		NewProjectRow(status.ProjectDisplay{Repo: "user/repo"}),
		NewProjectRow(status.ProjectDisplay{Repo: "user/repo2"}),
	}

	// Test down navigation
	keyMsg = tea.KeyMsg{Type: tea.KeyDown}
	newModel, _ = model.Update(keyMsg)
	model = newModel.(StatusModel)
	if model.cursor != 1 {
		t.Errorf("cursor should be 1 after down, got %d", model.cursor)
	}

	// Test up navigation
	keyMsg = tea.KeyMsg{Type: tea.KeyUp}
	newModel, _ = model.Update(keyMsg)
	model = newModel.(StatusModel)
	if model.cursor != 0 {
		t.Errorf("cursor should be 0 after up, got %d", model.cursor)
	}

	// Test j/k navigation
	keyMsg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}}
	newModel, _ = model.Update(keyMsg)
	model = newModel.(StatusModel)
	if model.cursor != 1 {
		t.Error("j should move cursor down")
	}

	keyMsg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}}
	newModel, _ = model.Update(keyMsg)
	model = newModel.(StatusModel)
	if model.cursor != 0 {
		t.Error("k should move cursor up")
	}

	// Test detail toggle
	keyMsg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}}
	newModel, _ = model.Update(keyMsg)
	model = newModel.(StatusModel)
	if !model.showDetail {
		t.Error("d should toggle detail view")
	}

	// Test refresh
	keyMsg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}}
	newModel, cmd = model.Update(keyMsg)
	model = newModel.(StatusModel)
	if model.viewState != ViewStateLoading {
		t.Error("r should trigger reload")
	}
	if cmd == nil {
		t.Error("r should return loadProjects command")
	}
}

func TestStatusModelUpdateWithMessages(t *testing.T) {
	cfg := &config.Config{
		Roots:   map[string]string{"sandbox": "/tmp"},
		Default: config.DefaultConfig{Root: "sandbox"},
	}
	appInstance := app.New(cfg)

	model := NewStatusModel(appInstance)

	// Test projectsLoadedMsg
	projects := []ProjectRow{
		NewProjectRow(status.ProjectDisplay{Repo: "test/repo"}),
	}
	msg := projectsLoadedMsg{projects: projects}
	newModel, _ := model.Update(msg)
	model = newModel.(StatusModel)

	if model.viewState != ViewStateList {
		t.Error("viewState should be ViewStateList after projectsLoadedMsg")
	}
	if len(model.projects) != 1 {
		t.Error("projects should be loaded")
	}

	// Test errorMsg with GhqxError
	ghqxErr := domain.NewError(domain.ErrCodeConfigNotFound, "test error").WithHint("test hint")
	errMsg := errorMsg{err: ghqxErr}
	newModel, _ = model.Update(errMsg)
	model = newModel.(StatusModel)

	if model.viewState != ViewStateError {
		t.Error("viewState should be ViewStateError after errorMsg")
	}
	if model.message == nil {
		t.Fatal("message should not be nil")
	}
	if model.message.Type != MessageTypeError {
		t.Error("message type should be Error")
	}

	// Test errorMsg with regular error
	model = NewStatusModel(appInstance)
	regularErr := os.ErrNotExist
	errMsg = errorMsg{err: regularErr}
	newModel, _ = model.Update(errMsg)
	model = newModel.(StatusModel)

	if model.viewState != ViewStateError {
		t.Error("viewState should be ViewStateError")
	}

	// Test WindowSizeMsg
	model = NewStatusModel(appInstance)
	sizeMsg := tea.WindowSizeMsg{Width: 100, Height: 50}
	newModel, _ = model.Update(sizeMsg)
	model = newModel.(StatusModel)

	if model.width != 100 || model.height != 50 {
		t.Error("window size not updated")
	}
}

func TestStatusModelView(t *testing.T) {
	cfg := &config.Config{
		Roots:   map[string]string{"sandbox": "/tmp"},
		Default: config.DefaultConfig{Root: "sandbox"},
	}
	appInstance := app.New(cfg)

	model := NewStatusModel(appInstance)

	// Test loading view
	model.viewState = ViewStateLoading
	view := model.View()
	if view == "" {
		t.Error("loading view should not be empty")
	}

	// Test error view
	model.viewState = ViewStateError
	model.err = os.ErrNotExist
	model.message = &Message{
		Text: "test error",
		Type: MessageTypeError,
		Hint: "test hint",
	}
	view = model.View()
	if view == "" {
		t.Error("error view should not be empty")
	}
	if !strings.Contains(view, "test error") {
		t.Error("error view should contain error message")
	}

	// Test list view
	model.viewState = ViewStateList
	model.projects = []ProjectRow{
		NewProjectRow(status.ProjectDisplay{
			Repo:       "user/repo",
			Workspace:  "sandbox",
			GitManaged: "Managed",
			Status:     "clean",
		}),
	}
	view = model.View()
	if view == "" {
		t.Error("list view should not be empty")
	}
	if !strings.Contains(view, "user/repo") {
		t.Error("list view should contain project name")
	}

	// Test detail view
	model.showDetail = true
	view = model.View()
	if view == "" {
		t.Error("detail view should not be empty")
	}
}

func TestRenderProjectRow(t *testing.T) {
	cfg := &config.Config{
		Roots:   map[string]string{"sandbox": "/tmp"},
		Default: config.DefaultConfig{Root: "sandbox"},
	}
	appInstance := app.New(cfg)

	model := NewStatusModel(appInstance)
	row := ProjectRow{
		ProjectDisplay: status.ProjectDisplay{
			Repo:       "user/repo",
			Workspace:  "sandbox",
			GitManaged: "Managed",
			Status:     "clean",
		},
	}

	// Test unselected row
	rendered := model.renderProjectRow(row, false)
	if rendered == "" {
		t.Error("rendered row should not be empty")
	}
	if !strings.Contains(rendered, "user/repo") {
		t.Error("rendered row should contain repo name")
	}

	// Test selected row
	rendered = model.renderProjectRow(row, true)
	if rendered == "" {
		t.Error("rendered selected row should not be empty")
	}
	if !strings.Contains(rendered, ">") {
		t.Error("selected row should contain cursor indicator")
	}
}

func TestRenderHelp(t *testing.T) {
	cfg := &config.Config{
		Roots:   map[string]string{"sandbox": "/tmp"},
		Default: config.DefaultConfig{Root: "sandbox"},
	}
	appInstance := app.New(cfg)

	model := NewStatusModel(appInstance)
	help := model.renderHelp()

	if help == "" {
		t.Error("help text should not be empty")
	}
}

func TestHandleKeyPressBoundaries(t *testing.T) {
	cfg := &config.Config{
		Roots:   map[string]string{"sandbox": "/tmp"},
		Default: config.DefaultConfig{Root: "sandbox"},
	}
	appInstance := app.New(cfg)

	model := NewStatusModel(appInstance)
	model.viewState = ViewStateList
	model.projects = []ProjectRow{
		NewProjectRow(status.ProjectDisplay{Repo: "repo1"}),
		NewProjectRow(status.ProjectDisplay{Repo: "repo2"}),
	}
	model.cursor = 0

	// Test up at top (should stay at 0)
	keyMsg := tea.KeyMsg{Type: tea.KeyUp}
	newModel, _ := model.Update(keyMsg)
	model = newModel.(StatusModel)
	if model.cursor != 0 {
		t.Error("cursor should stay at 0 when at top")
	}

	// Move to bottom
	model.cursor = 1

	// Test down at bottom (should stay at 1)
	keyMsg = tea.KeyMsg{Type: tea.KeyDown}
	newModel, _ = model.Update(keyMsg)
	model = newModel.(StatusModel)
	if model.cursor != 1 {
		t.Error("cursor should stay at 1 when at bottom")
	}
}

func TestCtrlCQuit(t *testing.T) {
	cfg := &config.Config{
		Roots:   map[string]string{"sandbox": "/tmp"},
		Default: config.DefaultConfig{Root: "sandbox"},
	}
	appInstance := app.New(cfg)

	model := NewStatusModel(appInstance)

	keyMsg := tea.KeyMsg{Type: tea.KeyCtrlC}
	_, cmd := model.Update(keyMsg)

	if cmd == nil {
		t.Error("Ctrl+C should return quit command")
	}
}
