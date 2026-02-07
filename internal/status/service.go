package status

import (
	"sync"

	"github.com/mi8bi/ghqx/internal/config"
	"github.com/mi8bi/ghqx/internal/domain"
	"github.com/mi8bi/ghqx/internal/fs"
	"github.com/mi8bi/ghqx/internal/git"
)

// Service handles status operations across all roots.
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

// Options configures status scanning behavior.
type Options struct {
	CheckDirty     bool
	LoadBranch     bool
	CountWorktrees bool
	RootFilter     string
}

// GetAll scans all roots and returns all projects.
func (s *Service) GetAll(opts Options) ([]domain.Project, error) {
	var allProjects []domain.Project
	var mu sync.Mutex
	var wg sync.WaitGroup
	errChan := make(chan error, len(s.cfg.Roots))

	for name, path := range s.cfg.Roots {
		if opts.RootFilter != "" && name != opts.RootFilter {
			continue
		}

		wg.Add(1)
		go func(rootName string, rootPath string) {
			defer wg.Done()

			zone := domain.DetermineZone(domain.RootName(rootName))
			isSandbox := zone == domain.ZoneSandbox

			projects, err := s.scanner.ScanRoot(domain.RootName(rootName), rootPath, isSandbox)
			if err != nil {
				errChan <- err
				return
			}

			// Set zone for all projects
			for i := range projects {
				projects[i].Zone = zone
			}

			// Enrich projects with git status if requested
			if opts.CheckDirty || opts.LoadBranch || opts.CountWorktrees {
				for i := range projects {
					if projects[i].HasGit {
						s.enrichProject(&projects[i], opts)
					}
				}
			}

			mu.Lock()
			allProjects = append(allProjects, projects...)
			mu.Unlock()
		}(name, path)
	}

	wg.Wait()
	close(errChan)

	// Return first error if any occurred
	if err := <-errChan; err != nil {
		return nil, err
	}

	return allProjects, nil
}

// enrichProject adds git information to a project.
func (s *Service) enrichProject(project *domain.Project, opts Options) {
	if opts.CheckDirty && project.HasGit {
		dirty, err := s.git.IsDirty(project.Path)
		if err == nil {
			project.Dirty = dirty
			if dirty {
				project.Type = domain.ProjectTypeDirty
			}
		}
	}

	if opts.LoadBranch && project.HasGit {
		branch, err := s.git.GetBranch(project.Path)
		if err == nil {
			project.Branch = branch
		}
	}

	if opts.CountWorktrees && project.HasGit {
		count, err := s.git.CountWorktrees(project.Path)
		if err == nil {
			project.WorktreeCount = count
		}
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

// GetByRoot returns all projects in a specific root.
func (s *Service) GetByRoot(rootName string) ([]domain.Project, error) {
	return s.GetAll(Options{RootFilter: rootName})
}
