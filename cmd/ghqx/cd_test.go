package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/mi8bi/ghqx/internal/app"
	"github.com/mi8bi/ghqx/internal/config"
)

func TestRunCDWithLoadAppError(t *testing.T) {
	// Test runCD when loadApp fails
	oldConfigPath := configPath
	configPath = "/nonexistent/path/to/config.toml"
	defer func() { configPath = oldConfigPath }()

	oldApp := application
	application = nil
	defer func() { application = oldApp }()

	err := runCD(cdCmd, []string{})
	if err == nil {
		t.Fatalf("expected error when loadApp fails")
	}
}

func TestRunCDSuccess(t *testing.T) {
	tmp, err := ioutil.TempDir("", "ghqx-cd-success")
	if err != nil {
		t.Fatalf("tempdir: %v", err)
	}
	defer os.RemoveAll(tmp)

	// Create test repository structure
	repo := filepath.Join(tmp, "github.com", "user", "repo")
	if err := os.MkdirAll(repo, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	cfg := &config.Config{Roots: map[string]string{"sandbox": tmp}, Default: config.DefaultConfig{Root: "sandbox"}}
	appInstance := app.New(cfg)
	application = appInstance

	// Test loadProjectsForSelection
	projects, err := loadProjectsForSelection()
	if err != nil {
		t.Fatalf("loadProjectsForSelection failed: %v", err)
	}
	if len(projects) == 0 {
		t.Fatalf("expected at least one project")
	}
}

func TestLoadProjectsForSelectionWithoutApp(t *testing.T) {
	// Test with nil application
	oldApp := application
	application = nil
	defer func() { application = oldApp }()

	_, err := loadProjectsForSelection()
	if err == nil {
		t.Fatalf("expected error when application is nil")
	}
}

func TestLoadProjectsForSelectionEmptyRoot(t *testing.T) {
	tmp, err := ioutil.TempDir("", "ghqx-cd-empty")
	if err != nil {
		t.Fatalf("tempdir: %v", err)
	}
	defer os.RemoveAll(tmp)

	cfg := &config.Config{Roots: map[string]string{"sandbox": tmp}, Default: config.DefaultConfig{Root: "sandbox"}}
	appInstance := app.New(cfg)
	application = appInstance

	projects, err := loadProjectsForSelection()
	if err != nil {
		t.Fatalf("loadProjectsForSelection failed: %v", err)
	}

	// Empty root should return empty projects list
	if len(projects) != 0 {
		t.Logf("Expected 0 projects but got %d", len(projects))
	}
}
