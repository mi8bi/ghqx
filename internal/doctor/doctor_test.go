package doctor

import (
	"os"
	"testing"

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
