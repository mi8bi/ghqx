package configtui

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mi8bi/ghqx/internal/config"
)

func TestRun(t *testing.T) {
	tmp, err := os.MkdirTemp("", "ghqx-config-run")
	if err != nil {
		t.Fatalf("tempdir: %v", err)
	}
	defer os.RemoveAll(tmp)

	cfgPath := filepath.Join(tmp, "config.toml")
	cfg := &config.Config{
		Roots:   map[string]string{"dev": "/tmp/dev"},
		Default: config.DefaultConfig{Root: "dev"},
	}

	// We can't fully test the interactive TUI, but we can verify setup
	// The Run function will create a tea.Program which requires terminal
	// For unit tests, we verify the function signature and basic setup

	// Verify function exists and accepts correct parameters
	_ = Run

	// Test that we can create a model
	model := NewModel(cfg, cfgPath)
	if model.editor == nil {
		t.Fatal("model editor should not be nil")
	}

	if model.editor.ConfigPath != cfgPath {
		t.Error("config path mismatch")
	}
}

func TestRunWithValidConfig(t *testing.T) {
	tmp, err := os.MkdirTemp("", "ghqx-config-valid")
	if err != nil {
		t.Fatalf("tempdir: %v", err)
	}
	defer os.RemoveAll(tmp)

	cfgPath := filepath.Join(tmp, "config.toml")
	cfg := &config.Config{
		Roots: map[string]string{
			"dev":     filepath.Join(tmp, "dev"),
			"release": filepath.Join(tmp, "release"),
			"sandbox": filepath.Join(tmp, "sandbox"),
		},
		Default: config.DefaultConfig{Root: "dev"},
	}

	// Create root directories
	for _, path := range cfg.Roots {
		if err := os.MkdirAll(path, 0755); err != nil {
			t.Fatalf("mkdir: %v", err)
		}
	}

	// Save config
	loader := config.NewLoader()
	if err := loader.Save(cfg, cfgPath); err != nil {
		t.Fatalf("save config: %v", err)
	}

	// Create model
	model := NewModel(cfg, cfgPath)

	// Verify setup
	if len(model.editor.Fields) == 0 {
		t.Error("fields should be populated")
	}

	if model.state != EditStateList {
		t.Error("initial state should be EditStateList")
	}
}
