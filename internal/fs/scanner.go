package fs

import (
	"os"
	"path/filepath"

	"github.com/mi8bi/ghqx/internal/domain"
)

// Scanner handles filesystem operations for project discovery.
type Scanner struct{}

// NewScanner creates a new filesystem scanner.
func NewScanner() *Scanner {
	return &Scanner{}
}

// ScanRoot scans a root directory and returns all projects found.
func (s *Scanner) ScanRoot(rootName domain.RootName, rootPath string, isSandbox bool) ([]domain.Project, error) {
	if _, err := os.Stat(rootPath); os.IsNotExist(err) {
		return nil, domain.ErrRootDirNotExist(string(rootName), rootPath)
	}

	return s.scanGhqRoot(rootName, rootPath)
}

// scanGhqRoot scans a directory recursively for git repositories.
func (s *Scanner) scanGhqRoot(rootName domain.RootName, rootPath string) ([]domain.Project, error) {
	var projects []domain.Project

	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Skip errors (e.g., permission denied)
		}

		if !info.IsDir() {
			return nil // Only interested in directories
		}

		// Check if this directory is a Git repository
		if !s.hasGitDir(path) {
			return nil // Not a git repo, continue walking
		}

		// Calculate relative path from the rootPath
		relPath, _ := filepath.Rel(rootPath, path)

		// Use the relative path as the project name
		projectName := filepath.ToSlash(relPath)
		if projectName == "." || projectName == "" {
			projectName = filepath.Base(path)
		}

		// Determine project type based on root name
		projectType := domain.ProjectTypeDev
		switch rootName {
		case "sandbox":
			projectType = domain.ProjectTypeSandbox
		case "release":
			projectType = domain.ProjectTypeRelease
		}

		projects = append(projects, domain.Project{
			Name:        projectName,
			DisplayName: domain.FormatDisplayName(projectName),
			Root:        rootName,
			Path:        path,
			Zone:        domain.DetermineZone(rootName),
			Type:        projectType,
			HasGit:      true,
		})

		// Skip descending into this git repository
		return filepath.SkipDir
	})

	if err != nil {
		return nil, domain.ErrFSScanRoot(err)
	}

	return projects, nil
}

// hasGitDir checks if a directory contains a .git subdirectory.
func (s *Scanner) hasGitDir(path string) bool {
	gitPath := filepath.Join(path, ".git")
	info, err := os.Stat(gitPath)
	return err == nil && info.IsDir()
}

// HasGitDir checks if a directory contains a .git subdirectory (public version).
func (s Scanner) HasGitDir(path string) bool {
	return s.hasGitDir(path)
}

// EnsureDir creates a directory if it doesn't exist.
func (s *Scanner) EnsureDir(path string) error {
	if err := os.MkdirAll(path, 0755); err != nil {
		return domain.ErrFSCreateDir(err)
	}
	return nil
}

// IsSafeName checks if a name is safe for filesystem use.
func IsSafeName(name string) bool {
	if name == "" || name == "." || name == ".." {
		return false
	}
	// Disallow path separators and other dangerous characters
	forbidden := []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|"}
	for _, char := range forbidden {
		if contains(name, char) {
			return false
		}
	}
	return true
}

// contains is a helper function for string containment check
func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
