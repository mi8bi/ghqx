package ghq

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/mi8bi/ghqx/internal/config"
)

func TestGetWithoutGhqReturnsError(t *testing.T) {
	cfg := &config.Config{Roots: map[string]string{"sandbox": "/tmp"}, Default: config.DefaultConfig{Root: "sandbox"}}
	c := NewClient(cfg)

	// Ensure ghq cannot be found by clearing PATH
	origPath := os.Getenv("PATH")
	os.Setenv("PATH", "")
	defer os.Setenv("PATH", origPath)

	err := c.Get(GetOptions{Repository: "github.com/user/repo", Workspace: "sandbox"})
	if err == nil {
		t.Fatalf("expected error when ghq is not available")
	}
}

// Additional tests for better coverage

func TestNewClient(t *testing.T) {
	cfg := &config.Config{
		Roots:   map[string]string{"sandbox": "/tmp"},
		Default: config.DefaultConfig{Root: "sandbox"},
	}
	c := NewClient(cfg)

	if c == nil {
		t.Fatal("NewClient returned nil")
	}

	if c.cfg != cfg {
		t.Error("config not set correctly")
	}

	if c.timeout == 0 {
		t.Error("timeout should be set")
	}

	if c.timeout != 30*time.Second {
		t.Errorf("expected timeout 30s, got %v", c.timeout)
	}
}

func TestGetWithInvalidWorkspace(t *testing.T) {
	cfg := &config.Config{
		Roots:   map[string]string{"sandbox": "/tmp/sandbox"},
		Default: config.DefaultConfig{Root: "sandbox"},
	}
	c := NewClient(cfg)

	err := c.Get(GetOptions{Repository: "test/repo", Workspace: "nonexistent"})
	if err == nil {
		t.Fatal("expected error when workspace doesn't exist")
	}
}

func TestGetWithValidWorkspaceButNoGhq(t *testing.T) {
	tmp, err := os.MkdirTemp("", "ghqx-ghq-test")
	if err != nil {
		t.Fatalf("tempdir: %v", err)
	}
	defer os.RemoveAll(tmp)

	sandboxPath := filepath.Join(tmp, "sandbox")
	if err := os.MkdirAll(sandboxPath, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	cfg := &config.Config{
		Roots:   map[string]string{"sandbox": sandboxPath},
		Default: config.DefaultConfig{Root: "sandbox"},
	}
	c := NewClient(cfg)

	// Clear PATH to ensure ghq is not available
	origPath := os.Getenv("PATH")
	os.Setenv("PATH", "")
	defer os.Setenv("PATH", origPath)

	err = c.Get(GetOptions{Repository: "test/repo", Workspace: "sandbox"})
	if err == nil {
		t.Fatal("expected error when ghq not available")
	}
}

func TestHasGhq(t *testing.T) {
	cfg := &config.Config{
		Roots:   map[string]string{"sandbox": "/tmp"},
		Default: config.DefaultConfig{Root: "sandbox"},
	}
	c := NewClient(cfg)

	// Save original PATH
	origPath := os.Getenv("PATH")
	defer os.Setenv("PATH", origPath)

	// Test with empty PATH (ghq shouldn't be found)
	os.Setenv("PATH", "")
	if c.hasGhq() {
		t.Error("hasGhq should return false when PATH is empty")
	}

	// Restore PATH
	os.Setenv("PATH", origPath)
}

func TestHasGhqWithAvailableGhq(t *testing.T) {
	// Skip if ghq is not available
	if _, err := exec.LookPath("ghq"); err != nil {
		t.Skip("ghq not available, skipping test")
	}

	cfg := &config.Config{
		Roots:   map[string]string{"sandbox": "/tmp"},
		Default: config.DefaultConfig{Root: "sandbox"},
	}
	c := NewClient(cfg)

	if !c.hasGhq() {
		t.Error("hasGhq should return true when ghq is available")
	}
}

func TestGetOptionsStruct(t *testing.T) {
	opts := GetOptions{
		Repository: "github.com/user/repo",
		Workspace:  "dev",
	}

	if opts.Repository != "github.com/user/repo" {
		t.Errorf("unexpected repository: %s", opts.Repository)
	}

	if opts.Workspace != "dev" {
		t.Errorf("unexpected workspace: %s", opts.Workspace)
	}
}

func TestGetWithDifferentWorkspaces(t *testing.T) {
	tmp, err := os.MkdirTemp("", "ghqx-ghq-workspaces")
	if err != nil {
		t.Fatalf("tempdir: %v", err)
	}
	defer os.RemoveAll(tmp)

	cfg := &config.Config{
		Roots: map[string]string{
			"sandbox": filepath.Join(tmp, "sandbox"),
			"dev":     filepath.Join(tmp, "dev"),
			"release": filepath.Join(tmp, "release"),
		},
		Default: config.DefaultConfig{Root: "sandbox"},
	}

	// Create root directories
	for _, path := range cfg.Roots {
		if err := os.MkdirAll(path, 0755); err != nil {
			t.Fatalf("mkdir: %v", err)
		}
	}

	c := NewClient(cfg)

	// Clear PATH to avoid actual ghq execution
	origPath := os.Getenv("PATH")
	os.Setenv("PATH", "")
	defer os.Setenv("PATH", origPath)

	testCases := []string{"sandbox", "dev", "release"}

	for _, workspace := range testCases {
		err := c.Get(GetOptions{
			Repository: "test/repo",
			Workspace:  workspace,
		})
		// Should fail because ghq is not available, but not because workspace doesn't exist
		if err == nil {
			t.Errorf("expected error for workspace %s", workspace)
		}
	}
}

func TestGetWithEmptyRepository(t *testing.T) {
	tmp, err := os.MkdirTemp("", "ghqx-ghq-empty-repo")
	if err != nil {
		t.Fatalf("tempdir: %v", err)
	}
	defer os.RemoveAll(tmp)

	cfg := &config.Config{
		Roots:   map[string]string{"sandbox": tmp},
		Default: config.DefaultConfig{Root: "sandbox"},
	}
	c := NewClient(cfg)

	// Even with empty repository, should check for ghq first
	origPath := os.Getenv("PATH")
	os.Setenv("PATH", "")
	defer os.Setenv("PATH", origPath)

	err = c.Get(GetOptions{
		Repository: "",
		Workspace:  "sandbox",
	})

	if err == nil {
		t.Fatal("expected error with empty repository")
	}
}

func TestClientTimeout(t *testing.T) {
	cfg := &config.Config{
		Roots:   map[string]string{"sandbox": "/tmp"},
		Default: config.DefaultConfig{Root: "sandbox"},
	}
	c := NewClient(cfg)

	expectedTimeout := 30 * time.Second
	if c.timeout != expectedTimeout {
		t.Errorf("expected timeout %v, got %v", expectedTimeout, c.timeout)
	}
}

func TestGetWithSpecialCharactersInRepository(t *testing.T) {
	tmp, err := os.MkdirTemp("", "ghqx-ghq-special")
	if err != nil {
		t.Fatalf("tempdir: %v", err)
	}
	defer os.RemoveAll(tmp)

	cfg := &config.Config{
		Roots:   map[string]string{"sandbox": tmp},
		Default: config.DefaultConfig{Root: "sandbox"},
	}
	c := NewClient(cfg)

	// Clear PATH
	origPath := os.Getenv("PATH")
	os.Setenv("PATH", "")
	defer os.Setenv("PATH", origPath)

	// Test with URL-like repository
	err = c.Get(GetOptions{
		Repository: "https://github.com/user/repo.git",
		Workspace:  "sandbox",
	})

	// Should fail because ghq is not available
	if err == nil {
		t.Fatal("expected error when ghq not available")
	}
}

func TestGetRootPath(t *testing.T) {
	cfg := &config.Config{
		Roots: map[string]string{
			"sandbox": "/tmp/sandbox",
			"dev":     "/tmp/dev",
		},
		Default: config.DefaultConfig{Root: "sandbox"},
	}
	c := NewClient(cfg)

	// Test that correct root path is retrieved
	for workspace, expectedPath := range cfg.Roots {
		path, exists := c.cfg.GetRoot(workspace)
		if !exists {
			t.Errorf("workspace %s should exist", workspace)
		}
		if path != expectedPath {
			t.Errorf("expected path %s for %s, got %s", expectedPath, workspace, path)
		}
	}
}
