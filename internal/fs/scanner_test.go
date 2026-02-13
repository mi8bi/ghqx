package fs

import (
	"os"
	"path/filepath"
	"testing"
)

func TestIsSafeNameAndContains(t *testing.T) {
	if !IsSafeName("good-name") {
		t.Fatalf("IsSafeName should accept good-name")
	}
	if IsSafeName("bad/name") {
		t.Fatalf("IsSafeName should reject path separators")
	}
}

func TestHasGitDirAndEnsureDir(t *testing.T) {
	tmp, err := os.MkdirTemp("", "ghqx-fs-test")
	if err != nil {
		t.Fatalf("tempdir: %v", err)
	}
	defer os.RemoveAll(tmp)

	s := NewScanner()

	// EnsureDir should create a nested dir
	nested := filepath.Join(tmp, "a", "b")
	if err := s.EnsureDir(nested); err != nil {
		t.Fatalf("EnsureDir failed: %v", err)
	}

	// create .git dir and test HasGitDir
	repo := filepath.Join(tmp, "repo")
	if err := os.MkdirAll(filepath.Join(repo, ".git"), 0755); err != nil {
		t.Fatalf("mkdir .git: %v", err)
	}

	if !s.HasGitDir(repo) {
		t.Fatalf("HasGitDir should detect .git")
	}
}

func TestScanRootReturnsProjects(t *testing.T) {
	tmp, err := os.MkdirTemp("", "ghqx-scan-root")
	if err != nil {
		t.Fatalf("tempdir: %v", err)
	}
	defer os.RemoveAll(tmp)

	// create structure github.com/user/repo
	repo := filepath.Join(tmp, "github.com", "user", "repo")
	if err := os.MkdirAll(repo, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	s := NewScanner()
	projects, err := s.ScanRoot("sandbox", tmp)
	if err != nil {
		t.Fatalf("ScanRoot failed: %v", err)
	}
	if len(projects) == 0 {
		t.Fatalf("expected at least one project")
	}
}
