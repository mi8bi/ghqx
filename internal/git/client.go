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
