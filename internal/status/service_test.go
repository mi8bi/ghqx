package status

import (
	"os"
	"os/exec"
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

// Additional tests for better coverage

func TestNewService(t *testing.T) {
	cfg := &config.Config{
		Roots:   map[string]string{"sandbox": "/tmp"},
		Default: config.DefaultConfig{Root: "sandbox"},
	}
	s := NewService(cfg)

	if s == nil {
		t.Fatal("NewService returned nil")
	}

	if s.cfg != cfg {
		t.Error("config not set correctly")
	}

	if s.scanner == nil {
		t.Error("scanner should not be nil")
	}

	if s.git == nil {
		t.Error("git client should not be nil")
	}
}

func TestGetAllWithCheckDirty(t *testing.T) {
	// Skip if git is not available
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available, skipping test")
	}

	tmp, err := os.MkdirTemp("", "ghqx-status-dirty")
	if err != nil {
		t.Fatalf("tempdir: %v", err)
	}
	defer os.RemoveAll(tmp)

	// Create git repository
	repo := filepath.Join(tmp, "github.com", "user", "repo")
	if err := os.MkdirAll(repo, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	// Initialize git repo
	cmd := exec.Command("git", "init")
	cmd.Dir = repo
	if err := cmd.Run(); err != nil {
		t.Skip("failed to init git repo")
	}

	// Configure git
	exec.Command("git", "config", "user.email", "test@example.com").Run()
	exec.Command("git", "config", "user.name", "Test").Run()

	cfg := &config.Config{Roots: map[string]string{"dev": tmp}}
	s := NewService(cfg)

	projects, err := s.GetAll(Options{CheckDirty: true})
	if err != nil {
		t.Fatalf("GetAll failed: %v", err)
	}

	if len(projects) == 0 {
		t.Fatal("expected at least one project")
	}

	// Project should have git
	if !projects[0].HasGit {
		t.Error("project should be detected as git repo")
	}
}

func TestGetAllWithLoadBranch(t *testing.T) {
	// Skip if git is not available
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available, skipping test")
	}

	tmp, err := os.MkdirTemp("", "ghqx-status-branch")
	if err != nil {
		t.Fatalf("tempdir: %v", err)
	}
	defer os.RemoveAll(tmp)

	// Create git repository
	repo := filepath.Join(tmp, "github.com", "user", "repo")
	if err := os.MkdirAll(repo, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	// Initialize git repo
	cmd := exec.Command("git", "init")
	cmd.Dir = repo
	if err := cmd.Run(); err != nil {
		t.Skip("failed to init git repo")
	}

	// Configure git
	exec.Command("git", "config", "user.email", "test@example.com").Run()
	exec.Command("git", "config", "user.name", "Test").Run()

	// Create initial commit
	testFile := filepath.Join(repo, "test.txt")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Skip("failed to create test file")
	}

	cmd = exec.Command("git", "add", ".")
	cmd.Dir = repo
	cmd.Run()

	cmd = exec.Command("git", "commit", "-m", "initial")
	cmd.Dir = repo
	if err := cmd.Run(); err != nil {
		t.Skip("failed to commit")
	}

	cfg := &config.Config{Roots: map[string]string{"dev": tmp}}
	s := NewService(cfg)

	projects, err := s.GetAll(Options{LoadBranch: true})
	if err != nil {
		t.Fatalf("GetAll failed: %v", err)
	}

	if len(projects) == 0 {
		t.Fatal("expected at least one project")
	}

	// Branch should be loaded
	if projects[0].Branch == "" {
		t.Error("branch should be loaded")
	}
}

func TestGetAllWithMultipleRoots(t *testing.T) {
	tmp, err := os.MkdirTemp("", "ghqx-status-multi")
	if err != nil {
		t.Fatalf("tempdir: %v", err)
	}
	defer os.RemoveAll(tmp)

	// Create multiple roots
	devRoot := filepath.Join(tmp, "dev")
	sandboxRoot := filepath.Join(tmp, "sandbox")

	for _, root := range []string{devRoot, sandboxRoot} {
		repo := filepath.Join(root, "github.com", "user", "repo")
		if err := os.MkdirAll(repo, 0755); err != nil {
			t.Fatalf("mkdir: %v", err)
		}
	}

	cfg := &config.Config{
		Roots: map[string]string{
			"dev":     devRoot,
			"sandbox": sandboxRoot,
		},
		Default: config.DefaultConfig{Root: "dev"},
	}
	s := NewService(cfg)

	projects, err := s.GetAll(Options{})
	if err != nil {
		t.Fatalf("GetAll failed: %v", err)
	}

	if len(projects) != 2 {
		t.Errorf("expected 2 projects (one from each root), got %d", len(projects))
	}
}

func TestGetAllWithRootFilter(t *testing.T) {
	tmp, err := os.MkdirTemp("", "ghqx-status-filter")
	if err != nil {
		t.Fatalf("tempdir: %v", err)
	}
	defer os.RemoveAll(tmp)

	// Create multiple roots
	devRoot := filepath.Join(tmp, "dev")
	sandboxRoot := filepath.Join(tmp, "sandbox")

	for _, root := range []string{devRoot, sandboxRoot} {
		repo := filepath.Join(root, "github.com", "user", "repo")
		if err := os.MkdirAll(repo, 0755); err != nil {
			t.Fatalf("mkdir: %v", err)
		}
	}

	cfg := &config.Config{
		Roots: map[string]string{
			"dev":     devRoot,
			"sandbox": sandboxRoot,
		},
		Default: config.DefaultConfig{Root: "dev"},
	}
	s := NewService(cfg)

	// Filter by dev only
	projects, err := s.GetAll(Options{}, "dev")
	if err != nil {
		t.Fatalf("GetAll failed: %v", err)
	}

	if len(projects) != 1 {
		t.Errorf("expected 1 project from dev root, got %d", len(projects))
	}

	if projects[0].Root != "dev" {
		t.Errorf("expected project from dev root, got %s", projects[0].Root)
	}
}

func TestFindProjectNotFound(t *testing.T) {
	tmp, err := os.MkdirTemp("", "ghqx-status-notfound")
	if err != nil {
		t.Fatalf("tempdir: %v", err)
	}
	defer os.RemoveAll(tmp)

	cfg := &config.Config{Roots: map[string]string{"sandbox": tmp}}
	s := NewService(cfg)

	_, err = s.FindProject("nonexistent/project")
	if err == nil {
		t.Fatal("expected error when project not found")
	}
}

func TestEnrichProjectsWithNonGitProject(t *testing.T) {
	tmp, err := os.MkdirTemp("", "ghqx-status-nongit")
	if err != nil {
		t.Fatalf("tempdir: %v", err)
	}
	defer os.RemoveAll(tmp)

	// Create non-git directory
	dir := filepath.Join(tmp, "github.com", "user", "repo")
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	cfg := &config.Config{Roots: map[string]string{"sandbox": tmp}}
	s := NewService(cfg)

	projects, err := s.GetAll(Options{CheckDirty: true})
	if err != nil {
		t.Fatalf("GetAll failed: %v", err)
	}

	if len(projects) == 0 {
		t.Fatal("expected at least one project")
	}

	// Should not have git
	if projects[0].HasGit {
		t.Error("project should not have git")
	}

	// Dirty flag should be false for non-git
	if projects[0].Dirty {
		t.Error("non-git project should not be dirty")
	}
}

func TestGetAllWithEmptyRoot(t *testing.T) {
	tmp, err := os.MkdirTemp("", "ghqx-status-empty")
	if err != nil {
		t.Fatalf("tempdir: %v", err)
	}
	defer os.RemoveAll(tmp)

	cfg := &config.Config{Roots: map[string]string{"sandbox": tmp}}
	s := NewService(cfg)

	projects, err := s.GetAll(Options{})
	if err != nil {
		t.Fatalf("GetAll failed: %v", err)
	}

	if len(projects) != 0 {
		t.Errorf("expected no projects in empty root, got %d", len(projects))
	}
}

func TestDetermineTargetRootsWithEmptyFilter(t *testing.T) {
	cfg := &config.Config{
		Roots: map[string]string{
			"dev":     "/tmp/dev",
			"sandbox": "/tmp/sandbox",
		},
	}
	s := NewService(cfg)

	// Empty filter should return all roots
	roots := s.determineTargetRoots([]string{""})
	if len(roots) != 2 {
		t.Errorf("expected 2 roots with empty filter, got %d", len(roots))
	}
}
