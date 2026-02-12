package main

import (
	"bufio"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mi8bi/ghqx/internal/config"
)

func TestRunConfigInitWithYesFlag(t *testing.T) {
	tmp, err := ioutil.TempDir("", "ghqx-config-init-yes")
	if err != nil {
		t.Fatalf("tempdir: %v", err)
	}
	defer os.RemoveAll(tmp)

	cfgPath := filepath.Join(tmp, "config.toml")

	oldConfigPath := configPath
	configPath = cfgPath
	defer func() { configPath = oldConfigPath }()

	oldInitYes := configInitYes
	configInitYes = true
	defer func() { configInitYes = oldInitYes }()

	err = runConfigInit(configInitCmd, []string{})
	if err != nil {
		t.Fatalf("runConfigInit with --yes failed: %v", err)
	}

	// Verify config file was created
	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		t.Fatalf("config file was not created")
	}

	// Verify content
	loader := config.NewLoader()
	cfg, err := loader.Load(cfgPath)
	if err != nil {
		t.Fatalf("failed to load created config: %v", err)
	}

	if len(cfg.Roots) == 0 {
		t.Fatalf("config roots are empty")
	}
}

func TestRunConfigInitFileExists(t *testing.T) {
	tmp, err := ioutil.TempDir("", "ghqx-config-exists")
	if err != nil {
		t.Fatalf("tempdir: %v", err)
	}
	defer os.RemoveAll(tmp)

	cfgPath := filepath.Join(tmp, "config.toml")

	// Create existing config file
	if err := ioutil.WriteFile(cfgPath, []byte("[roots]\ndev = \"/tmp/dev\"\n"), 0644); err != nil {
		t.Fatalf("failed to create existing config: %v", err)
	}

	oldConfigPath := configPath
	configPath = cfgPath
	defer func() { configPath = oldConfigPath }()

	err = runConfigInit(configInitCmd, []string{})
	if err == nil {
		t.Fatalf("expected error when config file already exists")
	}
}

func TestRunConfigShow(t *testing.T) {
	tmp, err := ioutil.TempDir("", "ghqx-config-show")
	if err != nil {
		t.Fatalf("tempdir: %v", err)
	}
	defer os.RemoveAll(tmp)

	cfgPath := filepath.Join(tmp, "config.toml")
	cfg := &config.Config{
		Roots:   map[string]string{"dev": filepath.Join(tmp, "dev")},
		Default: config.DefaultConfig{Root: "dev"},
	}

	loader := config.NewLoader()
	if err := loader.Save(cfg, cfgPath); err != nil {
		t.Fatalf("failed to save config: %v", err)
	}

	oldConfigPath := configPath
	configPath = cfgPath
	defer func() { configPath = oldConfigPath }()

	oldApp := application
	defer func() { application = oldApp }()

	err = runConfigShow(configShowCmd, []string{})
	if err != nil {
		t.Fatalf("runConfigShow failed: %v", err)
	}
}

func TestRunConfigShowWithoutApp(t *testing.T) {
	oldConfigPath := configPath
	configPath = "/nonexistent/config.toml"
	defer func() { configPath = oldConfigPath }()

	oldApp := application
	application = nil
	defer func() { application = oldApp }()

	err := runConfigShow(configShowCmd, []string{})
	if err == nil {
		t.Fatalf("expected error when loadApp fails")
	}
}

func TestRunConfigEdit(t *testing.T) {
	tmp, err := ioutil.TempDir("", "ghqx-config-edit")
	if err != nil {
		t.Fatalf("tempdir: %v", err)
	}
	defer os.RemoveAll(tmp)

	cfgPath := filepath.Join(tmp, "config.toml")
	cfg := &config.Config{
		Roots:   map[string]string{"dev": filepath.Join(tmp, "dev")},
		Default: config.DefaultConfig{Root: "dev"},
	}

	loader := config.NewLoader()
	if err := loader.Save(cfg, cfgPath); err != nil {
		t.Fatalf("failed to save config: %v", err)
	}

	oldConfigPath := configPath
	configPath = cfgPath
	defer func() { configPath = oldConfigPath }()

	// We can't fully test the TUI, but we can test that the command loads config
	// The actual TUI testing would require mocking tea.Program
	_, err = loader.Load(cfgPath)
	if err != nil {
		t.Fatalf("config edit preparation failed: %v", err)
	}
}

func TestPromptWithDefault(t *testing.T) {
	testCases := []struct {
		input        string
		defaultValue string
		expected     string
	}{
		{"", "default", "default"},
		{"custom", "default", "custom"},
		{"  custom  ", "default", "custom"},
	}

	for _, tc := range testCases {
		reader := bufio.NewReader(strings.NewReader(tc.input + "\n"))
		result := promptWithDefault(reader, "prompt", tc.defaultValue)
		if result != tc.expected {
			t.Errorf("promptWithDefault(%q, %q) = %q, want %q", tc.input, tc.defaultValue, result, tc.expected)
		}
	}
}

func TestPrintConfigSummary(t *testing.T) {
	cfg := &config.Config{
		Roots: map[string]string{
			"dev":     "/tmp/dev",
			"sandbox": "/tmp/sandbox",
		},
		Default: config.DefaultConfig{Root: "dev"},
	}

	// Capture output
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	printConfigSummary(cfg)

	w.Close()
	os.Stdout = oldStdout

	output, _ := ioutil.ReadAll(r)
	outputStr := string(output)

	if !strings.Contains(outputStr, "dev") {
		t.Errorf("printConfigSummary output missing 'dev'")
	}
	if !strings.Contains(outputStr, "sandbox") {
		t.Errorf("printConfigSummary output missing 'sandbox'")
	}
}

func TestPromptForConfig(t *testing.T) {
	input := "\n\n\n\n" // Press enter for all defaults

	// We need to temporarily redirect stdin
	oldStdin := os.Stdin
	r, w, _ := os.Pipe()
	os.Stdin = r

	go func() {
		w.Write([]byte(input))
		w.Close()
	}()
	defer func() { os.Stdin = oldStdin }()

	// Note: promptForConfig reads from bufio.NewReader(os.Stdin) internally
	// So we test the individual prompt function instead
	testReader := bufio.NewReader(strings.NewReader("\n"))
	result := promptWithDefault(testReader, "test", "default")
	if result != "default" {
		t.Errorf("expected 'default', got %q", result)
	}
}
