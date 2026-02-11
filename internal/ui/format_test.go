package ui

import (
	"os"
	"testing"

	"github.com/mi8bi/ghqx/internal/domain"
	"github.com/mi8bi/ghqx/internal/i18n"
)

func TestFormatMessages(t *testing.T) {
	t.Run("success/warning/info", func(t *testing.T) {
		if got := FormatSuccess("ok"); got != i18n.T("ui.success.prefix")+" ok\n" {
			t.Fatalf("FormatSuccess returned %q", got)
		}
		if got := FormatWarning("warn"); got != i18n.T("ui.warning.prefix")+" warn\n" {
			t.Fatalf("FormatWarning returned %q", got)
		}
		if got := FormatInfo("info"); got != i18n.T("ui.info.prefix")+" info\n" {
			t.Fatalf("FormatInfo returned %q", got)
		}
	})

	t.Run("error and detailed error", func(t *testing.T) {
		e := domain.NewError(domain.ErrCodeConfigNotFound, "config missing").WithHint("create config")

		// Without debug env
		os.Unsetenv("GHQX_DEBUG")
		out := FormatError(e)
		want := "\n" + i18n.T("ui.error.prefix") + ": " + e.Message + "\n\n" + i18n.T("ui.error.hintPrefix") + ":\n  " + e.Hint + "\n"
		if out != want {
			t.Fatalf("FormatError mismatch:\n got: %q\nwant: %q", out, want)
		}

		// With debug env and internal/cause
		e = e.WithInternal("internal-info").WithHint(e.Hint)
		e.Cause = domain.NewError(domain.ErrCodeUnknown, "cause")
		os.Setenv("GHQX_DEBUG", "1")
		defer os.Unsetenv("GHQX_DEBUG")

		dout := FormatDetailedError(e)
		if dout == "" {
			t.Fatalf("FormatDetailedError returned empty string")
		}
	})
}
