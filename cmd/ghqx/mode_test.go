package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mi8bi/ghqx/internal/app"
	"github.com/mi8bi/ghqx/internal/config"

	tea "github.com/charmbracelet/bubbletea"
)

func TestRunModeWithLoadAppError(t *testing.T) {
	oldConfigPath := configPath
	configPath = "/nonexistent/config.toml"
	defer func() { configPath = oldConfigPath }()

	oldApp := application
	application = nil
	defer func() { application = oldApp }()

	err := runMode(modeCmd, []string{})
	if err == nil {
		t.Fatalf("expected error when loadApp fails")
	}
}

func TestRunModeWithNoRoots(t *testing.T) {
	tmp, err := os.MkdirTemp("", "ghqx-mode-noroots")
	if err != nil {
		t.Fatalf("tempdir: %v", err)
	}
	defer os.RemoveAll(tmp)

	cfgPath := filepath.Join(tmp, "config.toml")
	cfg := &config.Config{
		Roots:   map[string]string{}, // Empty roots
		Default: config.DefaultConfig{Root: ""},
	}

	// This should fail validation, but let's create it directly
	if err := os.WriteFile(cfgPath, []byte("[roots]\n[default]\nroot = \"\"\n"), 0644); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	oldConfigPath := configPath
	configPath = cfgPath
	defer func() { configPath = oldConfigPath }()

	// Create app with empty roots (bypassing validation)
	cfg.Roots = map[string]string{} // Empty
	appInstance := app.New(cfg)
	application = appInstance

	// This should return an error about no roots
	err = runMode(modeCmd, []string{})
	if err == nil {
		t.Fatalf("expected error when no roots defined")
	}
}

func TestModeSelectorModelInit(t *testing.T) {
	model := ModeSelectorModel{
		workspaceNames: []string{"dev", "sandbox", "release"},
		cursor:         0,
	}

	cmd := model.Init()
	if cmd != nil {
		t.Error("Init should return nil")
	}
}

func TestModeSelectorModelUpdate(t *testing.T) {
	model := ModeSelectorModel{
		workspaceNames: []string{"dev", "sandbox", "release"},
		cursor:         0,
	}

	// Test up movement
	msg := tea.KeyMsg{Type: tea.KeyUp}
	updatedModel, _ := model.Update(msg)
	m := updatedModel.(ModeSelectorModel)
	if m.cursor != 0 {
		t.Errorf("cursor should stay at 0 when at top, got %d", m.cursor)
	}

	// Test down movement
	msg = tea.KeyMsg{Type: tea.KeyDown}
	updatedModel, _ = model.Update(msg)
	m = updatedModel.(ModeSelectorModel)
	if m.cursor != 1 {
		t.Errorf("cursor should move to 1, got %d", m.cursor)
	}

	// Test down movement with k
	msg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}}
	updatedModel, _ = m.Update(msg)
	m = updatedModel.(ModeSelectorModel)
	if m.cursor != 0 {
		t.Errorf("cursor should move up with 'k', got %d", m.cursor)
	}

	// Test down movement with j
	msg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}}
	updatedModel, _ = m.Update(msg)
	m = updatedModel.(ModeSelectorModel)
	if m.cursor != 1 {
		t.Errorf("cursor should move down with 'j', got %d", m.cursor)
	}

	// Test enter selection
	msg = tea.KeyMsg{Type: tea.KeyEnter}
	updatedModel, cmd := m.Update(msg)
	m = updatedModel.(ModeSelectorModel)
	if m.selected == "" {
		t.Error("selection should be set after Enter")
	}
	if cmd == nil {
		t.Error("should return tea.Quit command")
	}

	// Test quit with q
	model = ModeSelectorModel{
		workspaceNames: []string{"dev"},
		cursor:         0,
	}
	msg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}
	updatedModel, cmd = model.Update(msg)
	m = updatedModel.(ModeSelectorModel)
	if !m.quitting {
		t.Error("should set quitting to true")
	}
	if cmd == nil {
		t.Error("should return tea.Quit command")
	}

	// Test Ctrl+C
	model = ModeSelectorModel{
		workspaceNames: []string{"dev"},
		cursor:         0,
	}
	msg = tea.KeyMsg{Type: tea.KeyCtrlC}
	updatedModel, cmd = model.Update(msg)
	m = updatedModel.(ModeSelectorModel)
	if !m.quitting {
		t.Error("should set quitting to true on Ctrl+C")
	}

	// Test Esc
	model = ModeSelectorModel{
		workspaceNames: []string{"dev"},
		cursor:         0,
	}
	msg = tea.KeyMsg{Type: tea.KeyEsc}
	updatedModel, cmd = model.Update(msg)
	m = updatedModel.(ModeSelectorModel)
	if !m.quitting {
		t.Error("should set quitting to true on Esc")
	}
}

func TestModeSelectorModelView(t *testing.T) {
	model := ModeSelectorModel{
		workspaceNames: []string{"dev", "sandbox", "release"},
		cursor:         1,
		quitting:       false,
	}

	view := model.View()
	if view == "" {
		t.Error("View should return non-empty string")
	}

	// Test quitting state
	model.quitting = true
	view = model.View()
	if view != "" {
		t.Error("View should return empty string when quitting")
	}
}

func TestModeSelectorModelBoundaries(t *testing.T) {
	model := ModeSelectorModel{
		workspaceNames: []string{"dev", "sandbox"},
		cursor:         1,
	}

	// Test moving down at bottom
	msg := tea.KeyMsg{Type: tea.KeyDown}
	updatedModel, _ := model.Update(msg)
	m := updatedModel.(ModeSelectorModel)
	if m.cursor != 1 {
		t.Errorf("cursor should stay at bottom (1), got %d", m.cursor)
	}

	// Test moving up from bottom
	msg = tea.KeyMsg{Type: tea.KeyUp}
	updatedModel, _ = m.Update(msg)
	m = updatedModel.(ModeSelectorModel)
	if m.cursor != 0 {
		t.Errorf("cursor should move to 0, got %d", m.cursor)
	}
}

func TestRunModeWithSingleRoot(t *testing.T) {
	tmp, err := os.MkdirTemp("", "ghqx-mode-single")
	if err != nil {
		t.Fatalf("tempdir: %v", err)
	}
	defer os.RemoveAll(tmp)

	cfgPath := filepath.Join(tmp, "config.toml")
	cfg := &config.Config{
		Roots:   map[string]string{"dev": filepath.Join(tmp, "dev")},
		Default: config.DefaultConfig{Root: "dev"},
	}

	loader := config.NewLoader()
	if err := loader.Save(cfg, cfgPath); err != nil {
		t.Fatalf("failed to save config: %v", err)
	}

	oldConfigPath := configPath
	configPath = cfgPath
	defer func() { configPath = oldConfigPath }()

	appInstance := app.New(cfg)
	application = appInstance

	// We can't easily test the full TUI flow, but we can verify setup
	rootNames := []string{}
	for name := range application.Config.Roots {
		rootNames = append(rootNames, name)
	}

	if len(rootNames) != 1 {
		t.Errorf("expected 1 root, got %d", len(rootNames))
	}
}
