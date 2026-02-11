package domain

import (
	"errors"
	"testing"
)

func TestErrorConstructorsAndVars(t *testing.T) {
	if ErrConfigNotFoundAny == nil {
		t.Fatalf("ErrConfigNotFoundAny is nil")
	}

	e := ErrConfigNotFoundAt("/x")
	if e == nil || e.Code != ErrCodeConfigNotFound {
		t.Fatalf("ErrConfigNotFoundAt returned unexpected: %#v", e)
	}

	et := ErrConfigInvalidTOML(errors.New("cause"))
	if et == nil || et.Code != ErrCodeConfigInvalid {
		t.Fatalf("ErrConfigInvalidTOML returned unexpected: %#v", et)
	}

	if ErrConfigNoRoots == nil {
		t.Fatalf("ErrConfigNoRoots is nil")
	}

	if ErrConfigInvalidDefaultRoot == nil {
		t.Fatalf("ErrConfigInvalidDefaultRoot is nil")
	}

	r := ErrRootNotFound("name")
	if r == nil || r.Code != ErrCodeRootNotFound {
		t.Fatalf("ErrRootNotFound unexpected: %#v", r)
	}

	p := ErrProjectNotFound("p")
	if p == nil || p.Code != ErrCodeProjectNotFound {
		t.Fatalf("ErrProjectNotFound unexpected: %#v", p)
	}

	if ErrFSCreateDir(errors.New("x")) == nil {
		t.Fatalf("ErrFSCreateDir returned nil")
	}
	if ErrFSReadDir(errors.New("x")) == nil {
		t.Fatalf("ErrFSReadDir returned nil")
	}
	if ErrFSScanRoot(errors.New("x")) == nil {
		t.Fatalf("ErrFSScanRoot returned nil")
	}

	gt := ErrGitTimeout("op")
	if gt == nil || gt.Code != ErrCodeGitError {
		t.Fatalf("ErrGitTimeout unexpected: %#v", gt)
	}

	gc := ErrGitCommandFailed("op", errors.New("cause"))
	if gc == nil || gc.Code != ErrCodeGitError {
		t.Fatalf("ErrGitCommandFailed unexpected: %#v", gc)
	}
}
