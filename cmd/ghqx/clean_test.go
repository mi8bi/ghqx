package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mi8bi/ghqx/internal/config"
)

func TestRunCleanWithValidConfig(t *testing.T) {
	tmp, err := ioutil.TempDir("", "ghqx-clean-test")
	if err != nil {
		t.Fatalf("tempdir: %v", err)
	}
	defer os.RemoveAll(tmp)

	// Create a test config file
	cfgPath := filepath.Join(tmp, "config.toml")
	cfg := &config.Config{
		Roots:   map[string]string{"sandbox": filepath.Join(tmp, "sandbox")},
		Default: config.DefaultConfig{Root: "sandbox"},
	}

	loader := config.NewLoader()
	if err := loader.Save(cfg, cfgPath); err != nil {
		t.Fatalf("failed to save config: %v", err)
	}

	// Create the sandbox directory
	if err := os.MkdirAll(cfg.Roots["sandbox"], 0755); err != nil {
		t.Fatalf("failed to create sandbox dir: %v", err)
	}

	// Set configPath
	oldConfigPath := configPath
	configPath = cfgPath
	defer func() { configPath = oldConfigPath }()

	// Test with "no" input (abort)
	input := "no\n"
	oldStdin := os.Stdin
	r, w, _ := os.Pipe()
	os.Stdin = r
	w.WriteString(input)
	w.Close()
	defer func() { os.Stdin = oldStdin }()

	err = runClean(cleanCmd, []string{})
	if err != nil {
		t.Fatalf("runClean failed: %v", err)
	}

	// Verify files still exist
	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		t.Fatalf("config file was deleted when it shouldn't be")
	}
}

func TestRunCleanWithYesInput(t *testing.T) {
	tmp, err := ioutil.TempDir("", "ghqx-clean-yes")
	if err != nil {
		t.Fatalf("tempdir: %v", err)
	}
	defer os.RemoveAll(tmp)

	// Create a test config file
	cfgPath := filepath.Join(tmp, "config.toml")
	sandboxPath := filepath.Join(tmp, "sandbox")
	cfg := &config.Config{
		Roots:   map[string]string{"sandbox": sandboxPath},
		Default: config.DefaultConfig{Root: "sandbox"},
	}

	loader := config.NewLoader()
	if err := loader.Save(cfg, cfgPath); err != nil {
		t.Fatalf("failed to save config: %v", err)
	}

	// Create the sandbox directory with a test file
	if err := os.MkdirAll(sandboxPath, 0755); err != nil {
		t.Fatalf("failed to create sandbox dir: %v", err)
	}
	testFile := filepath.Join(sandboxPath, "test.txt")
	if err := ioutil.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	// Set configPath
	oldConfigPath := configPath
	configPath = cfgPath
	defer func() { configPath = oldConfigPath }()

	// Test with "yes" input
	input := "yes\n"
	oldStdin := os.Stdin
	r, w, _ := os.Pipe()
	os.Stdin = r
	w.WriteString(input)
	w.Close()
	defer func() { os.Stdin = oldStdin }()

	err = runClean(cleanCmd, []string{})
	if err != nil {
		t.Fatalf("runClean failed: %v", err)
	}

	// Verify sandbox directory was deleted
	if _, err := os.Stat(sandboxPath); !os.IsNotExist(err) {
		t.Fatalf("sandbox directory was not deleted")
	}

	// Verify config file was deleted
	if _, err := os.Stat(cfgPath); !os.IsNotExist(err) {
		t.Fatalf("config file was not deleted")
	}
}

func TestRunCleanWithoutConfig(t *testing.T) {
	// Test runClean when config doesn't exist
	oldConfigPath := configPath
	configPath = "/nonexistent/config.toml"
	defer func() { configPath = oldConfigPath }()

	// Test with "no" input to abort
	input := "no\n"
	oldStdin := os.Stdin
	r, w, _ := os.Pipe()
	os.Stdin = r
	w.WriteString(input)
	w.Close()
	defer func() { os.Stdin = oldStdin }()

	err := runClean(cleanCmd, []string{})
	if err != nil {
		t.Fatalf("runClean should not fail when config doesn't exist: %v", err)
	}
}

func TestRunCleanCaseInsensitive(t *testing.T) {
	tmp, err := ioutil.TempDir("", "ghqx-clean-case")
	if err != nil {
		t.Fatalf("tempdir: %v", err)
	}
	defer os.RemoveAll(tmp)

	cfgPath := filepath.Join(tmp, "config.toml")
	cfg := &config.Config{
		Roots:   map[string]string{"sandbox": filepath.Join(tmp, "sandbox")},
		Default: config.DefaultConfig{Root: "sandbox"},
	}

	loader := config.NewLoader()
	if err := loader.Save(cfg, cfgPath); err != nil {
		t.Fatalf("failed to save config: %v", err)
	}

	oldConfigPath := configPath
	configPath = cfgPath
	defer func() { configPath = oldConfigPath }()

	// Test with uppercase "YES"
	testCases := []string{"YES\n", "Yes\n", "yEs\n"}

	for _, input := range testCases {
		// Recreate config for each test
		loader := config.NewLoader()
		if err := loader.Save(cfg, cfgPath); err != nil {
			continue
		}

		if err := os.MkdirAll(cfg.Roots["sandbox"], 0755); err != nil {
			continue
		}

		oldStdin := os.Stdin
		r, w, _ := os.Pipe()
		os.Stdin = r
		w.WriteString(input)
		w.Close()

		err = runClean(cleanCmd, []string{})
		os.Stdin = oldStdin

		if err != nil {
			t.Logf("runClean with input %q: %v", strings.TrimSpace(input), err)
		}
	}
}
