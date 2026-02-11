package git

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

func TestHasGit(t *testing.T) {
	// HasGit should return true on machines with git installed
	_ = HasGit()
}

func TestInitCommitBranchIsDirty(t *testing.T) {
	if !HasGit() {
		t.Skip("git not available")
	}

	tmp, err := ioutil.TempDir("", "ghqx-git-test")
	if err != nil {
		t.Fatalf("tempdir: %v", err)
	}
	defer os.RemoveAll(tmp)

	c := NewClientWithTimeout(2 * time.Second)

	// init
	if err := c.Init(tmp); err != nil {
		t.Fatalf("git init failed: %v", err)
	}

	// create file and commit
	f := filepath.Join(tmp, "f.txt")
	if err := ioutil.WriteFile(f, []byte("x"), 0644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	// configure git user locally to allow commit
	execCommand := func(name string, args ...string) error {
		cmd := exec.Command(name, args...)
		cmd.Dir = tmp
		return cmd.Run()
	}
	_ = execCommand("git", "config", "user.email", "test@example.com")
	_ = execCommand("git", "config", "user.name", "test")

	// try committing using system git; if it fails, skip remaining assertions
	cmd := exec.Command("git", "add", ".")
	cmd.Dir = tmp
	if err := cmd.Run(); err != nil {
		t.Skipf("git add failed: %v", err)
	}
	cmd = exec.Command("git", "commit", "-m", "msg")
	cmd.Dir = tmp
	if err := cmd.Run(); err != nil {
		t.Skipf("git commit failed: %v", err)
	}

	// branch
	br, err := c.GetBranch(tmp)
	if err != nil {
		t.Fatalf("GetBranch failed: %v", err)
	}
	if br == "" {
		t.Fatalf("expected non-empty branch")
	}

	// IsDirty should be false right after commit
	dirty, err := c.IsDirty(tmp)
	if err != nil {
		t.Fatalf("IsDirty failed: %v", err)
	}
	if dirty {
		t.Fatalf("expected clean repo after commit")
	}
}
