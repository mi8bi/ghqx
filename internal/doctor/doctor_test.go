package doctor

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mi8bi/ghqx/internal/config"
	"github.com/mi8bi/ghqx/internal/i18n"
)

func TestRunChecksWhenToolsMissing(t *testing.T) {
	// Ensure locale messages available
	i18n.SetLocale(i18n.LocaleEN)

	// Empty PATH to simulate missing ghq/git
	origPath := os.Getenv("PATH")
	os.Setenv("PATH", "")
	defer os.Setenv("PATH", origPath)

	s := NewService()
	res := s.RunChecks()
	if len(res) != 3 {
		t.Fatalf("expected 3 checks, got %d", len(res))
	}

	// Config check may pass or fail depending on whether default config exists
	// We only check that ghq and git checks fail
	if res[1].OK {
		t.Fatalf("expected ghq check to fail when ghq missing")
	}
	if res[2].OK {
		t.Fatalf("expected git check to fail when git missing")
	}
}

// Additional tests for better coverage

func TestNewService(t *testing.T) {
	s := NewService()
	if s == nil {
		t.Fatal("NewService returned nil")
	}
	if s.configLoader == nil {
		t.Fatal("configLoader should not be nil")
	}
	if s.configPath != "" {
		t.Error("configPath should be empty for NewService")
	}
}

func TestNewServiceWithConfigPath(t *testing.T) {
	testPath := "/test/path/config.toml"
	s := NewServiceWithConfigPath(testPath)
	if s == nil {
		t.Fatal("NewServiceWithConfigPath returned nil")
	}
	if s.configPath != testPath {
		t.Errorf("expected configPath %s, got %s", testPath, s.configPath)
	}
}

func TestCheckConfigWithValidConfig(t *testing.T) {
	tmp, err := os.MkdirTemp("", "ghqx-doctor-valid")
	if err != nil {
		t.Fatalf("tempdir: %v", err)
	}
	defer os.RemoveAll(tmp)

	cfgPath := filepath.Join(tmp, "config.toml")
	cfg := &config.Config{
		Roots:   map[string]string{"dev": filepath.Join(tmp, "dev")},
		Default: config.DefaultConfig{Root: "dev"},
	}

	if err := os.MkdirAll(cfg.Roots["dev"], 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	loader := config.NewLoader()
	if err := loader.Save(cfg, cfgPath); err != nil {
		t.Fatalf("save config: %v", err)
	}

	s := NewServiceWithConfigPath(cfgPath)
	result := s.CheckConfig()

	if !result.OK {
		t.Errorf("CheckConfig should pass with valid config: %s", result.Message)
	}
	if result.Name == "" {
		t.Error("result Name should not be empty")
	}
}

func TestCheckConfigWithInvalidConfig(t *testing.T) {
	s := NewServiceWithConfigPath("/nonexistent/config.toml")
	result := s.CheckConfig()

	if result.OK {
		t.Error("CheckConfig should fail with nonexistent config")
	}
	if result.Hint == "" {
		t.Error("result Hint should not be empty on failure")
	}
}

func TestCheckGhqWhenAvailable(t *testing.T) {
	// Skip if ghq is not available
	if _, err := exec.LookPath("ghq"); err != nil {
		t.Skip("ghq not available, skipping test")
	}

	s := NewService()
	result := s.CheckGhq()

	if !result.OK {
		t.Errorf("CheckGhq should pass when ghq is available: %s", result.Message)
	}
	if !strings.Contains(result.Message, "ghq found") && !strings.Contains(result.Message, "が見つかりました") {
		t.Error("message should indicate ghq was found")
	}
}

func TestCheckGhqWhenNotAvailable(t *testing.T) {
	origPath := os.Getenv("PATH")
	os.Setenv("PATH", "")
	defer os.Setenv("PATH", origPath)

	s := NewService()
	result := s.CheckGhq()

	if result.OK {
		t.Error("CheckGhq should fail when ghq not in PATH")
	}
	if result.Hint == "" {
		t.Error("should provide hint when ghq not found")
	}
}

func TestCheckGhqExecutionFailure(t *testing.T) {
	// Create a fake "ghq" that fails
	tmp, err := os.MkdirTemp("", "ghqx-fake-ghq")
	if err != nil {
		t.Fatalf("tempdir: %v", err)
	}
	defer os.RemoveAll(tmp)

	fakeGhq := filepath.Join(tmp, "ghq")
	var content string
	if os.PathSeparator == '\\' {
		// Windows batch file
		content = "@echo off\nexit /b 1\n"
		fakeGhq += ".bat"
	} else {
		// Unix shell script
		content = "#!/bin/sh\nexit 1\n"
	}

	if err := os.WriteFile(fakeGhq, []byte(content), 0755); err != nil {
		t.Fatalf("write fake ghq: %v", err)
	}

	origPath := os.Getenv("PATH")
	os.Setenv("PATH", tmp)
	defer os.Setenv("PATH", origPath)

	s := NewService()
	result := s.CheckGhq()

	// Should find ghq but fail to execute
	if result.OK {
		t.Error("CheckGhq should fail when ghq execution fails")
	}
}

func TestCheckGitWhenAvailable(t *testing.T) {
	// Skip if git is not available
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available, skipping test")
	}

	s := NewService()
	result := s.CheckGit()

	if !result.OK {
		t.Errorf("CheckGit should pass when git is available: %s", result.Message)
	}
	if !strings.Contains(result.Message, "git found") && !strings.Contains(result.Message, "が見つかりました") {
		t.Error("message should indicate git was found")
	}
}

func TestCheckGitWhenNotAvailable(t *testing.T) {
	origPath := os.Getenv("PATH")
	os.Setenv("PATH", "")
	defer os.Setenv("PATH", origPath)

	s := NewService()
	result := s.CheckGit()

	if result.OK {
		t.Error("CheckGit should fail when git not in PATH")
	}
	if result.Hint == "" {
		t.Error("should provide hint when git not found")
	}
}

func TestCheckGitExecutionFailure(t *testing.T) {
	// Create a fake "git" that fails
	tmp, err := os.MkdirTemp("", "ghqx-fake-git")
	if err != nil {
		t.Fatalf("tempdir: %v", err)
	}
	defer os.RemoveAll(tmp)

	fakeGit := filepath.Join(tmp, "git")
	var content string
	if os.PathSeparator == '\\' {
		// Windows batch file
		content = "@echo off\nexit /b 1\n"
		fakeGit += ".bat"
	} else {
		// Unix shell script
		content = "#!/bin/sh\nexit 1\n"
	}

	if err := os.WriteFile(fakeGit, []byte(content), 0755); err != nil {
		t.Fatalf("write fake git: %v", err)
	}

	origPath := os.Getenv("PATH")
	os.Setenv("PATH", tmp)
	defer os.Setenv("PATH", origPath)

	s := NewService()
	result := s.CheckGit()

	// Should find git but fail to execute
	if result.OK {
		t.Error("CheckGit should fail when git execution fails")
	}
}

func TestRunChecksWithAllAvailable(t *testing.T) {
	// Skip if tools are not available
	if _, err := exec.LookPath("ghq"); err != nil {
		t.Skip("ghq not available")
	}
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available")
	}

	tmp, err := os.MkdirTemp("", "ghqx-doctor-all")
	if err != nil {
		t.Fatalf("tempdir: %v", err)
	}
	defer os.RemoveAll(tmp)

	cfgPath := filepath.Join(tmp, "config.toml")
	cfg := &config.Config{
		Roots:   map[string]string{"dev": filepath.Join(tmp, "dev")},
		Default: config.DefaultConfig{Root: "dev"},
	}

	if err := os.MkdirAll(cfg.Roots["dev"], 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	loader := config.NewLoader()
	if err := loader.Save(cfg, cfgPath); err != nil {
		t.Fatalf("save config: %v", err)
	}

	s := NewServiceWithConfigPath(cfgPath)
	results := s.RunChecks()

	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}

	// All checks should pass
	for _, result := range results {
		if !result.OK {
			t.Errorf("check %s failed: %s", result.Name, result.Message)
		}
	}
}

func TestCheckResultStructure(t *testing.T) {
	result := CheckResult{
		Name:    "test",
		OK:      true,
		Message: "test message",
		Hint:    "test hint",
	}

	if result.Name != "test" {
		t.Error("Name mismatch")
	}
	if !result.OK {
		t.Error("OK mismatch")
	}
	if result.Message != "test message" {
		t.Error("Message mismatch")
	}
	if result.Hint != "test hint" {
		t.Error("Hint mismatch")
	}
}
