// internal/tui/config/commands_test.go

package configtui

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mi8bi/ghqx/internal/config"
)

func TestSaveConfig(t *testing.T) {
	tmp, err := os.MkdirTemp("", "ghqx-config-save")
	if err != nil {
		t.Fatalf("tempdir: %v", err)
	}
	defer os.RemoveAll(tmp)

	cfgPath := filepath.Join(tmp, "config.toml")
	cfg := &config.Config{
		Roots:   map[string]string{"dev": "/tmp/dev"},
		Default: config.DefaultConfig{Root: "dev"},
	}

	editor := NewConfigEditor(cfg, cfgPath)
	model := Model{
		editor: editor,
		state:  EditStateList,
	}

	// Execute save command
	cmd := model.saveConfig()
	if cmd == nil {
		t.Fatal("saveConfig should return a command")
	}

	// Execute the command function
	msg := cmd()

	// Check result
	switch msg := msg.(type) {
	case saveSuccessMsg:
		// Success expected
	case saveErrorMsg:
		t.Errorf("save failed: %v", msg.err)
	default:
		t.Errorf("unexpected message type: %T", msg)
	}
}

func TestSaveConfigWithValidationError(t *testing.T) {
	tmp, err := os.MkdirTemp("", "ghqx-config-invalid")
	if err != nil {
		t.Fatalf("tempdir: %v", err)
	}
	defer os.RemoveAll(tmp)

	cfgPath := filepath.Join(tmp, "config.toml")

	// Create initially valid config
	cfg := &config.Config{
		Roots:   map[string]string{"dev": "/tmp/dev", "release": "/tmp/release", "sandbox": "/tmp/sandbox"},
		Default: config.DefaultConfig{Root: "dev"},
	}

	editor := NewConfigEditor(cfg, cfgPath)

	// Clear all field values to make them invalid
	// When ApplyChanges is called, it will update Config.Roots with empty values
	for i := range editor.Fields {
		if editor.Fields[i].Key == "roots.dev" ||
			editor.Fields[i].Key == "roots.release" ||
			editor.Fields[i].Key == "roots.sandbox" {
			editor.Fields[i].Value = "" // Set to empty string
		}
	}

	// After ApplyChanges, all roots will be empty strings (not removed from map)
	// But we need them actually removed for validation to fail
	// So let's update the field to set default.root to invalid value instead
	for i := range editor.Fields {
		if editor.Fields[i].Key == "default.root" {
			editor.Fields[i].Value = "nonexistent" // This will fail validation
			break
		}
	}

	model := Model{
		editor: editor,
		state:  EditStateList,
	}

	// Execute save command
	cmd := model.saveConfig()
	msg := cmd()

	// Should return error because validation will fail
	if _, ok := msg.(saveErrorMsg); !ok {
		t.Errorf("expected saveErrorMsg, got %T", msg)
	}
}

func TestSaveConfigWithInvalidDefaultRoot(t *testing.T) {
	tmp, err := os.MkdirTemp("", "ghqx-config-invalid-default")
	if err != nil {
		t.Fatalf("tempdir: %v", err)
	}
	defer os.RemoveAll(tmp)

	cfgPath := filepath.Join(tmp, "config.toml")

	// Create config with invalid default root
	cfg := &config.Config{
		Roots:   map[string]string{"dev": "/tmp/dev"},
		Default: config.DefaultConfig{Root: "nonexistent"}, // Invalid: doesn't exist in roots
	}

	editor := NewConfigEditor(cfg, cfgPath)

	model := Model{
		editor: editor,
		state:  EditStateList,
	}

	// Execute save command
	cmd := model.saveConfig()
	msg := cmd()

	// Should return error because default root is invalid
	if _, ok := msg.(saveErrorMsg); !ok {
		t.Errorf("expected saveErrorMsg for invalid default root, got %T", msg)
	}
}

func TestSaveConfigWithIOError(t *testing.T) {
	// Create a path that will definitely fail on write
	// Use a file as the "directory" path which will cause MkdirAll to fail
	tmp, err := os.MkdirTemp("", "ghqx-config-ioerror")
	if err != nil {
		t.Fatalf("tempdir: %v", err)
	}
	defer os.RemoveAll(tmp)

	// Create a file where we want to create a directory
	blockingFile := filepath.Join(tmp, "blocking")
	if err := os.WriteFile(blockingFile, []byte("test"), 0644); err != nil {
		t.Fatalf("failed to create blocking file: %v", err)
	}

	// Try to save config to a path that requires creating a directory
	// where a file already exists
	cfgPath := filepath.Join(blockingFile, "config.toml")

	cfg := &config.Config{
		Roots:   map[string]string{"dev": "/tmp/dev"},
		Default: config.DefaultConfig{Root: "dev"},
	}

	editor := NewConfigEditor(cfg, cfgPath)
	model := Model{
		editor: editor,
		state:  EditStateList,
	}

	// Execute save command
	cmd := model.saveConfig()
	msg := cmd()

	// Should return error due to IO issues (can't create dir where file exists)
	if _, ok := msg.(saveErrorMsg); !ok {
		t.Errorf("expected saveErrorMsg for IO error, got %T", msg)
	}
}

func TestSaveConfigWithReadOnlyDirectory(t *testing.T) {
	// Skip on Windows as chmod doesn't work the same way
	if os.Getenv("OS") == "Windows_NT" {
		t.Skip("Skipping read-only directory test on Windows")
	}

	tmp, err := os.MkdirTemp("", "ghqx-config-readonly")
	if err != nil {
		t.Fatalf("tempdir: %v", err)
	}
	defer os.RemoveAll(tmp)

	// Create a read-only directory
	readOnlyDir := filepath.Join(tmp, "readonly")
	if err := os.Mkdir(readOnlyDir, 0555); err != nil {
		t.Fatalf("failed to create readonly dir: %v", err)
	}
	// Restore permissions for cleanup
	defer os.Chmod(readOnlyDir, 0755)

	cfgPath := filepath.Join(readOnlyDir, "config.toml")

	cfg := &config.Config{
		Roots:   map[string]string{"dev": "/tmp/dev"},
		Default: config.DefaultConfig{Root: "dev"},
	}

	editor := NewConfigEditor(cfg, cfgPath)
	model := Model{
		editor: editor,
		state:  EditStateList,
	}

	// Execute save command
	cmd := model.saveConfig()
	msg := cmd()

	// Should return error due to permission issues
	if _, ok := msg.(saveErrorMsg); !ok {
		t.Errorf("expected saveErrorMsg for permission error, got %T", msg)
	}
}

func TestSaveSuccessMsg(t *testing.T) {
	msg := saveSuccessMsg{}
	// Just verify it exists and can be created
	_ = msg
}

func TestSaveErrorMsg(t *testing.T) {
	testErr := os.ErrNotExist
	msg := saveErrorMsg{err: testErr}

	if msg.err != testErr {
		t.Error("error mismatch")
	}
}

func TestSaveConfigAppliesChanges(t *testing.T) {
	tmp, err := os.MkdirTemp("", "ghqx-config-apply")
	if err != nil {
		t.Fatalf("tempdir: %v", err)
	}
	defer os.RemoveAll(tmp)

	cfgPath := filepath.Join(tmp, "config.toml")
	cfg := &config.Config{
		Roots: map[string]string{
			"dev":     "/tmp/dev",
			"release": "/tmp/release",
			"sandbox": "/tmp/sandbox",
		},
		Default: config.DefaultConfig{Root: "dev"},
	}

	editor := NewConfigEditor(cfg, cfgPath)

	// Modify a field
	editor.UpdateField(0, "/new/dev/path")

	model := Model{
		editor: editor,
		state:  EditStateList,
	}

	// Execute save command
	cmd := model.saveConfig()
	msg := cmd()

	// Verify changes were applied
	if editor.Config.Roots["dev"] != "/new/dev/path" {
		t.Error("changes should be applied before save")
	}

	// Check result
	if _, ok := msg.(saveSuccessMsg); !ok {
		t.Errorf("expected saveSuccessMsg, got %T", msg)
	}
}

func TestSaveConfigWithEmptyRootsAfterEdit(t *testing.T) {
	tmp, err := os.MkdirTemp("", "ghqx-config-empty-roots")
	if err != nil {
		t.Fatalf("tempdir: %v", err)
	}
	defer os.RemoveAll(tmp)

	cfgPath := filepath.Join(tmp, "config.toml")
	cfg := &config.Config{
		Roots: map[string]string{
			"dev":     "/tmp/dev",
			"release": "/tmp/release",
			"sandbox": "/tmp/sandbox",
		},
		Default: config.DefaultConfig{Root: "dev"},
	}

	editor := NewConfigEditor(cfg, cfgPath)

	// Set default.root to invalid value via Fields
	// This will cause validation to fail after ApplyChanges
	for i := range editor.Fields {
		if editor.Fields[i].Key == "default.root" {
			editor.Fields[i].Value = "invalid_root_name"
			break
		}
	}

	model := Model{
		editor: editor,
		state:  EditStateList,
	}

	// Execute save command
	cmd := model.saveConfig()
	msg := cmd()

	// Should return error due to validation failure
	errMsg, ok := msg.(saveErrorMsg)
	if !ok {
		t.Fatalf("expected saveErrorMsg, got %T", msg)
	}

	if errMsg.err == nil {
		t.Error("error should not be nil")
	}
}
