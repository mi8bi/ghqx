package ghq

import (
	"os"
	"path/filepath"
	"testing"

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

func TestGetWithInvalidWorkspace(t *testing.T) {
	cfg := &config.Config{
		Roots:   map[string]string{"sandbox": "/tmp/sandbox"},
		Default: config.DefaultConfig{Root: "sandbox"},
	}
	c := NewClient(cfg)

	err := c.Get(GetOptions{Repository: "test/repo", Workspace: "nonexistent"})
	if err == nil {
		t.Fatalf("expected error when workspace doesn't exist")
	}
}

func TestNewClient(t *testing.T) {
	cfg := &config.Config{
		Roots:   map[string]string{"sandbox": "/tmp"},
		Default: config.DefaultConfig{Root: "sandbox"},
	}
	c := NewClient(cfg)

	if c == nil {
		t.Fatal("NewClient returned nil")
	}

	if c.cfg == nil {
		t.Fatal("Client config is nil")
	}

	if c.timeout == 0 {
		t.Fatal("Client timeout is zero")
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
