package main

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mi8bi/ghqx/internal/config"
	"github.com/mi8bi/ghqx/internal/doctor"
	"github.com/mi8bi/ghqx/internal/i18n"
)

func TestRunDoctorWithInvalidConfig(t *testing.T) {
	// Test the doctor service with invalid config path
	doctorService := doctor.NewServiceWithConfigPath("/nonexistent/config.toml")
	results := doctorService.RunChecks()

	// Config check should fail
	foundFailedConfigCheck := false
	for _, res := range results {
		if strings.Contains(res.Name, "config") && !res.OK {
			foundFailedConfigCheck = true
			break
		}
	}

	if !foundFailedConfigCheck {
		t.Error("expected config check to fail with nonexistent config")
	}
}

func TestDoctorServiceChecks(t *testing.T) {
	// This test focuses on the doctor service itself, not the command execution
	tmp, err := os.MkdirTemp("", "ghqx-doctor-service")
	if err != nil {
		t.Fatalf("tempdir: %v", err)
	}
	defer os.RemoveAll(tmp)

	// Create a valid config file
	cfgPath := filepath.Join(tmp, "config.toml")
	cfg := &config.Config{
		Roots:   map[string]string{"dev": filepath.Join(tmp, "dev")},
		Default: config.DefaultConfig{Root: "dev"},
	}

	// Create the dev directory
	if err := os.MkdirAll(cfg.Roots["dev"], 0755); err != nil {
		t.Fatalf("failed to create dev dir: %v", err)
	}

	loader := config.NewLoader()
	if err := loader.Save(cfg, cfgPath); err != nil {
		t.Fatalf("failed to save config: %v", err)
	}

	// Test the doctor service directly with explicit config path
	doctorService := doctor.NewServiceWithConfigPath(cfgPath)
	results := doctorService.RunChecks()

	if len(results) == 0 {
		t.Fatal("expected at least one check result")
	}

	// Verify that config check exists and passes
	foundConfigCheck := false
	for _, res := range results {
		if strings.Contains(res.Name, "config") {
			foundConfigCheck = true
			if !res.OK {
				t.Errorf("config check failed: %s", res.Message)
			}
		}
	}

	if !foundConfigCheck {
		t.Error("config check not found in results")
	}
}

func TestDoctorOutputFormat(t *testing.T) {
	// Test the output format without actually running the command
	tmp, err := os.MkdirTemp("", "ghqx-doctor-output")
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

	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	doctorService := doctor.NewServiceWithConfigPath(cfgPath)
	results := doctorService.RunChecks()

	// Simulate the output logic from runDoctor
	for _, res := range results {
		if res.OK {
			os.Stdout.WriteString(i18n.T("doctor.result.ok") + " " + res.Message + "\n")
		} else {
			os.Stdout.WriteString(i18n.T("doctor.result.ng") + " " + res.Message + "\n")
			if res.Hint != "" {
				os.Stdout.WriteString("     " + i18n.T("doctor.result.hint") + ": " + res.Hint + "\n")
			}
		}
	}

	w.Close()
	os.Stdout = oldStdout
	output, _ := io.ReadAll(r)
	outputStr := string(output)

	// Verify output format
	if !strings.Contains(outputStr, i18n.T("doctor.result.ok")) && !strings.Contains(outputStr, i18n.T("doctor.result.ng")) {
		t.Error("output should contain OK or NG markers")
	}
}

func TestRunDoctorWithEmptyPath(t *testing.T) {
	// Save original PATH
	oldPath := os.Getenv("PATH")
	defer os.Setenv("PATH", oldPath)

	// Set PATH to empty
	os.Setenv("PATH", "")

	// Create a valid config to pass config check
	tmp, err := os.MkdirTemp("", "ghqx-doctor-nopath")
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
		t.Fatalf("failed to create dev dir: %v", err)
	}

	loader := config.NewLoader()
	if err := loader.Save(cfg, cfgPath); err != nil {
		t.Fatalf("failed to save config: %v", err)
	}

	// Test doctor service
	doctorService := doctor.NewServiceWithConfigPath(cfgPath)
	results := doctorService.RunChecks()

	// ghq and git checks should fail
	ghqFailed := false
	gitFailed := false
	for _, res := range results {
		if strings.Contains(res.Name, "ghq") && !res.OK {
			ghqFailed = true
		}
		if strings.Contains(res.Name, "git") && !res.OK {
			gitFailed = true
		}
	}

	if !ghqFailed {
		t.Error("expected ghq check to fail when PATH is empty")
	}
	if !gitFailed {
		t.Error("expected git check to fail when PATH is empty")
	}
}

// Test runDoctor command directly
func TestRunDoctorCommand(t *testing.T) {
	tmp, err := os.MkdirTemp("", "ghqx-doctor-cmd")
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
		t.Fatalf("failed to create dev dir: %v", err)
	}

	loader := config.NewLoader()
	if err := loader.Save(cfg, cfgPath); err != nil {
		t.Fatalf("failed to save config: %v", err)
	}

	oldConfigPath := configPath
	configPath = cfgPath
	defer func() { configPath = oldConfigPath }()

	// Capture output
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err = runDoctor(doctorCmd, []string{})

	w.Close()
	os.Stdout = oldStdout
	output, _ := io.ReadAll(r)
	outputStr := string(output)

	// Should contain check results
	if !strings.Contains(outputStr, i18n.T("doctor.check.config.ok")) {
		t.Error("output should contain config check")
	}
}
