package main

import (
	"io"

	"os"
	"runtime"
	"strings"
	"testing"
)

func TestRunVersionDefault(t *testing.T) {
	// Save original values
	origVersion := version
	origCommit := commit
	origBuildTime := buildTime
	defer func() {
		version = origVersion
		commit = origCommit
		buildTime = origBuildTime
	}()

	// Set test values
	version = "v1.0.0"
	commit = "abc123"
	buildTime = "2026-02-11T00:00:00Z"

	oldVerbose := versionVerbose
	versionVerbose = false
	defer func() { versionVerbose = oldVerbose }()

	// Capture output
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	runVersion(versionCmd, []string{})

	w.Close()
	os.Stdout = oldStdout
	output, _ := io.ReadAll(r)
	outputStr := string(output)

	// Verify output contains version
	if !strings.Contains(outputStr, "v1.0.0") {
		t.Errorf("version output should contain version number, got: %s", outputStr)
	}

	// Should not contain verbose info
	if strings.Contains(outputStr, "commit") {
		t.Errorf("default output should not contain commit info")
	}
	if strings.Contains(outputStr, "built at") {
		t.Errorf("default output should not contain build time")
	}
}

func TestRunVersionVerbose(t *testing.T) {
	// Save original values
	origVersion := version
	origCommit := commit
	origBuildTime := buildTime
	defer func() {
		version = origVersion
		commit = origCommit
		buildTime = origBuildTime
	}()

	// Set test values
	version = "v1.0.0"
	commit = "abc123"
	buildTime = "2026-02-11T00:00:00Z"

	oldVerbose := versionVerbose
	versionVerbose = true
	defer func() { versionVerbose = oldVerbose }()

	// Capture output
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	runVersion(versionCmd, []string{})

	w.Close()
	os.Stdout = oldStdout
	output, _ := io.ReadAll(r)
	outputStr := string(output)

	// Verify output contains all info
	if !strings.Contains(outputStr, "v1.0.0") {
		t.Errorf("verbose output should contain version")
	}
	if !strings.Contains(outputStr, "commit") {
		t.Errorf("verbose output should contain commit label")
	}
	if !strings.Contains(outputStr, "abc123") {
		t.Errorf("verbose output should contain commit hash")
	}
	if !strings.Contains(outputStr, "built at") {
		t.Errorf("verbose output should contain build time label")
	}
	if !strings.Contains(outputStr, "2026-02-11T00:00:00Z") {
		t.Errorf("verbose output should contain build time")
	}
	if !strings.Contains(outputStr, "go version") {
		t.Errorf("verbose output should contain go version label")
	}
	if !strings.Contains(outputStr, runtime.Version()) {
		t.Errorf("verbose output should contain actual go version")
	}
}

func TestVersionDefaultValues(t *testing.T) {
	// Test that default values are reasonable
	// (These are set at build time, but have defaults)

	// Save originals
	origVersion := version
	origCommit := commit
	origBuildTime := buildTime
	defer func() {
		version = origVersion
		commit = origCommit
		buildTime = origBuildTime
	}()

	// Reset to defaults
	version = "dev"
	commit = "none"
	buildTime = "unknown"

	oldVerbose := versionVerbose
	versionVerbose = true
	defer func() { versionVerbose = oldVerbose }()

	// Capture output
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	runVersion(versionCmd, []string{})

	w.Close()
	os.Stdout = oldStdout
	output, _ := io.ReadAll(r)
	outputStr := string(output)

	// Verify defaults appear in output
	if !strings.Contains(outputStr, "dev") {
		t.Errorf("output should contain default version 'dev'")
	}
	if !strings.Contains(outputStr, "none") {
		t.Errorf("output should contain default commit 'none'")
	}
	if !strings.Contains(outputStr, "unknown") {
		t.Errorf("output should contain default buildTime 'unknown'")
	}
}

func TestVersionCommandSetup(t *testing.T) {
	// Verify versionCmd is properly configured
	if versionCmd.Use != "version" {
		t.Errorf("versionCmd.Use should be 'version', got %q", versionCmd.Use)
	}

	// Verify flags are registered
	flag := versionCmd.Flags().Lookup("verbose")
	if flag == nil {
		t.Error("--verbose flag should be registered")
	}
}

func TestVersionWithDifferentFormats(t *testing.T) {
	testCases := []struct {
		version   string
		commit    string
		buildTime string
		verbose   bool
	}{
		{"v1.0.0", "abc123", "2026-02-11T00:00:00Z", false},
		{"v1.0.0", "abc123", "2026-02-11T00:00:00Z", true},
		{"v2.5.3", "def456", "2026-03-15T10:30:00Z", true},
		{"dev", "none", "unknown", true},
	}

	for _, tc := range testCases {
		// Save originals
		origVersion := version
		origCommit := commit
		origBuildTime := buildTime
		origVerbose := versionVerbose

		// Set test values
		version = tc.version
		commit = tc.commit
		buildTime = tc.buildTime
		versionVerbose = tc.verbose

		// Capture output
		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		runVersion(versionCmd, []string{})

		w.Close()
		os.Stdout = oldStdout
		output, _ := io.ReadAll(r)
		outputStr := string(output)

		// Verify version always appears
		if !strings.Contains(outputStr, tc.version) {
			t.Errorf("output should contain version %q", tc.version)
		}

		// Restore originals
		version = origVersion
		commit = origCommit
		buildTime = origBuildTime
		versionVerbose = origVerbose
	}
}

func TestVersionOutputFormat(t *testing.T) {
	// Save originals
	origVersion := version
	origCommit := commit
	origBuildTime := buildTime
	defer func() {
		version = origVersion
		commit = origCommit
		buildTime = origBuildTime
	}()

	version = "v1.2.3"
	commit = "abc123"
	buildTime = "2026-02-11T00:00:00Z"

	// Test non-verbose format
	oldVerbose := versionVerbose
	versionVerbose = false
	defer func() { versionVerbose = oldVerbose }()

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	runVersion(versionCmd, []string{})

	w.Close()
	os.Stdout = oldStdout
	output, _ := io.ReadAll(r)
	outputStr := string(output)

	// Should be single line format
	lines := strings.Split(strings.TrimSpace(outputStr), "\n")
	if len(lines) != 1 {
		t.Errorf("non-verbose output should be single line, got %d lines", len(lines))
	}

	// Should match pattern "ghqx v1.2.3"
	expectedPattern := "ghqx v1.2.3"
	if !strings.Contains(outputStr, expectedPattern) {
		t.Errorf("output should contain %q, got: %s", expectedPattern, outputStr)
	}
}

func TestVersionVerboseOutputFormat(t *testing.T) {
	// Save originals
	origVersion := version
	origCommit := commit
	origBuildTime := buildTime
	defer func() {
		version = origVersion
		commit = origCommit
		buildTime = origBuildTime
	}()

	version = "v1.2.3"
	commit = "abc123"
	buildTime = "2026-02-11T00:00:00Z"

	oldVerbose := versionVerbose
	versionVerbose = true
	defer func() { versionVerbose = oldVerbose }()

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	runVersion(versionCmd, []string{})

	w.Close()
	os.Stdout = oldStdout
	output, _ := io.ReadAll(r)
	outputStr := string(output)

	// Should be multi-line format with at least 4 lines
	lines := strings.Split(strings.TrimSpace(outputStr), "\n")
	if len(lines) < 4 {
		t.Errorf("verbose output should have at least 4 lines, got %d", len(lines))
	}
}
