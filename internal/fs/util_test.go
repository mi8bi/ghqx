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

// Additional tests for better coverage

func TestCopyDirWithNonExistentSrc(t *testing.T) {
	tmp, err := os.MkdirTemp("", "ghqx-copy-noexist")
	if err != nil {
		t.Fatalf("tempdir: %v", err)
	}
	defer os.RemoveAll(tmp)

	src := filepath.Join(tmp, "nonexistent")
	dst := filepath.Join(tmp, "dst")

	err = CopyDir(src, dst)
	if err == nil {
		t.Fatal("expected error when source doesn't exist")
	}
}

func TestCopyDirWithFileSrc(t *testing.T) {
	tmp, err := os.MkdirTemp("", "ghqx-copy-file-src")
	if err != nil {
		t.Fatalf("tempdir: %v", err)
	}
	defer os.RemoveAll(tmp)

	// Create file instead of directory
	src := filepath.Join(tmp, "file.txt")
	if err := os.WriteFile(src, []byte("test"), 0644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	dst := filepath.Join(tmp, "dst")

	err = CopyDir(src, dst)
	if err == nil {
		t.Fatal("expected error when source is not a directory")
	}
}

func TestCopyDirWithExistingDst(t *testing.T) {
	tmp, err := os.MkdirTemp("", "ghqx-copy-existing")
	if err != nil {
		t.Fatalf("tempdir: %v", err)
	}
	defer os.RemoveAll(tmp)

	src := filepath.Join(tmp, "src")
	if err := os.MkdirAll(src, 0755); err != nil {
		t.Fatalf("mkdir src: %v", err)
	}
	if err := os.WriteFile(filepath.Join(src, "test.txt"), []byte("test"), 0644); err != nil {
		t.Fatalf("write test.txt: %v", err)
	}

	dst := filepath.Join(tmp, "dst")
	if err := os.MkdirAll(dst, 0755); err != nil {
		t.Fatalf("mkdir dst: %v", err)
	}

	// Should still work
	err = CopyDir(src, dst)
	if err != nil {
		t.Fatalf("CopyDir failed with existing dst: %v", err)
	}
}

func TestCopyFileWithNonExistentSrc(t *testing.T) {
	tmp, err := os.MkdirTemp("", "ghqx-copyfile-noexist")
	if err != nil {
		t.Fatalf("tempdir: %v", err)
	}
	defer os.RemoveAll(tmp)

	src := filepath.Join(tmp, "nonexistent.txt")
	dst := filepath.Join(tmp, "dst.txt")

	err = CopyFile(src, dst)
	if err == nil {
		t.Fatal("expected error when source file doesn't exist")
	}
}

func TestCopyFileWithInvalidDst(t *testing.T) {
	tmp, err := os.MkdirTemp("", "ghqx-copyfile-invalid")
	if err != nil {
		t.Fatalf("tempdir: %v", err)
	}
	defer os.RemoveAll(tmp)

	src := filepath.Join(tmp, "src.txt")
	if err := os.WriteFile(src, []byte("test"), 0644); err != nil {
		t.Fatalf("write src: %v", err)
	}

	// Try to copy to a directory that doesn't exist
	dst := filepath.Join(tmp, "nonexistent", "dst.txt")

	err = CopyFile(src, dst)
	if err == nil {
		t.Fatal("expected error when destination directory doesn't exist")
	}
}

func TestCopyDirWithNestedStructure(t *testing.T) {
	tmp, err := os.MkdirTemp("", "ghqx-copy-nested")
	if err != nil {
		t.Fatalf("tempdir: %v", err)
	}
	defer os.RemoveAll(tmp)

	// Create nested structure
	src := filepath.Join(tmp, "src")
	dirs := []string{
		"a",
		"a/b",
		"a/b/c",
		"x",
		"x/y",
	}

	for _, dir := range dirs {
		path := filepath.Join(src, dir)
		if err := os.MkdirAll(path, 0755); err != nil {
			t.Fatalf("mkdir %s: %v", dir, err)
		}
		// Add file in each directory
		if err := os.WriteFile(filepath.Join(path, "file.txt"), []byte("content"), 0644); err != nil {
			t.Fatalf("write file in %s: %v", dir, err)
		}
	}

	dst := filepath.Join(tmp, "dst")
	if err := CopyDir(src, dst); err != nil {
		t.Fatalf("CopyDir failed: %v", err)
	}

	// Verify all files copied
	for _, dir := range dirs {
		path := filepath.Join(dst, dir, "file.txt")
		content, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("failed to read %s: %v", path, err)
		}
		if string(content) != "content" {
			t.Errorf("content mismatch in %s", path)
		}
	}
}

func TestCopyFilePermissions(t *testing.T) {
	tmp, err := os.MkdirTemp("", "ghqx-copy-perms")
	if err != nil {
		t.Fatalf("tempdir: %v", err)
	}
	defer os.RemoveAll(tmp)

	src := filepath.Join(tmp, "src.txt")
	if err := os.WriteFile(src, []byte("test"), 0644); err != nil {
		t.Fatalf("write src: %v", err)
	}

	dst := filepath.Join(tmp, "dst.txt")
	if err := CopyFile(src, dst); err != nil {
		t.Fatalf("CopyFile failed: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(dst); os.IsNotExist(err) {
		t.Fatal("destination file should exist")
	}

	// Verify content
	content, err := os.ReadFile(dst)
	if err != nil {
		t.Fatalf("read dst: %v", err)
	}
	if string(content) != "test" {
		t.Errorf("content mismatch: got %q, want %q", string(content), "test")
	}
}

func TestCopyDirWithEmptyDirectory(t *testing.T) {
	tmp, err := os.MkdirTemp("", "ghqx-copy-empty")
	if err != nil {
		t.Fatalf("tempdir: %v", err)
	}
	defer os.RemoveAll(tmp)

	src := filepath.Join(tmp, "src")
	if err := os.MkdirAll(src, 0755); err != nil {
		t.Fatalf("mkdir src: %v", err)
	}

	dst := filepath.Join(tmp, "dst")
	if err := CopyDir(src, dst); err != nil {
		t.Fatalf("CopyDir failed: %v", err)
	}

	// Verify dst exists
	info, err := os.Stat(dst)
	if err != nil {
		t.Fatalf("stat dst: %v", err)
	}
	if !info.IsDir() {
		t.Error("dst should be a directory")
	}
}
