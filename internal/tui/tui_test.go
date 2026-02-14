package tui

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mi8bi/ghqx/internal/app"
	"github.com/mi8bi/ghqx/internal/config"
)

func TestRunStatus(t *testing.T) {
	tmp, err := os.MkdirTemp("", "ghqx-tui-run")
	if err != nil {
		t.Fatalf("tempdir: %v", err)
	}
	defer os.RemoveAll(tmp)

	cfg := &config.Config{
		Roots:   map[string]string{"sandbox": tmp},
		Default: config.DefaultConfig{Root: "sandbox"},
	}
	appInstance := app.New(cfg)

	// We can't fully test the interactive TUI, but we can verify it doesn't panic
	// In a real test environment, the TUI will exit immediately
	// This test mainly ensures the function signature and basic setup works

	// Note: Actual TUI testing would require mocking tea.Program or using test mode
	// For now, we just verify the function exists and accepts correct parameters
	if appInstance == nil {
		t.Fatal("app instance should not be nil")
	}

	// Cannot actually run the TUI in tests as it requires terminal
	// But we can test that the function signature is correct
	_ = RunStatus
}

func TestRunStatusSetup(t *testing.T) {
	// Test that we can set up the model correctly
	cfg := &config.Config{
		Roots:   map[string]string{"sandbox": "/tmp"},
		Default: config.DefaultConfig{Root: "sandbox"},
	}
	appInstance := app.New(cfg)

	model := NewStatusModel(appInstance)

	if model.app != appInstance {
		t.Error("model app should match provided app")
	}

	if model.viewState != ViewStateLoading {
		t.Error("initial viewState should be Loading")
	}
}

func TestRunStatusWithValidApp(t *testing.T) {
	tmp, err := os.MkdirTemp("", "ghqx-tui-valid")
	if err != nil {
		t.Fatalf("tempdir: %v", err)
	}
	defer os.RemoveAll(tmp)

	// Create test repository
	repo := filepath.Join(tmp, "github.com", "user", "repo")
	if err := os.MkdirAll(repo, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	cfg := &config.Config{
		Roots:   map[string]string{"sandbox": tmp},
		Default: config.DefaultConfig{Root: "sandbox"},
	}
	appInstance := app.New(cfg)

	// Verify app is set up correctly for TUI
	if appInstance.Status == nil {
		t.Fatal("app Status service should not be nil")
	}

	if appInstance.Config == nil {
		t.Fatal("app Config should not be nil")
	}
}
