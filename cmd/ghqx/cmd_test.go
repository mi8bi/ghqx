package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mi8bi/ghqx/internal/app"
	"github.com/mi8bi/ghqx/internal/config"
	"github.com/mi8bi/ghqx/internal/status"
)

func TestExtractRepoNameAndContains(t *testing.T) {
	cases := map[string]string{
		"github.com/user/repo": "repo",
		"user/repo":            "repo",
		"simple":               "simple",
	}
	for in, want := range cases {
		if got := extractRepoName(in); got != want {
			t.Fatalf("extractRepoName(%q) = %q, want %q", in, got, want)
		}
	}

	if !contains("a/b/c", "b/c") {
		t.Fatalf("contains should detect suffix with slash")
	}
	if contains("abc", "bc") {
		t.Fatalf("contains should require slash boundary")
	}
}

func TestPadRightAndTruncate(t *testing.T) {
	s := "hello"
	p := padRight(s, 10)
	if len(p) < len(s) {
		t.Fatalf("padRight shortened string")
	}

	long := "abcdefghijklmnopqrstuvwxyz"
	tr := truncateString(long, 5)
	if tr == long {
		t.Fatalf("truncateString did not shorten long string")
	}
}

func TestOutputTablesAndCheckRepositoryExists(t *testing.T) {
	// sample display projects
	projects := []status.ProjectDisplay{
		{Repo: "mi8bi/ghqx", Workspace: "sandbox", GitManaged: "Managed", Status: "clean", FullPath: "/tmp/ghqx"},
	}

	if err := outputCompactTable(projects); err != nil {
		t.Fatalf("outputCompactTable failed: %v", err)
	}
	if err := outputVerboseTable(projects); err != nil {
		t.Fatalf("outputVerboseTable failed: %v", err)
	}

	// Prepare a real app with a temp root so checkRepositoryExists can scan
	tmp, err := os.MkdirTemp("", "ghqx-cmd-check")
	if err != nil {
		t.Fatalf("tempdir: %v", err)
	}
	defer os.RemoveAll(tmp)

	repo := filepath.Join(tmp, "github.com", "user", "repo")
	if err := os.MkdirAll(repo, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	cfg := &config.Config{Roots: map[string]string{"sandbox": tmp}, Default: config.DefaultConfig{Root: "sandbox"}}
	appInstance := app.New(cfg)
	application = appInstance

	// Now check repository exists
	got := checkRepositoryExists("github.com/user/repo")
	if got == "" {
		t.Fatalf("expected checkRepositoryExists to find repo")
	}
}
