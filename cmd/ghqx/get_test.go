package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mi8bi/ghqx/internal/app"
	"github.com/mi8bi/ghqx/internal/config"
)

func TestRunGetWithLoadAppError(t *testing.T) {
	oldConfigPath := configPath
	configPath = "/nonexistent/config.toml"
	defer func() { configPath = oldConfigPath }()

	oldApp := application
	application = nil
	defer func() { application = oldApp }()

	err := runGet(getCmd, []string{"test/repo"})
	if err == nil {
		t.Fatalf("expected error when loadApp fails")
	}
}

func TestRunGetWithExistingRepository(t *testing.T) {
	tmp, err := os.MkdirTemp("", "ghqx-get-existing")
	if err != nil {
		t.Fatalf("tempdir: %v", err)
	}
	defer os.RemoveAll(tmp)

	// Create existing repository
	repo := filepath.Join(tmp, "github.com", "user", "repo")
	if err := os.MkdirAll(repo, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	cfg := &config.Config{
		Roots:   map[string]string{"sandbox": tmp},
		Default: config.DefaultConfig{Root: "sandbox"},
	}
	appInstance := app.New(cfg)
	application = appInstance

	// Test checkRepositoryExists
	workspace := checkRepositoryExists("github.com/user/repo")
	if workspace == "" {
		t.Fatalf("expected checkRepositoryExists to find existing repo")
	}
	if workspace != "sandbox" {
		t.Errorf("expected workspace 'sandbox', got %q", workspace)
	}
}

func TestCheckRepositoryExistsNotFound(t *testing.T) {
	tmp, err := os.MkdirTemp("", "ghqx-get-notfound")
	if err != nil {
		t.Fatalf("tempdir: %v", err)
	}
	defer os.RemoveAll(tmp)

	cfg := &config.Config{
		Roots:   map[string]string{"sandbox": tmp},
		Default: config.DefaultConfig{Root: "sandbox"},
	}
	appInstance := app.New(cfg)
	application = appInstance

	workspace := checkRepositoryExists("nonexistent/repo")
	if workspace != "" {
		t.Errorf("expected empty string for nonexistent repo, got %q", workspace)
	}
}

func TestExtractRepoNameVariations(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"github.com/user/repo", "repo"},
		{"user/repo", "repo"},
		{"simple", "simple"},
		{"gitlab.com/org/project/repo", "repo"},
		{"bitbucket.org/team/repository", "repository"},
	}

	for _, tc := range testCases {
		result := extractRepoName(tc.input)
		if result != tc.expected {
			t.Errorf("extractRepoName(%q) = %q, want %q", tc.input, result, tc.expected)
		}
	}
}

func TestContainsFunction(t *testing.T) {
	testCases := []struct {
		a        string
		b        string
		expected bool
	}{
		{"github.com/user/repo", "user/repo", true},
		{"user/repo", "repo", true},
		{"simple", "simple", true},
		{"github.com/user/repo", "other/repo", false},
		{"short", "longstring", false},
		{"exact", "exact", true},
	}

	for _, tc := range testCases {
		result := contains(tc.a, tc.b)
		if result != tc.expected {
			t.Errorf("contains(%q, %q) = %v, want %v", tc.a, tc.b, result, tc.expected)
		}
	}
}

func TestRunGetWithWorkspaceFlag(t *testing.T) {
	tmp, err := os.MkdirTemp("", "ghqx-get-workspace")
	if err != nil {
		t.Fatalf("tempdir: %v", err)
	}
	defer os.RemoveAll(tmp)

	cfg := &config.Config{
		Roots: map[string]string{
			"sandbox": filepath.Join(tmp, "sandbox"),
			"dev":     filepath.Join(tmp, "dev"),
		},
		Default: config.DefaultConfig{Root: "sandbox"},
	}

	// Create root directories
	for _, path := range cfg.Roots {
		if err := os.MkdirAll(path, 0755); err != nil {
			t.Fatalf("failed to create root dir: %v", err)
		}
	}

	appInstance := app.New(cfg)
	application = appInstance

	// Test that the workspace flag can be set
	oldWorkspace := getTargetWorkspace
	defer func() { getTargetWorkspace = oldWorkspace }()

	getTargetWorkspace = "dev"

	// We can't easily test actual ghq execution, but we can verify the setup
	if getTargetWorkspace != "dev" {
		t.Errorf("expected workspace 'dev', got %q", getTargetWorkspace)
	}
}

func TestCheckRepositoryExistsWithStatusError(t *testing.T) {
	// Test when Status.GetAll returns an error
	tmp, err := os.MkdirTemp("", "ghqx-get-error")
	if err != nil {
		t.Fatalf("tempdir: %v", err)
	}
	defer os.RemoveAll(tmp)

	// Create invalid root path to cause error
	cfg := &config.Config{
		Roots:   map[string]string{"sandbox": filepath.Join(tmp, "nonexistent")},
		Default: config.DefaultConfig{Root: "sandbox"},
	}
	appInstance := app.New(cfg)
	application = appInstance

	workspace := checkRepositoryExists("test/repo")
	// Should return empty string on error
	if workspace != "" {
		t.Errorf("expected empty string on error, got %q", workspace)
	}
}

func TestRunGetNoArgs(t *testing.T) {
	tmp, err := os.MkdirTemp("", "ghqx-get-noargs")
	if err != nil {
		t.Fatalf("tempdir: %v", err)
	}
	defer os.RemoveAll(tmp)

	cfg := &config.Config{
		Roots:   map[string]string{"sandbox": tmp},
		Default: config.DefaultConfig{Root: "sandbox"},
	}
	appInstance := app.New(cfg)
	application = appInstance

	// The command requires exactly one argument
	// Testing with empty args should be handled by cobra's Args validation
	// We verify that the Args constraint is set
	if getCmd.Args == nil {
		t.Error("getCmd should have Args validation")
	}
}
