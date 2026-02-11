package status

import (
	"sync"

	"github.com/mi8bi/ghqx/internal/config"
	"github.com/mi8bi/ghqx/internal/domain"
	"github.com/mi8bi/ghqx/internal/fs"
	"github.com/mi8bi/ghqx/internal/git"
)

// Service handles status operations across all workspace roots.
// It coordinates project discovery with git status enrichment and supports
// concurrent scanning for improved performance across multiple roots.
type Service struct {
	// cfg holds the application configuration with root paths
	cfg *config.Config
	// scanner handles filesystem traversal and project discovery
	scanner *fs.Scanner
	// git provides git-related operations (status, branch info, etc.)
	git *git.Client
}

// NewService creates a new status service with default dependencies.
func NewService(cfg *config.Config) *Service {
	return &Service{
		cfg:     cfg,
		scanner: fs.NewScanner(),
		git:     git.NewClient(),
	}
}

// Options configures the behavior of status scanning operations.
type Options struct {
	// CheckDirty determines whether to check for uncommitted changes in git repos
	CheckDirty bool
	// LoadBranch determines whether to load the current branch name for git repos
	LoadBranch bool
}

// GetAll scans all configured roots and returns all discovered projects.
// Optionally filters to a single root using rootFilter.
// Scanning is performed concurrently across all roots for better performance.
func (s *Service) GetAll(opts Options, rootFilter ...string) ([]domain.Project, error) {
	// Determine which roots to scan based on filter
	targetRoots := s.determineTargetRoots(rootFilter)

	// Prepare result containers and synchronization primitives
	var allProjects []domain.Project
	var mu sync.Mutex     // Protects allProjects slice
	var wg sync.WaitGroup // Tracks goroutine completion
	var errors []error
	var errMu sync.Mutex // Protects errors slice

	// Launch concurrent scan for each root
	for name, path := range targetRoots {
		wg.Add(1)
		go s.scanRoot(name, path, opts, &allProjects, &errors, &mu, &errMu, &wg)
	}

	// Wait for all scans to complete
	wg.Wait()

	// Return first error that occurred during scanning
	if len(errors) > 0 {
		return nil, errors[0]
	}

	return allProjects, nil
}

// determineTargetRoots returns the roots to scan based on the filter.
// If a valid root filter is provided, returns only that root.
// Otherwise returns all configured roots.
func (s *Service) determineTargetRoots(rootFilter []string) map[string]string {
	if len(rootFilter) > 0 && rootFilter[0] != "" {
		rootName := rootFilter[0]
		if rootPath, exists := s.cfg.Roots[rootName]; exists {
			return map[string]string{rootName: rootPath}
		}
	}
	// Return all configured roots
	return s.cfg.Roots
}

// scanRoot scans a single root directory for projects.
// This function runs concurrently and safely appends results to shared slices.
func (s *Service) scanRoot(
	rootName, rootPath string,
	opts Options,
	allProjects *[]domain.Project,
	errors *[]error,
	mu, errMu *sync.Mutex,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	// Determine workspace type for all discovered projects in this root
	workspaceType := domain.DetermineWorkspaceType(domain.RootName(rootName))

	// Scan the root directory for projects
	projects, err := s.scanner.ScanRoot(domain.RootName(rootName), rootPath)
	if err != nil {
		// Record error in thread-safe manner
		errMu.Lock()
		*errors = append(*errors, err)
		errMu.Unlock()
		return
	}

	// Set workspace type for all discovered projects
	for i := range projects {
		projects[i].WorkspaceType = workspaceType
	}

	// Enrich projects with git status if requested
	if opts.CheckDirty || opts.LoadBranch {
		s.enrichProjects(projects, opts)
	}

	// Append projects to result in thread-safe manner
	mu.Lock()
	*allProjects = append(*allProjects, projects...)
	mu.Unlock()
}

// enrichProjects adds git information to projects that have git.
func (s *Service) enrichProjects(projects []domain.Project, opts Options) {
	for i := range projects {
		// Only enrich projects that are git repositories
		if projects[i].HasGit {
			s.enrichProject(&projects[i], opts)
		}
	}
}

// enrichProject adds requested git information to a single project.
// Checks dirty status and/or loads branch info based on Options.
func (s *Service) enrichProject(project *domain.Project, opts Options) {
	if opts.CheckDirty {
		s.updateDirtyStatus(project)
	}

	if opts.LoadBranch {
		s.updateBranchInfo(project)
	}
}

// updateDirtyStatus checks for uncommitted changes and updates project status.
func (s *Service) updateDirtyStatus(project *domain.Project) {
	dirty, err := s.git.IsDirty(project.Path)
	if err == nil {
		project.Dirty = dirty
		// Mark as dirty type if repository has changes
		if dirty {
			project.Type = domain.ProjectTypeDirty
		}
	}
}

// updateBranchInfo loads the current branch name for a project.
func (s *Service) updateBranchInfo(project *domain.Project) {
	branch, err := s.git.GetBranch(project.Path)
	if err == nil {
		project.Branch = branch
	}
}

// FindProject searches for a project by its full name across all roots.
// Returns a pointer to the project if found, or an error if not found.
func (s *Service) FindProject(name string) (*domain.Project, error) {
	// Scan all projects without git enrichment for efficiency
	projects, err := s.GetAll(Options{})
	if err != nil {
		return nil, err
	}

	// Search through projects for matching name
	for _, p := range projects {
		if p.Name == name {
			return &p, nil
		}
	}

	// Not found
	return nil, domain.ErrProjectNotFound(name)
}
