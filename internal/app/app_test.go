package app

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/mi8bi/ghqx/internal/config"
)

func TestNewAndNewFromConfigPath(t *testing.T) {
	cfg := config.NewDefaultConfig()
	a := New(cfg)
	if a == nil {
		t.Fatalf("New returned nil")
	}
	if a.Status == nil || a.Config == nil {
		t.Fatalf("App fields not initialized")
	}

	// create a temp config file and save using loader
	tmp, err := ioutil.TempDir("", "ghqx-app-test")
	if err != nil {
		t.Fatalf("tempdir: %v", err)
	}
	defer os.RemoveAll(tmp)

	path := filepath.Join(tmp, "cfg.toml")
	loader := config.NewLoader()
	if err := loader.Save(cfg, path); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	a2, err := NewFromConfigPath(path)
	if err != nil {
		t.Fatalf("NewFromConfigPath failed: %v", err)
	}
	if a2 == nil || a2.Config == nil {
		t.Fatalf("NewFromConfigPath returned invalid App")
	}
}
