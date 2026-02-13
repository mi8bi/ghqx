package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mi8bi/ghqx/internal/domain"
)

func TestValidateAndGetters(t *testing.T) {
	c := &Config{Roots: map[string]string{"dev": "/tmp/dev"}, Default: DefaultConfig{Root: "dev"}}
	if err := c.Validate(); err != nil {
		t.Fatalf("unexpected validate error: %v", err)
	}

	if p, ok := c.GetRoot("dev"); !ok || p != "/tmp/dev" {
		t.Fatalf("GetRoot failed: %v %v", p, ok)
	}

	if c.GetDefaultRoot() != "dev" {
		t.Fatalf("GetDefaultRoot mismatch")
	}
}

func TestValidateErrors(t *testing.T) {
	c := &Config{Roots: map[string]string{}}
	if err := c.Validate(); err == nil {
		t.Fatalf("expected error for no roots")
	}

	c = &Config{Roots: map[string]string{"a": "/tmp/a"}, Default: DefaultConfig{Root: "missing"}}
	if err := c.Validate(); err == nil {
		t.Fatalf("expected error for invalid default root")
	}
}

func TestNewDefaultConfigCreatesRoots(t *testing.T) {
	cfg := NewDefaultConfig()
	if len(cfg.Roots) == 0 {
		t.Fatalf("NewDefaultConfig should populate Roots")
	}
	// ensure GetDefaultRoot returns something
	if cfg.GetDefaultRoot() == "" {
		t.Fatalf("Default root should not be empty")
	}

	// Test EnsureRootDirectories creates directories
	tmpDir, err := os.MkdirTemp("", "ghqx-config-test")
	if err != nil {
		t.Fatalf("tempdir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	cfg2 := &Config{Roots: map[string]string{"r1": filepath.Join(tmpDir, "r1")}}
	if err := EnsureRootDirectories(cfg2); err != nil {
		t.Fatalf("EnsureRootDirectories failed: %v", err)
	}
	if _, err := os.Stat(cfg2.Roots["r1"]); os.IsNotExist(err) {
		t.Fatalf("expected directory created: %v", cfg2.Roots["r1"])
	}
}

func TestEnsureDirectoryErrors(t *testing.T) {
	// create a file where a directory is expected
	tmpDir, err := os.MkdirTemp("", "ghqx-ensure-test")
	if err != nil {
		t.Fatalf("tempdir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	filePath := filepath.Join(tmpDir, "file")
	if err := os.WriteFile(filePath, []byte("x"), 0644); err != nil {
		t.Fatalf("writefile: %v", err)
	}

	if err := ensureDirectory(filePath); err == nil {
		// ensureDirectory should return an error when path exists but is not a dir
		if _, ok := err.(*domain.GhqxError); !ok {
			t.Fatalf("expected domain.GhqxError, got %T", err)
		}
	}
}
