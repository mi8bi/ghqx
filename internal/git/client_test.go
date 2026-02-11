package git

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestHasGitUnavailableWhenPathEmpty(t *testing.T) {
	origPath := os.Getenv("PATH")
	os.Setenv("PATH", "")
	defer os.Setenv("PATH", origPath)
	// Verify git cannot be executed when PATH is empty
	cmd := exec.Command("git", "--version")
	if err := cmd.Run(); err == nil {
		t.Fatalf("expected git --version to fail when PATH is empty")
	}
}

func TestIsDirtyAndGetBranchNonRepo(t *testing.T) {
	tmp, err := ioutil.TempDir("", "ghqx-git-nonrepo")
	if err != nil {
		t.Fatalf("tempdir: %v", err)
	}
	defer os.RemoveAll(tmp)

	// Skip if git not available
	if exec.Command("git", "--version").Run() != nil {
		t.Skip("git not available on this system")
	}

	c := NewClient()

	dirty, err := c.IsDirty(tmp)
	if err == nil {
		if dirty {
			t.Fatalf("expected non-repo to be not dirty")
		}
	}

	if _, err := c.GetBranch(tmp); err == nil {
		t.Fatalf("expected GetBranch to fail on non-git directory")
	}
}

func TestInitCommitFlowWhenGitAvailable(t *testing.T) {
	// Skip if git not available
	if exec.Command("git", "--version").Run() != nil {
		t.Skip("git not available on this system")
	}

	tmp, err := ioutil.TempDir("", "ghqx-git-repo")
	if err != nil {
		t.Fatalf("tempdir: %v", err)
	}
	defer os.RemoveAll(tmp)

	// Initialize a real git repo using git CLI
	cmd := exec.Command("git", "init")
	cmd.Dir = tmp
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git init failed: %v, out: %s", err, string(out))
	}

	// configure local user to allow commits
	cmd = exec.Command("git", "config", "user.email", "test@example.com")
	cmd.Dir = tmp
	_ = cmd.Run()
	cmd = exec.Command("git", "config", "user.name", "Test")
	cmd.Dir = tmp
	_ = cmd.Run()
	// disable GPG signing for tests (CI may have gpg configured)
	cmd = exec.Command("git", "config", "commit.gpgsign", "false")
	cmd.Dir = tmp
	_ = cmd.Run()

	// create a file and commit using git CLI
	f := filepath.Join(tmp, "file.txt")
	if err := ioutil.WriteFile(f, []byte("x"), 0644); err != nil {
		t.Fatalf("write file: %v", err)
	}
	cmd = exec.Command("git", "add", ".")
	cmd.Dir = tmp
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git add failed: %v, out: %s", err, string(out))
	}
	cmd = exec.Command("git", "commit", "-m", "msg")
	cmd.Dir = tmp
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git commit failed: %v, out: %s", err, string(out))
	}

	c := NewClient()

	dirty, err := c.IsDirty(tmp)
	if err != nil {
		t.Fatalf("IsDirty error: %v", err)
	}
	if dirty {
		t.Fatalf("expected repo to be clean after commit")
	}

	branch, err := c.GetBranch(tmp)
	if err != nil {
		t.Fatalf("GetBranch error: %v", err)
	}
	if branch == "" {
		t.Fatalf("expected non-empty branch name")
	}
}
