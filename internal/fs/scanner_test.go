package fs

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mi8bi/ghqx/internal/domain"
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

// Additional tests for better coverage

func TestNewScanner(t *testing.T) {
	s := NewScanner()
	if s == nil {
		t.Fatal("NewScanner returned nil")
	}
}

func TestScanRootWithNonExistentPath(t *testing.T) {
	s := NewScanner()
	_, err := s.ScanRoot("test", "/nonexistent/path")
	if err == nil {
		t.Fatal("expected error for nonexistent root")
	}
}

func TestScanRootWithGitRepo(t *testing.T) {
	tmp, err := os.MkdirTemp("", "ghqx-scan-git")
	if err != nil {
		t.Fatalf("tempdir: %v", err)
	}
	defer os.RemoveAll(tmp)

	// Create git repository
	repo := filepath.Join(tmp, "github.com", "user", "repo")
	if err := os.MkdirAll(repo, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	// Add .git directory
	gitDir := filepath.Join(repo, ".git")
	if err := os.Mkdir(gitDir, 0755); err != nil {
		t.Fatalf("mkdir .git: %v", err)
	}

	s := NewScanner()
	projects, err := s.ScanRoot("dev", tmp)
	if err != nil {
		t.Fatalf("ScanRoot failed: %v", err)
	}

	if len(projects) == 0 {
		t.Fatal("expected at least one project")
	}

	// Should detect git
	found := false
	for _, p := range projects {
		if p.HasGit {
			found = true
			if p.Type != domain.ProjectTypeDev {
				t.Errorf("expected ProjectTypeDev, got %v", p.Type)
			}
			break
		}
	}
	if !found {
		t.Error("should detect git repository")
	}
}

func TestScanRootWithMultipleProjects(t *testing.T) {
	tmp, err := os.MkdirTemp("", "ghqx-scan-multiple")
	if err != nil {
		t.Fatalf("tempdir: %v", err)
	}
	defer os.RemoveAll(tmp)

	// Create multiple projects
	projects := []string{
		"github.com/user1/repo1",
		"github.com/user1/repo2",
		"github.com/user2/repo1",
		"gitlab.com/org/project",
	}

	for _, proj := range projects {
		path := filepath.Join(tmp, proj)
		if err := os.MkdirAll(path, 0755); err != nil {
			t.Fatalf("mkdir %s: %v", proj, err)
		}
	}

	s := NewScanner()
	found, err := s.ScanRoot("sandbox", tmp)
	if err != nil {
		t.Fatalf("ScanRoot failed: %v", err)
	}

	if len(found) != len(projects) {
		t.Errorf("expected %d projects, got %d", len(projects), len(found))
	}
}

func TestScanRootWithNestedGit(t *testing.T) {
	tmp, err := os.MkdirTemp("", "ghqx-scan-nested")
	if err != nil {
		t.Fatalf("tempdir: %v", err)
	}
	defer os.RemoveAll(tmp)

	// Create parent repo
	parent := filepath.Join(tmp, "github.com", "user", "parent")
	if err := os.MkdirAll(parent, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.Mkdir(filepath.Join(parent, ".git"), 0755); err != nil {
		t.Fatalf("mkdir .git: %v", err)
	}

	// Create nested repo (should be skipped)
	nested := filepath.Join(parent, "nested")
	if err := os.MkdirAll(nested, 0755); err != nil {
		t.Fatalf("mkdir nested: %v", err)
	}
	if err := os.Mkdir(filepath.Join(nested, ".git"), 0755); err != nil {
		t.Fatalf("mkdir nested .git: %v", err)
	}

	s := NewScanner()
	projects, err := s.ScanRoot("dev", tmp)
	if err != nil {
		t.Fatalf("ScanRoot failed: %v", err)
	}

	// Should only find parent, nested should be skipped
	if len(projects) != 1 {
		t.Errorf("expected 1 project (nested should be skipped), got %d", len(projects))
	}
}

func TestScanRootWithIntermediateDirectories(t *testing.T) {
	tmp, err := os.MkdirTemp("", "ghqx-scan-intermediate")
	if err != nil {
		t.Fatalf("tempdir: %v", err)
	}
	defer os.RemoveAll(tmp)

	// Create structure with intermediate directories
	repo := filepath.Join(tmp, "github.com", "user", "repo")
	if err := os.MkdirAll(repo, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	s := NewScanner()
	projects, err := s.ScanRoot("sandbox", tmp)
	if err != nil {
		t.Fatalf("ScanRoot failed: %v", err)
	}

	// Should find only the complete project, not intermediate directories
	if len(projects) != 1 {
		t.Errorf("expected 1 project, got %d", len(projects))
	}

	if projects[0].Name != "github.com/user/repo" {
		t.Errorf("unexpected project name: %s", projects[0].Name)
	}
}

func TestHasGitDirWithFile(t *testing.T) {
	tmp, err := os.MkdirTemp("", "ghqx-git-file")
	if err != nil {
		t.Fatalf("tempdir: %v", err)
	}
	defer os.RemoveAll(tmp)

	// Create .git as a file (not directory)
	gitFile := filepath.Join(tmp, ".git")
	if err := os.WriteFile(gitFile, []byte("test"), 0644); err != nil {
		t.Fatalf("write .git file: %v", err)
	}

	s := NewScanner()
	if s.HasGitDir(tmp) {
		t.Error("HasGitDir should return false when .git is a file")
	}
}

func TestIsSafeNameEdgeCases(t *testing.T) {
	testCases := []struct {
		name     string
		expected bool
	}{
		{"", false},
		{".", false},
		{"..", false},
		{"valid-name", true},
		{"valid_name", true},
		{"name/with/slash", false},
		{"name\\with\\backslash", false},
		{"name:with:colon", false},
		{"name*with*asterisk", false},
		{"name?with?question", false},
		{"name\"with\"quote", false},
		{"name<with<less", false},
		{"name>with>greater", false},
		{"name|with|pipe", false},
		{"valid.name.123", true},
	}

	for _, tc := range testCases {
		result := IsSafeName(tc.name)
		if result != tc.expected {
			t.Errorf("IsSafeName(%q) = %v, want %v", tc.name, result, tc.expected)
		}
	}
}

func TestEnsureDirWithExistingDir(t *testing.T) {
	tmp, err := os.MkdirTemp("", "ghqx-ensure-existing")
	if err != nil {
		t.Fatalf("tempdir: %v", err)
	}
	defer os.RemoveAll(tmp)

	s := NewScanner()

	// Should succeed without error
	if err := s.EnsureDir(tmp); err != nil {
		t.Errorf("EnsureDir should succeed for existing directory: %v", err)
	}
}

func TestScanRootWithDifferentWorkspaceTypes(t *testing.T) {
	testCases := []struct {
		rootName     domain.RootName
		expectedType domain.ProjectType
	}{
		{"sandbox", domain.ProjectTypeSandboxGit},
		{"dev", domain.ProjectTypeDev},
		{"release", domain.ProjectTypeRelease},
	}

	for _, tc := range testCases {
		tmp, err := os.MkdirTemp("", "ghqx-workspace-"+string(tc.rootName))
		if err != nil {
			t.Fatalf("tempdir: %v", err)
		}
		defer os.RemoveAll(tmp)

		// Create git repo
		repo := filepath.Join(tmp, "github.com", "user", "repo")
		if err := os.MkdirAll(repo, 0755); err != nil {
			t.Fatalf("mkdir: %v", err)
		}
		if err := os.Mkdir(filepath.Join(repo, ".git"), 0755); err != nil {
			t.Fatalf("mkdir .git: %v", err)
		}

		s := NewScanner()
		projects, err := s.ScanRoot(tc.rootName, tmp)
		if err != nil {
			t.Fatalf("ScanRoot failed: %v", err)
		}

		if len(projects) == 0 {
			t.Fatal("expected at least one project")
		}

		if projects[0].Type != tc.expectedType {
			t.Errorf("expected type %v for %s, got %v", tc.expectedType, tc.rootName, projects[0].Type)
		}
	}
}

func TestScanRootWithNonGitDirectory(t *testing.T) {
	tmp, err := os.MkdirTemp("", "ghqx-non-git")
	if err != nil {
		t.Fatalf("tempdir: %v", err)
	}
	defer os.RemoveAll(tmp)

	// Create directory without .git
	dir := filepath.Join(tmp, "github.com", "user", "repo")
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	s := NewScanner()
	projects, err := s.ScanRoot("sandbox", tmp)
	if err != nil {
		t.Fatalf("ScanRoot failed: %v", err)
	}

	if len(projects) == 0 {
		t.Fatal("expected to find non-git directory")
	}

	if projects[0].HasGit {
		t.Error("project should not have git")
	}

	if projects[0].Type != domain.ProjectTypeDir {
		t.Errorf("expected ProjectTypeDir, got %v", projects[0].Type)
	}
}
