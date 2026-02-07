package git

import (
	"bytes"
	"context"
	"os/exec"
	"strings"
	"time"

	"github.com/mi8bi/ghqx/internal/domain"
)

// Client handles git operations with timeout support.
type Client struct {
	timeout time.Duration
}

// NewClient creates a new git client with default timeout.
func NewClient() *Client {
	return &Client{
		timeout: 150 * time.Millisecond,
	}
}

// NewClientWithTimeout creates a new git client with custom timeout.
func NewClientWithTimeout(timeout time.Duration) *Client {
	return &Client{
		timeout: timeout,
	}
}

// IsDirty checks if a repository has uncommitted changes.
func (c *Client) IsDirty(repoPath string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "git", "status", "--porcelain")
	cmd.Dir = repoPath

	output, err := cmd.Output()
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return false, domain.ErrGitTimeout("status")
		}
		return false, domain.ErrGitCommandFailed("status", err)
	}

	return len(bytes.TrimSpace(output)) > 0, nil
}

// GetBranch returns the current branch name.
func (c *Client) GetBranch(repoPath string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "git", "rev-parse", "--abbrev-ref", "HEAD")
	cmd.Dir = repoPath

	output, err := cmd.Output()
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", domain.ErrGitTimeout("branch")
		}
		return "", domain.ErrGitCommandFailed("branch", err)
	}

	return strings.TrimSpace(string(output)), nil
}

// ListWorktrees returns all worktrees for a repository.
func (c *Client) ListWorktrees(repoPath string) ([]domain.Worktree, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout*2)
	defer cancel()

	cmd := exec.CommandContext(ctx, "git", "worktree", "list", "--porcelain")
	cmd.Dir = repoPath

	output, err := cmd.Output()
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return nil, domain.ErrGitTimeout("worktree list")
		}
		return nil, domain.ErrGitWorktreeList(err)
	}

	return parseWorktreeList(string(output)), nil
}

// CountWorktrees returns the number of worktrees for a repository.
func (c *Client) CountWorktrees(repoPath string) (int, error) {
	worktrees, err := c.ListWorktrees(repoPath)
	if err != nil {
		return 0, err
	}
	return len(worktrees), nil
}

// parseWorktreeList parses the output of `git worktree list --porcelain`.
func parseWorktreeList(output string) []domain.Worktree {
	var worktrees []domain.Worktree
	var current domain.Worktree

	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			if current.Path != "" {
				worktrees = append(worktrees, current)
				current = domain.Worktree{}
			}
			continue
		}

		parts := strings.SplitN(line, " ", 2)
		if len(parts) < 2 {
			continue
		}

		key := parts[0]
		value := parts[1]

		switch key {
		case "worktree":
			current.Path = value
		case "branch":
			current.Branch = strings.TrimPrefix(value, "refs/heads/")
		case "bare":
			current.Bare = true
		case "locked":
			current.Locked = true
		}
	}

	// Add last worktree if exists
	if current.Path != "" {
		worktrees = append(worktrees, current)
	}

	return worktrees
}

// Init initializes a new git repository.
func (c *Client) Init(repoPath string) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "git", "init")
	cmd.Dir = repoPath

	if err := cmd.Run(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return domain.ErrGitTimeout("init")
		}
		return domain.ErrGitCommandFailed("init", err)
	}

	return nil
}

// Commit creates a commit with the given message.
func (c *Client) Commit(repoPath, message string) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout*2)
	defer cancel()

	// Stage all changes
	addCmd := exec.CommandContext(ctx, "git", "add", ".")
	addCmd.Dir = repoPath
	if err := addCmd.Run(); err != nil {
		return domain.ErrGitCommandFailed("add", err)
	}

	// Commit
	commitCmd := exec.CommandContext(ctx, "git", "commit", "-m", message)
	commitCmd.Dir = repoPath
	if err := commitCmd.Run(); err != nil {
		return domain.ErrGitCommandFailed("commit", err)
	}

	return nil
}

// HasGit checks if git is available in PATH.
func HasGit() bool {
	cmd := exec.Command("git", "--version")
	return cmd.Run() == nil
}
