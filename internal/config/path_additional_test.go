package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mi8bi/ghqx/internal/domain"
)

func TestEnsureRootDirectoriesFileConflict(t *testing.T) {
	tmp, err := os.MkdirTemp("", "ghqx-path-test")
	if err != nil {
		t.Fatalf("tempdir: %v", err)
	}
	defer os.RemoveAll(tmp)

	// Create a file where a directory is expected
	file := filepath.Join(tmp, "conflict")
	if err := os.WriteFile(file, []byte("x"), 0644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	cfg := &Config{Roots: map[string]string{"r1": file}}
	if err := EnsureRootDirectories(cfg); err == nil {
		t.Fatalf("expected EnsureRootDirectories to fail when root path is a file")
	} else {
		if ge, ok := err.(*domain.GhqxError); ok {
			if ge.Code != domain.ErrCodeFSError {
				t.Fatalf("unexpected error code: %v", ge.Code)
			}
		}
	}
}
