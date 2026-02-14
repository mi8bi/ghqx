package configtui

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mi8bi/ghqx/internal/config"
)

func TestNewModel(t *testing.T) {
	cfg := &config.Config{
		Roots:   map[string]string{"dev": "/tmp/dev"},
		Default: config.DefaultConfig{Root: "dev"},
	}

	model := NewModel(cfg, "/tmp/config.toml")

	if model.editor == nil {
		t.Fatal("editor should not be nil")
	}

	if model.cursor != 0 {
		t.Error("cursor should be 0 initially")
	}

	if model.state != EditStateList {
		t.Error("state should be EditStateList initially")
	}

	if model.messageType != MessageTypeNone {
		t.Error("messageType should be MessageTypeNone initially")
	}
}

func TestModelInit(t *testing.T) {
	cfg := &config.Config{
		Roots:   map[string]string{"dev": "/tmp/dev"},
		Default: config.DefaultConfig{Root: "dev"},
	}

	model := NewModel(cfg, "/tmp/config.toml")
	cmd := model.Init()

	if cmd != nil {
		t.Error("Init should return nil")
	}
}

func TestModelUpdateWithKeys(t *testing.T) {
	cfg := &config.Config{
		Roots: map[string]string{
			"dev":     "/tmp/dev",
			"release": "/tmp/release",
			"sandbox": "/tmp/sandbox",
		},
		Default: config.DefaultConfig{Root: "dev"},
	}

	model := NewModel(cfg, "/tmp/config.toml")

	// Test navigation down
	keyMsg := tea.KeyMsg{Type: tea.KeyDown}
	newModel, _ := model.Update(keyMsg)
	model = newModel.(Model)
	if model.cursor != 1 {
		t.Errorf("cursor should be 1 after down, got %d", model.cursor)
	}

	// Test navigation up
	keyMsg = tea.KeyMsg{Type: tea.KeyUp}
	newModel, _ = model.Update(keyMsg)
	model = newModel.(Model)
	if model.cursor != 0 {
		t.Errorf("cursor should be 0 after up, got %d", model.cursor)
	}

	// Test j/k navigation
	keyMsg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}}
	newModel, _ = model.Update(keyMsg)
	model = newModel.(Model)
	if model.cursor != 1 {
		t.Error("j should move cursor down")
	}

	keyMsg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}}
	newModel, _ = model.Update(keyMsg)
	model = newModel.(Model)
	if model.cursor != 0 {
		t.Error("k should move cursor up")
	}

	// Test entering edit mode
	keyMsg = tea.KeyMsg{Type: tea.KeyEnter}
	newModel, _ = model.Update(keyMsg)
	model = newModel.(Model)
	if model.state != EditStateEdit {
		t.Error("Enter should enter edit mode")
	}

	// Test exiting edit mode
	keyMsg = tea.KeyMsg{Type: tea.KeyEsc}
	newModel, _ = model.Update(keyMsg)
	model = newModel.(Model)
	if model.state != EditStateList {
		t.Error("Esc should exit edit mode")
	}

	// Test quit with unsaved changes
	model.editor.Modified = true
	keyMsg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}
	newModel, cmd := model.Update(keyMsg)
	model = newModel.(Model)
	if cmd != nil {
		t.Error("q should not quit when there are unsaved changes")
	}
	if model.message == "" {
		t.Error("should show warning about unsaved changes")
	}

	// Test force quit
	keyMsg = tea.KeyMsg{Type: tea.KeyCtrlQ}
	_, cmd = model.Update(keyMsg)
	if cmd == nil {
		t.Error("Ctrl+Q should quit")
	}

	// Test save command
	model = NewModel(cfg, "/tmp/config.toml")
	keyMsg = tea.KeyMsg{Type: tea.KeyCtrlS}
	_, cmd = model.Update(keyMsg)
	if cmd == nil {
		t.Error("Ctrl+S should trigger save")
	}
}

func TestModelUpdateEditMode(t *testing.T) {
	cfg := &config.Config{
		Roots: map[string]string{
			"dev":     "/tmp/dev",
			"release": "/tmp/release",
			"sandbox": "/tmp/sandbox",
		},
		Default: config.DefaultConfig{Root: "dev"},
	}

	model := NewModel(cfg, "/tmp/config.toml")
	model.state = EditStateEdit
	model.cursor = 0
	model.editValue = "/tmp/test"

	// Test text input (for string fields)
	keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}}
	newModel, _ := model.Update(keyMsg)
	model = newModel.(Model)
	if !strings.Contains(model.editValue, "a") {
		t.Error("should append character to edit value")
	}

	// Test backspace
	keyMsg = tea.KeyMsg{Type: tea.KeyBackspace}
	newModel, _ = model.Update(keyMsg)
	model = newModel.(Model)
	// editValue should be shorter

	// Test Enter to confirm
	keyMsg = tea.KeyMsg{Type: tea.KeyEnter}
	newModel, _ = model.Update(keyMsg)
	model = newModel.(Model)
	if model.state != EditStateList {
		t.Error("Enter should confirm and exit edit mode")
	}

	// Test selection field
	model.state = EditStateEdit
	model.cursor = 3 // Default root field (selection type)
	model.editValue = "dev"

	// Test left arrow
	keyMsg = tea.KeyMsg{Type: tea.KeyLeft}
	newModel, _ = model.Update(keyMsg)
	model = newModel.(Model)
	// editValue should cycle

	// Test right arrow
	keyMsg = tea.KeyMsg{Type: tea.KeyRight}
	newModel, _ = model.Update(keyMsg)
	model = newModel.(Model)
	// editValue should cycle

	// Test space for selection
	keyMsg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{' '}}
	newModel, _ = model.Update(keyMsg)
	model = newModel.(Model)
	// editValue should cycle
}

func TestModelUpdateWithMessages(t *testing.T) {
	cfg := &config.Config{
		Roots:   map[string]string{"dev": "/tmp/dev"},
		Default: config.DefaultConfig{Root: "dev"},
	}

	model := NewModel(cfg, "/tmp/config.toml")

	// Test saveSuccessMsg
	msg := saveSuccessMsg{}
	newModel, _ := model.Update(msg)
	model = newModel.(Model)

	if model.state != EditStateList {
		t.Error("state should be EditStateList after save success")
	}
	if model.messageType != MessageTypeSuccess {
		t.Error("messageType should be Success")
	}
	if model.editor.Modified {
		t.Error("Modified should be false after save")
	}

	// Test saveErrorMsg
	model = NewModel(cfg, "/tmp/config.toml")
	errMsg := saveErrorMsg{err: os.ErrNotExist}
	newModel, _ = model.Update(errMsg)
	model = newModel.(Model)

	if model.state != EditStateList {
		t.Error("state should be EditStateList after save error")
	}
	if model.messageType != MessageTypeError {
		t.Error("messageType should be Error")
	}

	// Test WindowSizeMsg
	model = NewModel(cfg, "/tmp/config.toml")
	sizeMsg := tea.WindowSizeMsg{Width: 100, Height: 50}
	newModel, _ = model.Update(sizeMsg)
	model = newModel.(Model)

	if model.width != 100 || model.height != 50 {
		t.Error("window size not updated")
	}
}

func TestModelView(t *testing.T) {
	cfg := &config.Config{
		Roots: map[string]string{
			"dev":     "/tmp/dev",
			"release": "/tmp/release",
			"sandbox": "/tmp/sandbox",
		},
		Default: config.DefaultConfig{Root: "dev"},
	}

	model := NewModel(cfg, "/tmp/config.toml")

	// Test list view
	view := model.View()
	if view == "" {
		t.Error("view should not be empty")
	}
	if !strings.Contains(view, "ghqx config edit") {
		t.Error("view should contain title")
	}

	// Test edit mode view
	model.state = EditStateEdit
	model.editValue = "test"
	view = model.View()
	if view == "" {
		t.Error("edit view should not be empty")
	}

	// Test saving state view
	model.state = EditStateSaving
	view = model.View()
	if view == "" {
		t.Error("saving view should not be empty")
	}

	// Test with message
	model.state = EditStateList
	model.message = "test message"
	model.messageType = MessageTypeInfo
	view = model.View()
	if !strings.Contains(view, "test message") {
		t.Error("view should contain message")
	}

	// Test with modified flag
	model.editor.Modified = true
	view = model.View()
	if !strings.Contains(view, "変更あり") {
		t.Error("view should show modified indicator")
	}
}

func TestRenderField(t *testing.T) {
	cfg := &config.Config{
		Roots:   map[string]string{"dev": "/tmp/dev"},
		Default: config.DefaultConfig{Root: "dev"},
	}

	model := NewModel(cfg, "/tmp/config.toml")
	field := model.editor.Fields[0]

	// Test unselected field
	rendered := model.renderField(field, false)
	if rendered == "" {
		t.Error("rendered field should not be empty")
	}
	if !strings.Contains(rendered, field.Name) {
		t.Error("rendered field should contain field name")
	}

	// Test selected field
	rendered = model.renderField(field, true)
	if rendered == "" {
		t.Error("rendered selected field should not be empty")
	}

	// Test field in edit mode
	model.state = EditStateEdit
	model.editValue = "test value"
	rendered = model.renderField(field, true)
	if rendered == "" {
		t.Error("rendered edit field should not be empty")
	}
}

func TestRenderHelp(t *testing.T) {
	cfg := &config.Config{
		Roots:   map[string]string{"dev": "/tmp/dev"},
		Default: config.DefaultConfig{Root: "dev"},
	}

	model := NewModel(cfg, "/tmp/config.toml")

	// Test list mode help
	help := model.renderHelp()
	if help == "" {
		t.Error("help should not be empty")
	}

	// Test edit mode help for string field
	model.state = EditStateEdit
	model.cursor = 0
	help = model.renderHelp()
	if help == "" {
		t.Error("edit help should not be empty")
	}

	// Test edit mode help for selection field
	model.cursor = 3
	help = model.renderHelp()
	if help == "" {
		t.Error("selection help should not be empty")
	}

	// Test help with modified flag
	model.state = EditStateList
	model.editor.Modified = true
	help = model.renderHelp()
	if !strings.Contains(help, "Ctrl+Q") {
		t.Error("help should show Ctrl+Q when modified")
	}
}

func TestHandleListKeysBoundaries(t *testing.T) {
	cfg := &config.Config{
		Roots: map[string]string{
			"dev":     "/tmp/dev",
			"release": "/tmp/release",
		},
		Default: config.DefaultConfig{Root: "dev"},
	}

	model := NewModel(cfg, "/tmp/config.toml")
	model.cursor = 0

	// Test up at top (should stay at 0)
	keyMsg := tea.KeyMsg{Type: tea.KeyUp}
	newModel, _ := model.Update(keyMsg)
	model = newModel.(Model)
	if model.cursor != 0 {
		t.Error("cursor should stay at 0 when at top")
	}

	// Move to bottom
	model.cursor = len(model.editor.Fields) - 1

	// Test down at bottom (should stay at last)
	keyMsg = tea.KeyMsg{Type: tea.KeyDown}
	newModel, _ = model.Update(keyMsg)
	model = newModel.(Model)
	if model.cursor != len(model.editor.Fields)-1 {
		t.Error("cursor should stay at last when at bottom")
	}
}

func TestQuitWithoutModifications(t *testing.T) {
	cfg := &config.Config{
		Roots:   map[string]string{"dev": "/tmp/dev"},
		Default: config.DefaultConfig{Root: "dev"},
	}

	model := NewModel(cfg, "/tmp/config.toml")
	model.editor.Modified = false

	keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}
	_, cmd := model.Update(keyMsg)

	if cmd == nil {
		t.Error("q should quit when no modifications")
	}
}

func TestSaveWithValidPath(t *testing.T) {
	tmp, err := os.MkdirTemp("", "ghqx-model-save")
	if err != nil {
		t.Fatalf("tempdir: %v", err)
	}
	defer os.RemoveAll(tmp)

	cfgPath := filepath.Join(tmp, "config.toml")
	cfg := &config.Config{
		Roots:   map[string]string{"dev": filepath.Join(tmp, "dev")},
		Default: config.DefaultConfig{Root: "dev"},
	}

	model := NewModel(cfg, cfgPath)

	// Trigger save
	keyMsg := tea.KeyMsg{Type: tea.KeyCtrlS}
	_, cmd := model.Update(keyMsg)

	if cmd == nil {
		t.Fatal("Ctrl+S should return save command")
	}

	// Execute the save command
	msg := cmd()

	// Should succeed
	if _, ok := msg.(saveSuccessMsg); !ok {
		t.Errorf("expected saveSuccessMsg, got %T", msg)
	}
}
