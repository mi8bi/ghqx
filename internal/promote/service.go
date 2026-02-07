package promote

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/mi8bi/ghqx/internal/config"
	"github.com/mi8bi/ghqx/internal/domain"
	"github.com/mi8bi/ghqx/internal/fs"
	"github.com/mi8bi/ghqx/internal/git"
)

// Service handles promote and undo operations.
type Service struct {
	cfg         *config.Config
	scanner     *fs.Scanner
	git         *git.Client
	historyPath string
}

// NewService creates a new promote service.
func NewService(cfg *config.Config) *Service {
	home, _ := os.UserHomeDir()
	historyPath := filepath.Join(home, ".config", "ghqx", "history.json")

	return &Service{
		cfg:         cfg,
		scanner:     fs.NewScanner(),
		git:         git.NewClient(),
		historyPath: historyPath,
	}
}

// Options configures promote behavior.
type Options struct {
	ProjectName string
	FromRoot    string
	ToRoot      string
	Force       bool
	DryRun      bool
	AutoGitInit bool
	AutoCommit  bool
}

// Promote moves a project from one root to another.
func (s *Service) Promote(opts Options) (*domain.PromoteRecord, error) {
	// Resolve root paths
	fromPath, exists := s.cfg.GetRoot(opts.FromRoot)
	if !exists {
		return nil, domain.ErrRootNotFound(opts.FromRoot)
	}

	toPath, exists := s.cfg.GetRoot(opts.ToRoot)
	if !exists {
		return nil, domain.ErrRootNotFound(opts.ToRoot)
	}

	// Validate project name
	if !fs.IsSafeName(opts.ProjectName) {
		return nil, domain.ErrProjectNameInvalid
	}

	srcPath := filepath.Join(fromPath, opts.ProjectName)
	if _, err := os.Stat(srcPath); os.IsNotExist(err) {
		return nil, domain.ErrPromoteSourceNotFound(opts.ProjectName)
	}

	// Check if git repo is dirty (unless forced)
	if !opts.Force && s.scanner.HasGitDir(srcPath) {
		dirty, err := s.git.IsDirty(srcPath)
		if err == nil && dirty {
			return nil, domain.ErrGitDirtyRepo
		}
	}

	// Determine destination path
	var dstPath string
	if opts.ToRoot == "sandbox" {
		// Sandbox is flat
		dstPath = filepath.Join(toPath, opts.ProjectName)
	} else {
		// ghq-style: need to preserve host/owner/repo structure
		// For sandbox->dev promotion, we need to create the ghq path
		// This is a simplified version; real implementation might need ghq integration
		dstPath = filepath.Join(toPath, "github.com", "user", opts.ProjectName)
	}

	// Check if destination exists
	if _, err := os.Stat(dstPath); err == nil {
		return nil, domain.ErrPromoteDestExists(dstPath)
	}

	if opts.DryRun {
		return &domain.PromoteRecord{
			ProjectName: opts.ProjectName,
			FromRoot:    domain.RootName(opts.FromRoot),
			FromPath:    srcPath,
			ToRoot:      domain.RootName(opts.ToRoot),
			ToPath:      dstPath,
		}, nil
	}

	// Ensure destination parent directory exists
	if err := s.scanner.EnsureDir(filepath.Dir(dstPath)); err != nil {
		return nil, err
	}

	// Move the directory
	if err := os.Rename(srcPath, dstPath); err != nil {
		return nil, domain.ErrPromoteMoveFailed(err)
	}

	// Auto git init if requested and no .git exists
	if opts.AutoGitInit && !s.scanner.HasGitDir(dstPath) {
		if err := s.git.Init(dstPath); err != nil {
			// Non-fatal, just log
		}
	}

	// Auto commit if requested
	if opts.AutoCommit && s.scanner.HasGitDir(dstPath) {
		msg := "ghqx: promoted from " + opts.FromRoot
		s.git.Commit(dstPath, msg)
	}

	// Record in history
	record := &domain.PromoteRecord{
		Timestamp:   time.Now().Unix(),
		ProjectName: opts.ProjectName,
		FromRoot:    domain.RootName(opts.FromRoot),
		FromPath:    srcPath,
		ToRoot:      domain.RootName(opts.ToRoot),
		ToPath:      dstPath,
	}

	if s.cfg.History.Enabled {
		s.saveHistory(record)
	}

	return record, nil
}

// Undo reverts the last promote operation.
func (s *Service) Undo(dryRun bool) (*domain.PromoteRecord, error) {
	if !s.cfg.History.Enabled {
		return nil, domain.ErrUndoDisabled
	}

	record, err := s.getLastHistory()
	if err != nil {
		return nil, err
	}

	if record == nil {
		return nil, domain.ErrUndoNoHistory
	}

	// Check if destination still exists
	if _, err := os.Stat(record.ToPath); os.IsNotExist(err) {
		return nil, domain.ErrUndoDestMissing
	}

	// Check if source location is now occupied
	if _, err := os.Stat(record.FromPath); err == nil {
		return nil, domain.ErrUndoSourceOccupied
	}

	if dryRun {
		return record, nil
	}

	// Move back
	if err := os.Rename(record.ToPath, record.FromPath); err != nil {
		return nil, domain.ErrPromoteMoveFailed(err)
	}

	// Remove from history
	s.popHistory()

	return record, nil
}

// saveHistory appends a record to the history file.
func (s *Service) saveHistory(record *domain.PromoteRecord) error {
	history, _ := s.loadHistory()
	history = append(history, record)

	// Enforce max limit
	if len(history) > s.cfg.History.Max {
		history = history[len(history)-s.cfg.History.Max:]
	}

	return s.writeHistory(history)
}

// getLastHistory returns the most recent promote record.
func (s *Service) getLastHistory() (*domain.PromoteRecord, error) {
	history, err := s.loadHistory()
	if err != nil {
		return nil, err
	}

	if len(history) == 0 {
		return nil, nil
	}

	return history[len(history)-1], nil
}

// popHistory removes the last record from history.
func (s *Service) popHistory() error {
	history, err := s.loadHistory()
	if err != nil {
		return err
	}

	if len(history) == 0 {
		return nil
	}

	history = history[:len(history)-1]
	return s.writeHistory(history)
}

// loadHistory reads the history file.
func (s *Service) loadHistory() ([]*domain.PromoteRecord, error) {
	data, err := os.ReadFile(s.historyPath)
	if os.IsNotExist(err) {
		return []*domain.PromoteRecord{}, nil
	}
	if err != nil {
		return nil, err
	}

	var history []*domain.PromoteRecord
	if err := json.Unmarshal(data, &history); err != nil {
		return nil, err
	}

	return history, nil
}

// writeHistory writes the history file.
func (s *Service) writeHistory(history []*domain.PromoteRecord) error {
	dir := filepath.Dir(s.historyPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(history, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(s.historyPath, data, 0644)
}
