package fs

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCopyFileAndCopyDir(t *testing.T) {
	tmp, err := os.MkdirTemp("", "ghqx-fs-copy")
	if err != nil {
		t.Fatalf("tempdir: %v", err)
	}
	defer os.RemoveAll(tmp)

	// create src structure
	src := filepath.Join(tmp, "src")
	if err := os.MkdirAll(filepath.Join(src, "sub"), 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	f1 := filepath.Join(src, "a.txt")
	if err := os.WriteFile(f1, []byte("hello"), 0644); err != nil {
		t.Fatalf("write f1: %v", err)
	}

	f2 := filepath.Join(src, "sub", "b.txt")
	if err := os.WriteFile(f2, []byte("world"), 0644); err != nil {
		t.Fatalf("write f2: %v", err)
	}

	// test CopyDir
	dst := filepath.Join(tmp, "dst")
	if err := CopyDir(src, dst); err != nil {
		t.Fatalf("CopyDir failed: %v", err)
	}

	// verify files copied
	got1, err := os.ReadFile(filepath.Join(dst, "a.txt"))
	if err != nil {
		t.Fatalf("read dst a: %v", err)
	}
	if string(got1) != "hello" {
		t.Fatalf("a.txt content mismatch: %s", string(got1))
	}

	got2, err := os.ReadFile(filepath.Join(dst, "sub", "b.txt"))
	if err != nil {
		t.Fatalf("read dst b: %v", err)
	}
	if string(got2) != "world" {
		t.Fatalf("b.txt content mismatch: %s", string(got2))
	}

	// test CopyFile by copying a single file
	singleSrc := filepath.Join(tmp, "single.txt")
	if err := os.WriteFile(singleSrc, []byte("single"), 0644); err != nil {
		t.Fatalf("write single: %v", err)
	}
	singleDst := filepath.Join(tmp, "single_dst.txt")
	if err := CopyFile(singleSrc, singleDst); err != nil {
		t.Fatalf("CopyFile failed: %v", err)
	}
	gotS, err := os.ReadFile(singleDst)
	if err != nil {
		t.Fatalf("read single dst: %v", err)
	}
	if string(gotS) != "single" {
		t.Fatalf("single content mismatch: %s", string(gotS))
	}
}
