package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCompareVersions(t *testing.T) {
	cases := []struct {
		a, b string
		want int
	}{
		{"0.1.0", "0.1.0", 0},
		{"0.1.0", "0.2.0", -1},
		{"0.2.0", "0.1.0", 1},
		{"1.0.0", "0.9.9", 1},
		{"1.2", "1.2.0", 0},
		{"1.2.0", "1.2", 0},
		{"1.2.1", "1.2", 1},
		{"0.0.1", "0.0.2", -1},
	}
	for _, c := range cases {
		got, err := compareVersions(c.a, c.b)
		if err != nil {
			t.Errorf("compareVersions(%q, %q) unexpected error: %v", c.a, c.b, err)
			continue
		}
		if got != c.want {
			t.Errorf("compareVersions(%q, %q) = %d, want %d", c.a, c.b, got, c.want)
		}
	}
}

func TestCompareVersionsInvalid(t *testing.T) {
	if _, err := compareVersions("dev", "0.1.0"); err == nil {
		t.Error("compareVersions(\"dev\", \"0.1.0\") expected an error, got nil")
	}
	if _, err := compareVersions("0.1.0", "not-a-version"); err == nil {
		t.Error("compareVersions(\"0.1.0\", \"not-a-version\") expected an error, got nil")
	}
}

func TestResolveDefaultProjectDirFindsAncestor(t *testing.T) {
	root := t.TempDir()
	if err := os.Mkdir(filepath.Join(root, ".agentic-cms"), 0o755); err != nil {
		t.Fatalf("mkdir .agentic-cms: %v", err)
	}
	nested := filepath.Join(root, "docs", "topic")
	if err := os.MkdirAll(nested, 0o755); err != nil {
		t.Fatalf("mkdir nested: %v", err)
	}

	got := resolveDefaultProjectDir(nested)

	wantRoot, err := filepath.EvalSymlinks(root)
	if err != nil {
		t.Fatalf("eval symlinks on root: %v", err)
	}
	gotResolved, err := filepath.EvalSymlinks(got)
	if err != nil {
		t.Fatalf("eval symlinks on result: %v", err)
	}
	if gotResolved != wantRoot {
		t.Errorf("resolveDefaultProjectDir(%q) = %q, want %q", nested, got, root)
	}
}

func TestResolveDefaultProjectDirNoAncestorFallsBackToStart(t *testing.T) {
	start := t.TempDir()

	got := resolveDefaultProjectDir(start)

	wantStart, err := filepath.EvalSymlinks(start)
	if err != nil {
		t.Fatalf("eval symlinks on start: %v", err)
	}
	gotResolved, err := filepath.EvalSymlinks(got)
	if err != nil {
		t.Fatalf("eval symlinks on result: %v", err)
	}
	if gotResolved != wantStart {
		t.Errorf("resolveDefaultProjectDir(%q) = %q, want %q", start, got, start)
	}
}
