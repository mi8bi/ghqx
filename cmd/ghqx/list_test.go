package main

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mi8bi/ghqx/internal/app"
	"github.com/mi8bi/ghqx/internal/config"
)

func TestRunListWithLoadAppError(t *testing.T) {
	oldConfigPath := configPath
	configPath = "/nonexistent/config.toml"
	defer func() { configPath = oldConfigPath }()

	oldApp := application
	application = nil
	defer func() { application = oldApp }()

	err := runList(listCmd, []string{})
	if err == nil {
		t.Fatal("expected error when loadApp fails")
	}
}

func TestRunListDefaultRoot(t *testing.T) {
	tmp, err := os.MkdirTemp("", "ghqx-list-default")
	if err != nil {
		t.Fatalf("tempdir: %v", err)
	}
	defer os.RemoveAll(tmp)

	// Create multiple roots
	devRoot := filepath.Join(tmp, "dev")
	sandboxRoot := filepath.Join(tmp, "sandbox")

	// Create repositories in different roots with .git directories
	devRepos := []string{"github.com/dev/repo1", "github.com/dev/repo2"}
	sandboxRepos := []string{"github.com/sandbox/repo3"}

	for _, repo := range devRepos {
		path := filepath.Join(devRoot, repo)
		if err := os.MkdirAll(filepath.Join(path, ".git"), 0755); err != nil {
			t.Fatalf("mkdir: %v", err)
		}
	}

	for _, repo := range sandboxRepos {
		path := filepath.Join(sandboxRoot, repo)
		if err := os.MkdirAll(filepath.Join(path, ".git"), 0755); err != nil {
			t.Fatalf("mkdir: %v", err)
		}
	}

	cfg := &config.Config{
		Roots: map[string]string{
			"dev":     devRoot,
			"sandbox": sandboxRoot,
		},
		Default: config.DefaultConfig{Root: "dev"}, // Default is dev
	}

	// Save old state and prevent loading real config
	oldApp := application
	oldConfigPath := configPath
	configPath = filepath.Join(tmp, "test-config.toml") // Point to non-existent file
	defer func() {
		application = oldApp
		configPath = oldConfigPath
	}()

	appInstance := app.New(cfg)
	application = appInstance

	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	listFullPath = false
	listAllRoots = false // Only default root
	err = runList(listCmd, []string{})

	w.Close()
	os.Stdout = oldStdout

	if err != nil {
		t.Fatalf("runList failed: %v", err)
	}

	var buf bytes.Buffer
	io.Copy(&buf, r)
	outputStr := buf.String()

	// Should contain dev repos
	for _, repo := range devRepos {
		normalizedOutput := filepath.ToSlash(outputStr)
		if !strings.Contains(normalizedOutput, repo) {
			t.Errorf("output should contain dev repo: %s\nGot output:\n%s", repo, outputStr)
		}
	}

	// Should NOT contain sandbox repos (not in default root)
	for _, repo := range sandboxRepos {
		normalizedOutput := filepath.ToSlash(outputStr)
		if strings.Contains(normalizedOutput, repo) {
			t.Errorf("output should NOT contain sandbox repo (not default root): %s\nGot output:\n%s", repo, outputStr)
		}
	}

	// Should NOT contain any paths outside test directory
	lines := strings.Split(strings.TrimSpace(outputStr), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		// Each line should be relative path matching our test repos
		// Normalize path separators for cross-platform testing
		normalizedLine := filepath.ToSlash(line)
		found := false
		for _, repo := range devRepos {
			if strings.Contains(normalizedLine, repo) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("unexpected output line (not from test repos): %s", line)
		}
	}
}

func TestRunListAllRoots(t *testing.T) {
	tmp, err := os.MkdirTemp("", "ghqx-list-all")
	if err != nil {
		t.Fatalf("tempdir: %v", err)
	}
	defer os.RemoveAll(tmp)

	// Create multiple roots
	devRoot := filepath.Join(tmp, "dev")
	sandboxRoot := filepath.Join(tmp, "sandbox")

	// Create repositories with .git directories
	devRepo := "github.com/dev/repo1"
	sandboxRepo := "github.com/sandbox/repo2"

	if err := os.MkdirAll(filepath.Join(devRoot, devRepo, ".git"), 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(sandboxRoot, sandboxRepo, ".git"), 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	cfg := &config.Config{
		Roots: map[string]string{
			"dev":     devRoot,
			"sandbox": sandboxRoot,
		},
		Default: config.DefaultConfig{Root: "dev"},
	}

	// Save old state
	oldApp := application
	oldConfigPath := configPath
	configPath = filepath.Join(tmp, "test-config.toml")
	defer func() {
		application = oldApp
		configPath = oldConfigPath
	}()

	appInstance := app.New(cfg)
	application = appInstance

	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	listFullPath = false
	listAllRoots = true // All roots
	err = runList(listCmd, []string{})

	w.Close()
	os.Stdout = oldStdout

	if err != nil {
		t.Fatalf("runList failed: %v", err)
	}

	var buf bytes.Buffer
	io.Copy(&buf, r)
	outputStr := buf.String()

	// Should contain both dev and sandbox repos
	normalizedOutput := filepath.ToSlash(outputStr)
	if !strings.Contains(normalizedOutput, devRepo) {
		t.Errorf("output should contain dev repo: %s\nGot output:\n%s", devRepo, outputStr)
	}
	if !strings.Contains(normalizedOutput, sandboxRepo) {
		t.Errorf("output should contain sandbox repo: %s\nGot output:\n%s", sandboxRepo, outputStr)
	}

	// Should only contain test repos
	lines := strings.Split(strings.TrimSpace(outputStr), "\n")
	expectedRepos := []string{devRepo, sandboxRepo}
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		// Normalize path separators for cross-platform testing
		normalizedLine := filepath.ToSlash(line)
		found := false
		for _, repo := range expectedRepos {
			if strings.Contains(normalizedLine, repo) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("unexpected output line (not from test repos): %s", line)
		}
	}
}

func TestRunListFullPath(t *testing.T) {
	tmp, err := os.MkdirTemp("", "ghqx-list-full")
	if err != nil {
		t.Fatalf("tempdir: %v", err)
	}
	defer os.RemoveAll(tmp)

	repo := "github.com/user/repo"
	path := filepath.Join(tmp, repo)
	if err := os.MkdirAll(filepath.Join(path, ".git"), 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	cfg := &config.Config{
		Roots:   map[string]string{"sandbox": tmp},
		Default: config.DefaultConfig{Root: "sandbox"},
	}

	// Save old state
	oldApp := application
	oldConfigPath := configPath
	configPath = filepath.Join(tmp, "test-config.toml")
	defer func() {
		application = oldApp
		configPath = oldConfigPath
	}()

	appInstance := app.New(cfg)
	application = appInstance

	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	listFullPath = true
	listAllRoots = false
	err = runList(listCmd, []string{})

	w.Close()
	os.Stdout = oldStdout

	if err != nil {
		t.Fatalf("runList failed: %v", err)
	}

	var buf bytes.Buffer
	io.Copy(&buf, r)
	outputStr := buf.String()

	// Should contain full absolute path
	expectedPath := filepath.Join(tmp, repo)
	if !strings.Contains(outputStr, expectedPath) {
		t.Errorf("output should contain full path: %s\nGot output:\n%s", expectedPath, outputStr)
	}

	// Should only contain paths within test directory
	lines := strings.Split(strings.TrimSpace(outputStr), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if !strings.HasPrefix(line, tmp) {
			t.Errorf("output line should be within test directory %s, got: %s", tmp, line)
		}
	}
}

func TestRunListRelativePath(t *testing.T) {
	tmp, err := os.MkdirTemp("", "ghqx-list-relative")
	if err != nil {
		t.Fatalf("tempdir: %v", err)
	}
	defer os.RemoveAll(tmp)

	repo := "github.com/user/repo"
	path := filepath.Join(tmp, repo)
	if err := os.MkdirAll(filepath.Join(path, ".git"), 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	cfg := &config.Config{
		Roots:   map[string]string{"sandbox": tmp},
		Default: config.DefaultConfig{Root: "sandbox"},
	}

	// Save old state
	oldApp := application
	oldConfigPath := configPath
	configPath = filepath.Join(tmp, "test-config.toml")
	defer func() {
		application = oldApp
		configPath = oldConfigPath
	}()

	appInstance := app.New(cfg)
	application = appInstance

	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	listFullPath = false
	listAllRoots = false
	err = runList(listCmd, []string{})

	w.Close()
	os.Stdout = oldStdout

	if err != nil {
		t.Fatalf("runList failed: %v", err)
	}

	var buf bytes.Buffer
	io.Copy(&buf, r)
	outputStr := buf.String()

	// Should contain relative path
	normalizedOutput := filepath.ToSlash(outputStr)
	if !strings.Contains(normalizedOutput, repo) {
		t.Errorf("output should contain relative path: %s\nGot output:\n%s", repo, outputStr)
	}

	// Should NOT contain absolute base path
	if strings.Contains(outputStr, tmp) {
		t.Errorf("output should not contain absolute base path\nGot output:\n%s", outputStr)
	}

	// Should only have lines with relative paths (no external repos)
	lines := strings.Split(strings.TrimSpace(outputStr), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		// Normalize path separators and check if it matches our test repo
		normalizedLine := filepath.ToSlash(line)
		if !strings.Contains(normalizedLine, repo) {
			t.Errorf("unexpected output line (not test repo): %s", line)
		}
	}
}

func TestRunListWithEmptyRoot(t *testing.T) {
	tmp, err := os.MkdirTemp("", "ghqx-list-empty")
	if err != nil {
		t.Fatalf("tempdir: %v", err)
	}
	defer os.RemoveAll(tmp)

	cfg := &config.Config{
		Roots:   map[string]string{"sandbox": tmp},
		Default: config.DefaultConfig{Root: "sandbox"},
	}

	// Save old state
	oldApp := application
	oldConfigPath := configPath
	configPath = filepath.Join(tmp, "test-config.toml")
	defer func() {
		application = oldApp
		configPath = oldConfigPath
	}()

	appInstance := app.New(cfg)
	application = appInstance

	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	listFullPath = false
	listAllRoots = false
	err = runList(listCmd, []string{})

	w.Close()
	os.Stdout = oldStdout

	if err != nil {
		t.Fatalf("runList should not fail with empty root: %v", err)
	}

	var buf bytes.Buffer
	io.Copy(&buf, r)
	outputStr := strings.TrimSpace(buf.String())

	// Empty output is expected
	if outputStr != "" {
		t.Errorf("expected empty output, got: %s", outputStr)
	}
}

func TestListCmdFlags(t *testing.T) {
	// Verify --full-path flag
	flag := listCmd.Flags().Lookup("full-path")
	if flag == nil {
		t.Error("--full-path flag should be registered")
	}
	if flag.Shorthand != "p" {
		t.Errorf("expected short flag 'p', got %q", flag.Shorthand)
	}

	// Verify --all flag
	allFlag := listCmd.Flags().Lookup("all")
	if allFlag == nil {
		t.Error("--all flag should be registered")
	}
	if allFlag.Shorthand != "a" {
		t.Errorf("expected short flag 'a', got %q", allFlag.Shorthand)
	}
}
