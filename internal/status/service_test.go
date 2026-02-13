package status

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mi8bi/ghqx/internal/config"
	"github.com/mi8bi/ghqx/internal/fs"
)

func TestGetAllAndFindProject(t *testing.T) {
	tmp, err := os.MkdirTemp("", "ghqx-status-test")
	if err != nil {
		t.Fatalf("tempdir: %v", err)
	}
	defer os.RemoveAll(tmp)

	// create github.com/user/repo structure
	repo := filepath.Join(tmp, "github.com", "user", "repo")
	if err := os.MkdirAll(repo, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	cfg := &config.Config{Roots: map[string]string{"sandbox": tmp}, Default: config.DefaultConfig{Root: "sandbox"}}
	s := NewService(cfg)

	projects, err := s.GetAll(Options{})
	if err != nil {
		t.Fatalf("GetAll failed: %v", err)
	}
	if len(projects) == 0 {
		t.Fatalf("expected projects found")
	}

	// FindProject should locate by name
	name := projects[0].Name
	p, err := s.FindProject(name)
	if err != nil {
		t.Fatalf("FindProject failed: %v", err)
	}
	if p.Name != name {
		t.Fatalf("FindProject returned wrong project")
	}
}

func TestDetermineTargetRoots(t *testing.T) {
	cfg := &config.Config{Roots: map[string]string{"dev": "/tmp/dev", "sandbox": "/tmp/s"}}
	s := &Service{cfg: cfg}

	// no filter -> return all
	all := s.determineTargetRoots([]string{})
	if len(all) != 2 {
		t.Fatalf("expected 2 roots, got %d", len(all))
	}

	// filter existing
	one := s.determineTargetRoots([]string{"dev"})
	if p, ok := one["dev"]; !ok || p != "/tmp/dev" {
		t.Fatalf("determineTargetRoots with filter failed: %v", one)
	}

	// filter non-existing -> return all
	all2 := s.determineTargetRoots([]string{"nope"})
	if len(all2) != 2 {
		t.Fatalf("expected fallback to all roots when filter missing")
	}
}

func TestGetAllScansTempRoot(t *testing.T) {
	// create temp root with structure github.com/user/repo
	root, err := os.MkdirTemp("", "ghqx-test-root")
	if err != nil {
		t.Fatalf("tempdir: %v", err)
	}
	defer os.RemoveAll(root)

	repoPath := filepath.Join(root, "github.com", "user", "repo")
	if err := os.MkdirAll(repoPath, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	cfg := &config.Config{Roots: map[string]string{"sandbox": root}}
	s := &Service{cfg: cfg, scanner: fs.NewScanner()}

	projects, err := s.GetAll(Options{})
	if err != nil {
		t.Fatalf("GetAll error: %v", err)
	}
	if len(projects) == 0 {
		t.Fatalf("expected at least one project from scan")
	}
}
