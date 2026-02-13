package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mi8bi/ghqx/internal/app"
	"github.com/mi8bi/ghqx/internal/config"
	"github.com/mi8bi/ghqx/internal/i18n"
)

func TestMatchLocaleStringCases(t *testing.T) {
	testCases := []struct {
		input    string
		expected i18n.Locale
	}{
		{"ja_JP.UTF-8", i18n.LocaleJA},
		{"en_US.UTF-8", i18n.LocaleEN},
		{"en", i18n.LocaleEN},
		{"ja", i18n.LocaleJA},
		{"EN", i18n.LocaleEN},
		{"JA", i18n.LocaleJA},
		{"en:ja", i18n.LocaleEN},       // First in list should win
		{"ja:en", i18n.LocaleJA},       // First in list should win
		{"en_US:ja_JP", i18n.LocaleEN}, // First in list should win
		{"fr_FR", ""},
		{"de_DE", ""},
		{"", ""},
	}

	for _, tc := range testCases {
		result := matchLocaleString(tc.input)
		if result != tc.expected {
			t.Errorf("matchLocaleString(%q) = %q, want %q", tc.input, result, tc.expected)
		}
	}
}

func TestGetOSLanguageLocaleWithLANG(t *testing.T) {
	// Save original env vars
	origLC := os.Getenv("LC_ALL")
	origLANG := os.Getenv("LANG")
	origLANGUAGE := os.Getenv("LANGUAGE")
	defer func() {
		// Restore original values
		if origLC != "" {
			os.Setenv("LC_ALL", origLC)
		} else {
			os.Unsetenv("LC_ALL")
		}
		if origLANG != "" {
			os.Setenv("LANG", origLANG)
		} else {
			os.Unsetenv("LANG")
		}
		if origLANGUAGE != "" {
			os.Setenv("LANGUAGE", origLANGUAGE)
		} else {
			os.Unsetenv("LANGUAGE")
		}
	}()

	// Test 1: LANG only
	os.Unsetenv("LC_ALL")
	os.Unsetenv("LANGUAGE")
	os.Setenv("LANG", "en_US.UTF-8")

	locale := getOSLanguageLocale()
	if locale != i18n.LocaleEN {
		t.Errorf("Test 1: expected LocaleEN with LANG=en_US.UTF-8, got %v", locale)
	}

	// Test 2: LC_ALL has highest priority
	os.Setenv("LC_ALL", "ja_JP.UTF-8")
	os.Setenv("LANG", "en_US.UTF-8")
	os.Unsetenv("LANGUAGE")

	locale = getOSLanguageLocale()
	if locale != i18n.LocaleJA {
		t.Errorf("Test 2: expected LocaleJA with LC_ALL=ja_JP.UTF-8, got %v", locale)
	}

	// Test 3: LANGUAGE (first language in list should win)
	os.Unsetenv("LC_ALL")
	os.Unsetenv("LANG")
	os.Setenv("LANGUAGE", "en:ja")

	locale = getOSLanguageLocale()
	if locale != i18n.LocaleEN {
		t.Errorf("Test 3: expected LocaleEN with LANGUAGE=en:ja (first in list), got %v", locale)
	}

	// Test 3b: LANGUAGE with ja first
	os.Setenv("LANGUAGE", "ja:en")

	locale = getOSLanguageLocale()
	if locale != i18n.LocaleJA {
		t.Errorf("Test 3b: expected LocaleJA with LANGUAGE=ja:en (first in list), got %v", locale)
	}

	// Test 4: Default (no env vars)
	os.Unsetenv("LC_ALL")
	os.Unsetenv("LANG")
	os.Unsetenv("LANGUAGE")

	locale = getOSLanguageLocale()
	if locale != i18n.LocaleJA {
		t.Errorf("Test 4: expected LocaleJA as default, got %v", locale)
	}
}

func TestDetermineLocaleWithGHQX_LANG(t *testing.T) {
	// Save original
	origGHQXLang := os.Getenv("GHQX_LANG")
	defer os.Setenv("GHQX_LANG", origGHQXLang)

	// Test GHQX_LANG=en
	os.Setenv("GHQX_LANG", "en")
	locale := determineLocale()
	if locale != i18n.LocaleEN {
		t.Errorf("expected LocaleEN with GHQX_LANG=en, got %v", locale)
	}

	// Test GHQX_LANG=ja
	os.Setenv("GHQX_LANG", "ja")
	locale = determineLocale()
	if locale != i18n.LocaleJA {
		t.Errorf("expected LocaleJA with GHQX_LANG=ja, got %v", locale)
	}

	// Test GHQX_LANG=invalid
	os.Setenv("GHQX_LANG", "invalid")
	locale = determineLocale()
	// Should fall back to OS language
	if locale == "" {
		t.Error("determineLocale should not return empty string")
	}
}

func TestInitLocale(t *testing.T) {
	// Test that initLocale sets a valid locale
	origGHQXLang := os.Getenv("GHQX_LANG")
	defer os.Setenv("GHQX_LANG", origGHQXLang)

	os.Setenv("GHQX_LANG", "en")
	initLocale()

	// Verify that locale was set (by checking a translation works)
	if i18n.T("root.command.short") == "" {
		t.Error("initLocale should set locale that enables translations")
	}
}

func TestLoadAppSuccess(t *testing.T) {
	tmp, err := os.MkdirTemp("", "ghqx-loadapp")
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

	err = loadApp()
	if err != nil {
		t.Fatalf("loadApp failed: %v", err)
	}

	if application == nil {
		t.Fatal("application should not be nil after loadApp")
	}
}

func TestLoadAppFailure(t *testing.T) {
	oldConfigPath := configPath
	configPath = "/nonexistent/config.toml"
	defer func() { configPath = oldConfigPath }()

	oldApp := application
	application = nil
	defer func() { application = oldApp }()

	err := loadApp()
	if err == nil {
		t.Fatal("expected error when config doesn't exist")
	}
}

func TestRootCommandSetup(t *testing.T) {
	// Verify that rootCmd has expected subcommands
	subcommands := []string{"status", "cd", "version", "config", "get", "doctor", "clean", "mode"}

	for _, cmdName := range subcommands {
		found := false
		for _, cmd := range rootCmd.Commands() {
			if cmd.Name() == cmdName {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected subcommand %q not found", cmdName)
		}
	}
}

func TestPersistentPreRunEForConfigInit(t *testing.T) {
	// Test that configInitCmd skips app loading
	oldConfigPath := configPath
	configPath = "/nonexistent/config.toml"
	defer func() { configPath = oldConfigPath }()

	oldApp := application
	application = nil
	defer func() { application = oldApp }()

	// PersistentPreRunE should not fail for configInitCmd
	err := rootCmd.PersistentPreRunE(configInitCmd, []string{})
	if err != nil {
		t.Errorf("PersistentPreRunE should not fail for configInitCmd: %v", err)
	}
}

func TestPersistentPreRunEForOtherCommands(t *testing.T) {
	tmp, err := os.MkdirTemp("", "ghqx-prerun")
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
	application = nil
	defer func() { application = oldApp }()

	// Test with statusCmd
	err = rootCmd.PersistentPreRunE(statusCmd, []string{})
	if err != nil {
		t.Errorf("PersistentPreRunE should succeed with valid config: %v", err)
	}

	if application == nil {
		t.Error("application should be loaded after PersistentPreRunE")
	}
}

func TestMainFunctionExists(t *testing.T) {
	// Verify main function would work (can't test directly as it calls os.Exit)
	// But we can verify rootCmd is properly initialized
	if rootCmd.Use != "ghqx" {
		t.Error("rootCmd.Use should be 'ghqx'")
	}
}

func TestNewFromConfigPathWrapper(t *testing.T) {
	tmp, err := os.MkdirTemp("", "ghqx-newfromconfig")
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

	// Test app.NewFromConfigPath through loadApp
	oldConfigPath := configPath
	configPath = cfgPath
	defer func() { configPath = oldConfigPath }()

	oldApp := application
	defer func() { application = oldApp }()

	appInstance, err := app.NewFromConfigPath(cfgPath)
	if err != nil {
		t.Fatalf("NewFromConfigPath failed: %v", err)
	}

	if appInstance == nil {
		t.Fatal("appInstance should not be nil")
	}
}
