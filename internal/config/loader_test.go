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
