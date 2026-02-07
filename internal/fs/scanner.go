package fs

import (
	"os"
	"path/filepath"
	"strings"

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

	if isSandbox {
		return s.scanSandbox(rootName, rootPath)
	}
	return s.scanGhqRoot(rootName, rootPath)
}

// scanSandbox scans a flat sandbox directory.
func (s *Scanner) scanSandbox(rootName domain.RootName, rootPath string) ([]domain.Project, error) {
	entries, err := os.ReadDir(rootPath)
	if err != nil {
		return nil, domain.ErrFSReadDir(err)
	}

	var projects []domain.Project
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		name := entry.Name()
		projectPath := filepath.Join(rootPath, name)

		hasGit := s.hasGitDir(projectPath)
		projectType := domain.ProjectTypeSandbox
		if hasGit {
			projectType = domain.ProjectTypeSandboxGit
		}

		projects = append(projects, domain.Project{
			Name:   name,
			Root:   rootName,
			Path:   projectPath,
			Type:   projectType,
			HasGit: hasGit,
		})
	}

	return projects, nil
}

// scanGhqRoot scans a ghq-style hierarchical directory (host/owner/repo).
func (s *Scanner) scanGhqRoot(rootName domain.RootName, rootPath string) ([]domain.Project, error) {
	var projects []domain.Project

	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Skip errors
		}

		if !info.IsDir() {
			return nil
		}

		// Check if this looks like a repository directory
		if s.hasGitDir(path) {
			relPath, _ := filepath.Rel(rootPath, path)
			name := filepath.ToSlash(relPath)

			projectType := domain.ProjectTypeDev
			if rootName == "release" {
				projectType = domain.ProjectTypeRelease
			}

			projects = append(projects, domain.Project{
				Name:   name,
				Root:   rootName,
				Path:   path,
				Type:   projectType,
				HasGit: true,
			})

			return filepath.SkipDir // Don't descend into repositories
		}

		return nil
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
func (s *Scanner) HasGitDir(path string) bool {
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
		if strings.Contains(name, char) {
			return false
		}
	}
	return true
}
