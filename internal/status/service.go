package status

import (
	"sync"

	"github.com/mi8bi/ghqx/internal/config"
	"github.com/mi8bi/ghqx/internal/domain"
	"github.com/mi8bi/ghqx/internal/fs"
	"github.com/mi8bi/ghqx/internal/git"
)

// Service handles status operations across all workspace roots.
// It provides concurrent scanning of multiple roots and enriches
// project data with git information.
type Service struct {
	cfg     *config.Config
	scanner *fs.Scanner
	git     *git.Client
}

// NewService creates a new status service.
func NewService(cfg *config.Config) *Service {
	return &Service{
		cfg:     cfg,
		scanner: fs.NewScanner(),
		git:     git.NewClient(),
	}
}

// Options configures the behavior of status scanning operations.
type Options struct {
	CheckDirty bool // Whether to check if git repos have uncommitted changes
	LoadBranch bool // Whether to load current branch information
}

// GetAll scans all configured roots and returns all projects.
// If rootFilter is provided, only scans the specified root.
// Scanning is performed concurrently for better performance.
func (s *Service) GetAll(opts Options, rootFilter ...string) ([]domain.Project, error) {
	targetRoots := s.determineTargetRoots(rootFilter)

	var allProjects []domain.Project
	var mu sync.Mutex
	var wg sync.WaitGroup
	var errors []error
	var errMu sync.Mutex

	for name, path := range targetRoots {
		wg.Add(1)
		go s.scanRoot(name, path, opts, &allProjects, &errors, &mu, &errMu, &wg)
	}

	wg.Wait()

	// Return first error if any occurred
	if len(errors) > 0 {
		return nil, errors[0]
	}

	return allProjects, nil
}

// determineTargetRoots returns the roots to scan based on the filter.
func (s *Service) determineTargetRoots(rootFilter []string) map[string]string {
	if len(rootFilter) > 0 && rootFilter[0] != "" {
		rootName := rootFilter[0]
		if rootPath, exists := s.cfg.Roots[rootName]; exists {
			return map[string]string{rootName: rootPath}
		}
	}
	return s.cfg.Roots
}

// scanRoot scans a single root directory for projects.
func (s *Service) scanRoot(
	rootName, rootPath string,
	opts Options,
	allProjects *[]domain.Project,
	errors *[]error,
	mu, errMu *sync.Mutex,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	workspaceType := domain.DetermineWorkspaceType(domain.RootName(rootName))

	projects, err := s.scanner.ScanRoot(domain.RootName(rootName), rootPath)
	if err != nil {
		errMu.Lock()
		*errors = append(*errors, err)
		errMu.Unlock()
		return
	}

	// Set workspace type for all projects
	for i := range projects {
		projects[i].WorkspaceType = workspaceType
	}

	// Enrich projects with git status if requested
	if opts.CheckDirty || opts.LoadBranch {
		s.enrichProjects(projects, opts)
	}

	mu.Lock()
	*allProjects = append(*allProjects, projects...)
	mu.Unlock()
}

// enrichProjects adds git information to all projects that have git.
func (s *Service) enrichProjects(projects []domain.Project, opts Options) {
	for i := range projects {
		if projects[i].HasGit {
			s.enrichProject(&projects[i], opts)
		}
	}
}

// enrichProject adds git information to a single project.
func (s *Service) enrichProject(project *domain.Project, opts Options) {
	if opts.CheckDirty {
		s.updateDirtyStatus(project)
	}

	if opts.LoadBranch {
		s.updateBranchInfo(project)
	}
}

// updateDirtyStatus checks and updates the dirty status of a project.
func (s *Service) updateDirtyStatus(project *domain.Project) {
	dirty, err := s.git.IsDirty(project.Path)
	if err == nil {
		project.Dirty = dirty
		if dirty {
			project.Type = domain.ProjectTypeDirty
		}
	}
}

// updateBranchInfo loads and updates the branch information of a project.
func (s *Service) updateBranchInfo(project *domain.Project) {
	branch, err := s.git.GetBranch(project.Path)
	if err == nil {
		project.Branch = branch
	}
}

// FindProject finds a project by name across all roots.
func (s *Service) FindProject(name string) (*domain.Project, error) {
	projects, err := s.GetAll(Options{})
	if err != nil {
		return nil, err
	}

	for _, p := range projects {
		if p.Name == name {
			return &p, nil
		}
	}

	return nil, domain.ErrProjectNotFound(name)
}
