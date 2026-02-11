package git

import (
	"bytes"
	"context"
	"os/exec"
	"strings"
	"time"

	"github.com/mi8bi/ghqx/internal/domain"
)

// Client handles git operations with configurable timeout support.
// Timeouts prevent hanging when git commands are slow or unresponsive.
type Client struct {
	// timeout defines the maximum duration for git commands
	timeout time.Duration
}

// NewClient creates a new git client with a default 150ms timeout.
func NewClient() *Client {
	return &Client{
		timeout: 150 * time.Millisecond,
	}
}

// NewClientWithTimeout creates a new git client with a custom timeout duration.
func NewClientWithTimeout(timeout time.Duration) *Client {
	return &Client{
		timeout: timeout,
	}
}

// IsDirty checks if a repository has uncommitted changes.
// Returns true if there are staged or unstaged modifications.
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

	// Non-empty output means there are changes
	return len(bytes.TrimSpace(output)) > 0, nil
}

// GetBranch returns the current branch name of a repository.
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

// Init initializes a new git repository in the given directory.
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

// Commit stages all changes and creates a commit with the given message.
func (c *Client) Commit(repoPath, message string) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout*2)
	defer cancel()

	// Stage all changes
	addCmd := exec.CommandContext(ctx, "git", "add", ".")
	addCmd.Dir = repoPath
	if err := addCmd.Run(); err != nil {
		return domain.ErrGitCommandFailed("add", err)
	}

	// Create commit
	commitCmd := exec.CommandContext(ctx, "git", "commit", "-m", message)
	commitCmd.Dir = repoPath
	if err := commitCmd.Run(); err != nil {
		return domain.ErrGitCommandFailed("commit", err)
	}

	return nil
}

// HasGit checks if the git command is available in the system PATH.
func HasGit() bool {
	cmd := exec.Command("git", "--version")
	return cmd.Run() == nil
}
