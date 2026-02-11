package config

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/mi8bi/ghqx/internal/domain"
)

func TestFindConfigPathExplicitAndEnv(t *testing.T) {
	tmp, err := ioutil.TempDir("", "ghqx-loader-find")
	if err != nil {
		t.Fatalf("tempdir: %v", err)
	}
	defer os.RemoveAll(tmp)

	// explicit path existing
	f := filepath.Join(tmp, "cfg.toml")
	if err := ioutil.WriteFile(f, []byte("[roots]\n"), 0644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	l := NewLoader()
	p, err := l.findConfigPath(f)
	if err != nil {
		t.Fatalf("findConfigPath explicit failed: %v", err)
	}
	if p != f {
		t.Fatalf("expected %s got %s", f, p)
	}

	// explicit path missing -> error ErrConfigNotFoundAt
	_, err = l.findConfigPath(filepath.Join(tmp, "nope.toml"))
	if err == nil {
		t.Fatalf("expected error for missing explicit path")
	}
	if ge, ok := err.(*domain.GhqxError); ok {
		if ge.Code != domain.ErrCodeConfigNotFound {
			t.Fatalf("unexpected error code: %v", ge.Code)
		}
	}

	// GHQX_CONFIG env
	tmpf, err := ioutil.TempFile("", "ghqx-env-*.toml")
	if err != nil {
		t.Fatalf("tempfile: %v", err)
	}
	tmpf.Close()
	defer os.Remove(tmpf.Name())

	os.Setenv("GHQX_CONFIG", tmpf.Name())
	defer os.Unsetenv("GHQX_CONFIG")

	p2, err := l.findConfigPath("")
	if err != nil {
		t.Fatalf("findConfigPath via GHQX_CONFIG failed: %v", err)
	}
	if p2 != tmpf.Name() {
		t.Fatalf("expected %s got %s", tmpf.Name(), p2)
	}
}

func TestLoadFromPathInvalidTOML(t *testing.T) {
	tmpf, err := ioutil.TempFile("", "ghqx-bad-*.toml")
	if err != nil {
		t.Fatalf("tempfile: %v", err)
	}
	// write invalid toml
	if _, err := tmpf.Write([]byte("not = [")); err != nil {
		t.Fatalf("write: %v", err)
	}
	tmpf.Close()
	defer os.Remove(tmpf.Name())

	l := NewLoader()
	_, err = l.loadFromPath(tmpf.Name())
	if err == nil {
		t.Fatalf("expected error for invalid TOML")
	}
	if ge, ok := err.(*domain.GhqxError); ok {
		if ge.Code != domain.ErrCodeConfigInvalid {
			t.Fatalf("unexpected code: %v", ge.Code)
		}
	}
}

func TestSaveValidateError(t *testing.T) {
	l := NewLoader()
	// invalid cfg (no roots)
	cfg := &Config{Roots: map[string]string{}}
	tmp, err := ioutil.TempDir("", "ghqx-save")
	if err != nil {
		t.Fatalf("tempdir: %v", err)
	}
	defer os.RemoveAll(tmp)

	path := filepath.Join(tmp, "cfg.toml")
	if err := l.Save(cfg, path); err == nil {
		t.Fatalf("expected Save to fail for invalid cfg")
	}
}
