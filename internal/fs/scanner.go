package fs

import (
	"os"
	"path/filepath"
	"strings"
	"sort"

	"github.com/mi8bi/ghqx/internal/domain"
)

// Scanner handles filesystem operations for project discovery.
type Scanner struct{}

// NewScanner creates a new filesystem scanner.
func NewScanner() *Scanner {
	return &Scanner{}
}

// ScanRoot scans a root directory and returns all projects found.
func (s *Scanner) ScanRoot(rootName domain.RootName, rootPath string) ([]domain.Project, error) {
	if _, err := os.Stat(rootPath); os.IsNotExist(err) {
		return nil, domain.ErrRootDirNotExist(string(rootName), rootPath)
	}

	return s.scanGhqRoot(rootName, rootPath)
}

// scanGhqRoot scans a directory recursively for projects (git repos or plain directories).
func (s *Scanner) scanGhqRoot(rootName domain.RootName, rootPath string) ([]domain.Project, error) {
	var potentialProjects []domain.Project

	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Skip errors (e.g., permission denied)
		}

		if !info.IsDir() {
			return nil // Only interested in directories
		}

		// Don't treat the root itself as a project
		if path == rootPath {
			return nil
		}

		hasGit := s.hasGitDir(path)

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
			Name:        projectName,
			DisplayName: domain.FormatDisplayName(projectName),
			Root:        rootName,
			Path:        path,
			WorkspaceType: domain.DetermineWorkspaceType(rootName), // Updated from Zone and DetermineZone
			Type:        projectType,
			HasGit:      hasGit,
		})

		// If it's a Git repository, skip descending into it
		if hasGit {
			return filepath.SkipDir
		}
		return nil
	})

	if err != nil {
		return nil, domain.ErrFSScanRoot(err)
	}

	// Post-process to filter out intermediate non-Git directories
	var actualProjects []domain.Project

	// Sort potential projects by path length (descending) to easily identify "leaf" projects
	sort.Slice(potentialProjects, func(i, j int) bool {
		return len(potentialProjects[i].Path) > len(potentialProjects[j].Path)
	})

	// Use a map to keep track of paths that are already part of a valid project
	// This helps filter out parent directories that are not projects themselves
	isSubPathOfProject := make(map[string]bool)

	for _, p := range potentialProjects {
		// If this path is already part of a deeper project, skip it
		if isSubPathOfProject[p.Path] {
			continue
		}

		if p.HasGit {
			actualProjects = append(actualProjects, p)
			// Mark all ancestors of this Git repo as "part of a project"
			markAncestors(rootPath, p.Path, isSubPathOfProject)
		} else {
			// For non-Git projects, we only add them if they are the "repo" part
			// of a host/user/repo structure, and not an intermediate host or user directory.
			// This heuristic is based on ghq's cloning structure.
			relPathSegments := strings.Split(p.Name, "/")
			if len(relPathSegments) == 3 { // e.g., github.com/user/repo
				actualProjects = append(actualProjects, p)
				markAncestors(rootPath, p.Path, isSubPathOfProject)
			}
		}
	}
	
	// Reverse the order to get original sorting or sort by DisplayName if desired
	sort.Slice(actualProjects, func(i, j int) bool {
		return actualProjects[i].Name < actualProjects[j].Name // Sort by Name ascending
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