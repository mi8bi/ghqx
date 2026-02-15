package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetDefaultConfigPathWithXDG(t *testing.T) {
	tmp, err := os.MkdirTemp("", "ghqx-xdg")
	if err != nil {
		t.Fatalf("tempdir: %v", err)
	}
	defer os.RemoveAll(tmp)

	os.Setenv("XDG_CONFIG_HOME", tmp)
	defer os.Unsetenv("XDG_CONFIG_HOME")

	got, err := GetDefaultConfigPath()
	if err != nil {
		t.Fatalf("GetDefaultConfigPath error: %v", err)
	}
	want := filepath.Join(tmp, "ghqx", "config.toml")
	if got != want {
		t.Fatalf("unexpected path: got %s want %s", got, want)
	}
}

func TestLoaderSaveAndLoad(t *testing.T) {
	tmp, err := os.MkdirTemp("", "ghqx-loader")
	if err != nil {
		t.Fatalf("tempdir: %v", err)
	}
	defer os.RemoveAll(tmp)

	cfg := &Config{Roots: map[string]string{"r": "/tmp/r"}, Default: DefaultConfig{Root: "r"}}
	loader := NewLoader()

	path := filepath.Join(tmp, "cfg.toml")
	if err := loader.Save(cfg, path); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Load via explicit path
	loaded, err := loader.Load(path)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if _, ok := loaded.GetRoot("r"); !ok {
		t.Fatalf("loaded config missing root")
	}
}

func TestFindConfigPathEnvVar(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "ghqx-config-*.toml")
	if err != nil {
		t.Fatalf("tempfile: %v", err)
	}
	tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	os.Setenv("GHQX_CONFIG", tmpFile.Name())
	defer os.Unsetenv("GHQX_CONFIG")

	loader := NewLoader()
	p, err := loader.findConfigPath("")
	if err != nil {
		t.Fatalf("findConfigPath failed: %v", err)
	}
	if p != tmpFile.Name() {
		t.Fatalf("unexpected path: got %s want %s", p, tmpFile.Name())
	}
}

// Additional tests for better coverage

func TestGetDefaultConfigPathWithoutXDG(t *testing.T) {
	// Unset XDG_CONFIG_HOME to test fallback
	oldXDG := os.Getenv("XDG_CONFIG_HOME")
	os.Unsetenv("XDG_CONFIG_HOME")
	defer func() {
		if oldXDG != "" {
			os.Setenv("XDG_CONFIG_HOME", oldXDG)
		}
	}()

	path, err := GetDefaultConfigPath()
	if err != nil {
		t.Fatalf("GetDefaultConfigPath error: %v", err)
	}
	
	if path == "" {
		t.Fatal("path should not be empty")
	}

	// Should contain .config/ghqx/config.toml
	if !filepath.IsAbs(path) {
		t.Error("path should be absolute")
	}
}

func TestLoaderSaveWithEmptyPath(t *testing.T) {
	tmp, err := os.MkdirTemp("", "ghqx-loader-empty")
	if err != nil {
		t.Fatalf("tempdir: %v", err)
	}
	defer os.RemoveAll(tmp)

	// Change HOME to temp dir for this test
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tmp)
	defer os.Setenv("HOME", oldHome)

	oldUserProfile := os.Getenv("USERPROFILE")
	os.Setenv("USERPROFILE", tmp)
	defer os.Setenv("USERPROFILE", oldUserProfile)

	cfg := &Config{
		Roots:   map[string]string{"dev": "/tmp/dev"},
		Default: DefaultConfig{Root: "dev"},
	}

	loader := NewLoader()
	
	// Save with empty path (should use default)
	err = loader.Save(cfg, "")
	if err != nil {
		t.Fatalf("Save with empty path failed: %v", err)
	}

	// Verify file was created at default location
	defaultPath, _ := GetDefaultConfigPath()
	if _, err := os.Stat(defaultPath); os.IsNotExist(err) {
		t.Error("config file should be created at default location")
	}
}

func TestFindConfigPathWithHomeConfig(t *testing.T) {
	tmp, err := os.MkdirTemp("", "ghqx-home-config")
	if err != nil {
		t.Fatalf("tempdir: %v", err)
	}
	defer os.RemoveAll(tmp)

	// Create ~/.ghqx.toml
	homeConfig := filepath.Join(tmp, ".ghqx.toml")
	if err := os.WriteFile(homeConfig, []byte("[roots]\ndev=\"/tmp/dev\"\n"), 0644); err != nil {
		t.Fatalf("write home config: %v", err)
	}

	// Set HOME
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tmp)
	defer os.Setenv("HOME", oldHome)

	oldUserProfile := os.Getenv("USERPROFILE")
	os.Setenv("USERPROFILE", tmp)
	defer os.Setenv("USERPROFILE", oldUserProfile)

	// Unset other env vars
	oldXDG := os.Getenv("XDG_CONFIG_HOME")
	os.Unsetenv("XDG_CONFIG_HOME")
	defer func() {
		if oldXDG != "" {
			os.Setenv("XDG_CONFIG_HOME", oldXDG)
		}
	}()

	oldGHQX := os.Getenv("GHQX_CONFIG")
	os.Unsetenv("GHQX_CONFIG")
	defer func() {
		if oldGHQX != "" {
			os.Setenv("GHQX_CONFIG", oldGHQX)
		}
	}()

	loader := NewLoader()
	path, err := loader.findConfigPath("")
	if err != nil {
		t.Fatalf("findConfigPath failed: %v", err)
	}

	if path != homeConfig {
		t.Errorf("expected %s, got %s", homeConfig, path)
	}
}

func TestFindConfigPathPriority(t *testing.T) {
	tmp, err := os.MkdirTemp("", "ghqx-priority")
	if err != nil {
		t.Fatalf("tempdir: %v", err)
	}
	defer os.RemoveAll(tmp)

	// Create config files in different locations
	xdgConfig := filepath.Join(tmp, "xdg", "ghqx", "config.toml")
	if err := os.MkdirAll(filepath.Dir(xdgConfig), 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(xdgConfig, []byte("[roots]\n"), 0644); err != nil {
		t.Fatalf("write xdg config: %v", err)
	}

	homeConfig := filepath.Join(tmp, ".config", "ghqx", "config.toml")
	if err := os.MkdirAll(filepath.Dir(homeConfig), 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(homeConfig, []byte("[roots]\n"), 0644); err != nil {
		t.Fatalf("write home config: %v", err)
	}

	// Test XDG_CONFIG_HOME priority
	os.Setenv("XDG_CONFIG_HOME", filepath.Join(tmp, "xdg"))
	defer os.Unsetenv("XDG_CONFIG_HOME")

	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tmp)
	defer os.Setenv("HOME", oldHome)

	loader := NewLoader()
	path, err := loader.findConfigPath("")
	if err != nil {
		t.Fatalf("findConfigPath failed: %v", err)
	}

	// XDG should have priority
	if path != xdgConfig {
		t.Errorf("expected XDG config, got %s", path)
	}
}

func TestLoadFromPathWithInvalidToml(t *testing.T) {
	tmp, err := os.MkdirTemp("", "ghqx-invalid-toml")
	if err != nil {
		t.Fatalf("tempdir: %v", err)
	}
	defer os.RemoveAll(tmp)

	invalidPath := filepath.Join(tmp, "invalid.toml")
	if err := os.WriteFile(invalidPath, []byte("invalid [[[[ toml"), 0644); err != nil {
		t.Fatalf("write invalid toml: %v", err)
	}

	loader := NewLoader()
	_, err = loader.loadFromPath(invalidPath)
	if err == nil {
		t.Fatal("expected error for invalid TOML")
	}
}

func TestLoadFromPathWithInvalidConfig(t *testing.T) {
	tmp, err := os.MkdirTemp("", "ghqx-invalid-config")
	if err != nil {
		t.Fatalf("tempdir: %v", err)
	}
	defer os.RemoveAll(tmp)

	invalidPath := filepath.Join(tmp, "invalid.toml")
	// Valid TOML but invalid config (no roots)
	if err := os.WriteFile(invalidPath, []byte("[default]\nroot=\"dev\"\n"), 0644); err != nil {
		t.Fatalf("write invalid config: %v", err)
	}

	loader := NewLoader()
	_, err = loader.loadFromPath(invalidPath)
	if err == nil {
		t.Fatal("expected validation error for invalid config")
	}
}

func TestSaveWithInvalidConfig(t *testing.T) {
	tmp, err := os.MkdirTemp("", "ghqx-save-invalid")
	if err != nil {
		t.Fatalf("tempdir: %v", err)
	}
	defer os.RemoveAll(tmp)

	invalidCfg := &Config{
		Roots:   map[string]string{}, // Invalid: no roots
		Default: DefaultConfig{Root: "dev"},
	}

	loader := NewLoader()
	path := filepath.Join(tmp, "config.toml")
	
	err = loader.Save(invalidCfg, path)
	if err == nil {
		t.Fatal("expected validation error when saving invalid config")
	}
}

func TestSaveCreateDirError(t *testing.T) {
	// Try to save to a path where parent is a file (not a directory)
	tmp, err := os.MkdirTemp("", "ghqx-save-dir-error")
	if err != nil {
		t.Fatalf("tempdir: %v", err)
	}
	defer os.RemoveAll(tmp)

	// Create a file
	blockingFile := filepath.Join(tmp, "blocking")
	if err := os.WriteFile(blockingFile, []byte("test"), 0644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	cfg := &Config{
		Roots:   map[string]string{"dev": "/tmp/dev"},
		Default: DefaultConfig{Root: "dev"},
	}

	loader := NewLoader()
	// Try to save to blocking/config.toml (can't create dir)
	path := filepath.Join(blockingFile, "config.toml")
	
	err = loader.Save(cfg, path)
	if err == nil {
		t.Fatal("expected error when can't create directory")
	}
}

func TestLoadWithNoConfigFound(t *testing.T) {
	// Clear all environment variables
	oldXDG := os.Getenv("XDG_CONFIG_HOME")
	os.Unsetenv("XDG_CONFIG_HOME")
	defer func() {
		if oldXDG != "" {
			os.Setenv("XDG_CONFIG_HOME", oldXDG)
		}
	}()

	oldGHQX := os.Getenv("GHQX_CONFIG")
	os.Unsetenv("GHQX_CONFIG")
	defer func() {
		if oldGHQX != "" {
			os.Setenv("GHQX_CONFIG", oldGHQX)
		}
	}()

	oldHome := os.Getenv("HOME")
	// Set HOME to non-existent directory
	os.Setenv("HOME", "/nonexistent")
	defer os.Setenv("HOME", oldHome)

	oldUserProfile := os.Getenv("USERPROFILE")
	os.Setenv("USERPROFILE", "/nonexistent")
	defer os.Setenv("USERPROFILE", oldUserProfile)

	loader := NewLoader()
	_, err := loader.Load("")
	if err == nil {
		t.Fatal("expected error when no config found")
	}
}