package configtui

import (
	"testing"

	"github.com/mi8bi/ghqx/internal/config"
)

func TestConfigEditorBuildAndApply(t *testing.T) {
	cfg := &config.Config{
		Roots:   map[string]string{"dev": "d", "release": "r", "sandbox": "s"},
		Default: config.DefaultConfig{Root: "dev"},
	}

	e := NewConfigEditor(cfg, "/tmp/config")
	if len(e.Fields) != 4 {
		t.Fatalf("expected 4 fields, got %d", len(e.Fields))
	}

	// Update a field and apply
	e.UpdateField(0, "d2")
	if !e.Modified {
		t.Fatalf("expected Modified after update")
	}

	e.ApplyChanges()
	if cfg.Roots["dev"] != "d2" {
		t.Fatalf("expected cfg.Roots[dev] updated")
	}

	// Helpers
	if boolToString(true) != "true" || boolToString(false) != "false" {
		t.Fatalf("boolToString mismatch")
	}
	if !stringToBool("yes") || stringToBool("no") {
		t.Fatalf("stringToBool mismatch")
	}
	if intToString(3) != "3" {
		t.Fatalf("intToString mismatch")
	}
	if stringToInt("bad") != 0 || stringToInt("5") != 5 {
		t.Fatalf("stringToInt mismatch")
	}
}

func TestBoolStringIntHelpers(t *testing.T) {
	if boolToString(true) != "true" {
		t.Fatalf("boolToString(true) != \"true\"")
	}
	if boolToString(false) != "false" {
		t.Fatalf("boolToString(false) != \"false\"")
	}

	for _, s := range []string{"true", "yes", "1"} {
		if !stringToBool(s) {
			t.Fatalf("stringToBool(%q) should be true", s)
		}
	}
	if stringToBool("no") {
		t.Fatalf("stringToBool(\"no\") should be false")
	}

	if intToString(123) != "123" {
		t.Fatalf("intToString(123) != \"123\"")
	}

	if stringToInt("45") != 45 {
		t.Fatalf("stringToInt(\"45\") != 45")
	}
	if stringToInt("bad") != 0 {
		t.Fatalf("stringToInt(\"bad\") should return 0 for invalid input")
	}
}
