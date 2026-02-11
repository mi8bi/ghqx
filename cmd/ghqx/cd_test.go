package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/mi8bi/ghqx/internal/app"
	"github.com/mi8bi/ghqx/internal/config"
)

func TestLoadProjectsForSelection(t *testing.T) {
	tmp, err := ioutil.TempDir("", "ghqx-cd-load")
	if err != nil {
		t.Fatalf("tempdir: %v", err)
	}
	defer os.RemoveAll(tmp)

	// create structure github.com/user/repo
	repo := filepath.Join(tmp, "github.com", "user", "repo")
	if err := os.MkdirAll(repo, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	cfg := &config.Config{Roots: map[string]string{"sandbox": tmp}, Default: config.DefaultConfig{Root: "sandbox"}}
	appInstance := app.New(cfg)
	application = appInstance

	projects, err := loadProjectsForSelection()
	if err != nil {
		t.Fatalf("loadProjectsForSelection failed: %v", err)
	}
	if len(projects) == 0 {
		t.Fatalf("expected at least one project")
	}
	if projects[0].Repo == "" {
		t.Fatalf("expected project to have Repo field set")
	}
}
