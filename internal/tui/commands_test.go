package tui

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mi8bi/ghqx/internal/app"
	"github.com/mi8bi/ghqx/internal/config"
)

func TestLoadProjects(t *testing.T) {
	tmp, err := os.MkdirTemp("", "ghqx-tui-commands")
	if err != nil {
		t.Fatalf("tempdir: %v", err)
	}
	defer os.RemoveAll(tmp)

	// Create test repository
	repo := filepath.Join(tmp, "github.com", "user", "repo")
	if err := os.MkdirAll(repo, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	cfg := &config.Config{
		Roots:   map[string]string{"sandbox": tmp},
		Default: config.DefaultConfig{Root: "sandbox"},
	}
	appInstance := app.New(cfg)

	model := NewStatusModel(appInstance)

	// Execute loadProjects command
	cmd := model.loadProjects()
	if cmd == nil {
		t.Fatal("loadProjects returned nil command")
	}

	// Execute the command function
	msg := cmd()

	// Check message type
	switch msg := msg.(type) {
	case projectsLoadedMsg:
		if len(msg.projects) == 0 {
			t.Error("expected at least one project")
		}
	case errorMsg:
		t.Errorf("loadProjects failed with error: %v", msg.err)
	default:
		t.Errorf("unexpected message type: %T", msg)
	}
}

func TestLoadProjectsWithError(t *testing.T) {
	// Create config with invalid root path
	cfg := &config.Config{
		Roots:   map[string]string{"sandbox": "/nonexistent/path"},
		Default: config.DefaultConfig{Root: "sandbox"},
	}
	appInstance := app.New(cfg)

	model := NewStatusModel(appInstance)

	// Execute loadProjects command
	cmd := model.loadProjects()
	msg := cmd()

	// Should return errorMsg
	if _, ok := msg.(errorMsg); !ok {
		t.Errorf("expected errorMsg, got %T", msg)
	}
}

func TestProjectsLoadedMsg(t *testing.T) {
	msg := projectsLoadedMsg{
		projects: []ProjectRow{},
	}

	if msg.projects == nil {
		t.Error("projects should not be nil")
	}

	if len(msg.projects) != 0 {
		t.Error("expected empty projects")
	}
}

func TestErrorMsg(t *testing.T) {
	testErr := os.ErrNotExist
	msg := errorMsg{err: testErr}

	if msg.err != testErr {
		t.Error("error mismatch")
	}
}
