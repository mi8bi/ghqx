package fs

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/mi8bi/ghqx/internal/domain"
)

// Scanner handles filesystem operations for project discovery.
// It traverses directory trees to find git repositories and structured directories.
type Scanner struct{}

// NewScanner creates a new filesystem scanner instance.
func NewScanner() *Scanner {
	return &Scanner{}
}

// ScanRoot scans a root directory and returns all discovered projects.
// Validates the root path exists and recursively searches for projects.
func (s *Scanner) ScanRoot(rootName domain.RootName, rootPath string) ([]domain.Project, error) {
	if _, err := os.Stat(rootPath); os.IsNotExist(err) {
		return nil, domain.ErrRootDirNotExist(string(rootName), rootPath)
	}

	return s.scanGhqRoot(rootName, rootPath)
}

// scanGhqRoot recursively scans a directory for projects (git repos or structured directories).
// Returns a filtered list of actual projects, not intermediate directories.
func (s *Scanner) scanGhqRoot(rootName domain.RootName, rootPath string) ([]domain.Project, error) {
	var potentialProjects []domain.Project

	// Walk directory tree to find projects
	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Skip directories we can't access (permission denied, etc.)
		}

		if !info.IsDir() {
			return nil // Only interested in directories
		}

		// Don't treat the root itself as a project
		if path == rootPath {
			return nil
		}

		hasGit := s.hasGitDir(path)

		// Compute project name relative to root
		relPath, _ := filepath.Rel(rootPath, path)
		projectName := filepath.ToSlash(relPath)
		if projectName == "." || projectName == "" {
			projectName = filepath.Base(path)
		}

		// Determine project type based on root name and git status
		projectType := domain.ProjectTypeDir // Default for non-git directories
		if hasGit {
			switch rootName {
			case "sandbox":
				projectType = domain.ProjectTypeSandboxGit
			case "release":
				projectType = domain.ProjectTypeRelease
			case "dev":
				projectType = domain.ProjectTypeDev
			default:
				projectType = domain.ProjectTypeDir // Fallback
			}
		}

		potentialProjects = append(potentialProjects, domain.Project{
			Name:          projectName,
			DisplayName:   domain.FormatDisplayName(projectName),
			Root:          rootName,
			Path:          path,
			WorkspaceType: domain.DetermineWorkspaceType(rootName),
			Type:          projectType,
			HasGit:        hasGit,
		})

		// If it's a Git repository, skip descending into it to avoid nested checks
		if hasGit {
			return filepath.SkipDir
		}
		return nil
	})

	if err != nil {
		return nil, domain.ErrFSScanRoot(err)
	}

	// Post-process: filter out intermediate directories that aren't actual projects
	var actualProjects []domain.Project

	// Sort by path length descending to identify leaf projects easily
	sort.Slice(potentialProjects, func(i, j int) bool {
		return len(potentialProjects[i].Path) > len(potentialProjects[j].Path)
	})

	// Use a map to track paths that are already part of a valid project
	isSubPathOfProject := make(map[string]bool)

	for _, p := range potentialProjects {
		// Skip if this path is a subpath of an already discovered project
		if isSubPathOfProject[p.Path] {
			continue
		}

		if p.HasGit {
			actualProjects = append(actualProjects, p)
			// Mark all parent directories as part of this project
			markAncestors(rootPath, p.Path, isSubPathOfProject)
		} else {
			// For non-git directories, only include if they form a complete host/user/repo path
			// e.g., github.com/user/repo has 3 components
			relPathSegments := strings.Split(p.Name, "/")
			if len(relPathSegments) == 3 {
				actualProjects = append(actualProjects, p)
				markAncestors(rootPath, p.Path, isSubPathOfProject)
			}
		}
	}

	// Sort final results by name for consistent output
	sort.Slice(actualProjects, func(i, j int) bool {
		return actualProjects[i].Name < actualProjects[j].Name
	})

	return actualProjects, nil
}

// markAncestors marks all parent directories of a project as being part of a project hierarchy.
func markAncestors(rootPath, projectPath string, isSubPathOfProject map[string]bool) {
	currentPath := projectPath
	for currentPath != rootPath && currentPath != filepath.Dir(rootPath) {
		isSubPathOfProject[currentPath] = true
		currentPath = filepath.Dir(currentPath)
	}
}

// markAncestors marks all parent directories of a project as being part of a project hierarchy.
func markAncestors(rootPath, projectPath string, isSubPathOfProject map[string]bool) {
	currentPath := projectPath
	for currentPath != rootPath && currentPath != filepath.Dir(rootPath) {
		isSubPathOfProject[currentPath] = true
		currentPath = filepath.Dir(currentPath)
	}
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