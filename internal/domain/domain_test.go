package domain

import (
	"errors"
	"testing"
)

func TestFormatDisplayName(t *testing.T) {
	cases := map[string]string{
		"github.com/user/repo": "user/repo",
		"example.com/a/b/c":    "b/c",
		"shortname":            "shortname",
	}

	for in, want := range cases {
		if got := FormatDisplayName(in); got != want {
			t.Fatalf("FormatDisplayName(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestDetermineWorkspaceType(t *testing.T) {
	if DetermineWorkspaceType("sandbox") != WorkspaceTypeSandbox {
		t.Fatalf("sandbox not mapped")
	}
	if DetermineWorkspaceType("dev") != WorkspaceTypeDev {
		t.Fatalf("dev not mapped")
	}
	if DetermineWorkspaceType("release") != WorkspaceTypeRelease {
		t.Fatalf("release not mapped")
	}
	if DetermineWorkspaceType("unknown-name") != WorkspaceTypeUnknown {
		t.Fatalf("unknown should map to unknown")
	}
}

func TestGhqxErrorBehaviors(t *testing.T) {
	e := NewError(ErrCodeConfigNotFound, "msg")
	if e.Error() == "" {
		t.Fatalf("Error() should return non-empty")
	}

	if e.IsUserError() {
		t.Fatalf("IsUserError false when no hint")
	}

	e = e.WithHint("do this")
	if !e.IsUserError() {
		t.Fatalf("IsUserError true after WithHint")
	}

	e = e.WithInternal("internal")
	if e.Internal != "internal" {
		t.Fatalf("WithInternal failed")
	}

	cause := errors.New("cause")
	e2 := NewErrorWithCause(ErrCodeUnknown, "m", cause)
	if un := e2.Unwrap(); un == nil {
		t.Fatalf("Unwrap should return cause")
	}

	det := e2.DetailedError()
	if det == "" {
		t.Fatalf("DetailedError should be non-empty")
	}
}
