package scaffold

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func requireGit(t *testing.T) {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available")
	}
}

func initRepo(t *testing.T, dir string) {
	t.Helper()
	run := func(args ...string) {
		cmd := exec.Command("git", append([]string{"-C", dir}, args...)...)
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v\n%s", args, err, out)
		}
	}
	run("init", "-q")
	run("config", "user.email", "test@test.local")
	run("config", "user.name", "test")
}

func TestInstallGitHookFresh(t *testing.T) {
	requireGit(t)
	dir := t.TempDir()
	initRepo(t, dir)

	res := &Result{}
	if err := InstallGitHook(dir, res); err != nil {
		t.Fatal(err)
	}
	if len(res.Created) != 1 || !strings.Contains(res.Created[0], "installed") {
		t.Fatalf("expected one created git-hook entry, got %v", res.Created)
	}
	hookPath := filepath.Join(dir, ".git", "hooks", "pre-commit")
	fi, err := os.Stat(hookPath)
	if err != nil {
		t.Fatal(err)
	}
	if fi.Mode()&0o111 == 0 {
		t.Error("installed pre-commit hook is not executable")
	}
	b, _ := os.ReadFile(hookPath)
	if !strings.Contains(string(b), gitHookMarkerBegin) {
		t.Error("installed hook missing managed block marker")
	}

	// Idempotent re-run: skip, no duplicate marker.
	res2 := &Result{}
	if err := InstallGitHook(dir, res2); err != nil {
		t.Fatal(err)
	}
	if len(res2.Skipped) != 1 {
		t.Fatalf("expected re-run to skip, got %v", res2)
	}
	b, _ = os.ReadFile(hookPath)
	if strings.Count(string(b), gitHookMarkerBegin) != 1 {
		t.Error("managed block duplicated on re-run")
	}
}

func TestInstallGitHookAppendsToExisting(t *testing.T) {
	requireGit(t)
	dir := t.TempDir()
	initRepo(t, dir)
	hooksDir := filepath.Join(dir, ".git", "hooks")
	if err := os.MkdirAll(hooksDir, 0o755); err != nil {
		t.Fatal(err)
	}
	hookPath := filepath.Join(hooksDir, "pre-commit")
	existing := "#!/usr/bin/env bash\necho existing\n"
	if err := os.WriteFile(hookPath, []byte(existing), 0o755); err != nil {
		t.Fatal(err)
	}

	res := &Result{}
	if err := InstallGitHook(dir, res); err != nil {
		t.Fatal(err)
	}
	if len(res.Merged) != 1 {
		t.Fatalf("expected one merged git-hook entry, got %v", res)
	}
	b, _ := os.ReadFile(hookPath)
	s := string(b)
	if !strings.HasPrefix(s, "#!/usr/bin/env bash\necho existing") {
		t.Error("existing hook content lost")
	}
	if !strings.Contains(s, gitHookMarkerBegin) {
		t.Error("managed block not appended")
	}
	fi, err := os.Stat(hookPath)
	if err != nil {
		t.Fatal(err)
	}
	if fi.Mode()&0o111 == 0 {
		t.Error("merged hook is not executable")
	}
}

func TestInstallGitHookRespectsHooksPath(t *testing.T) {
	requireGit(t)
	dir := t.TempDir()
	initRepo(t, dir)
	if err := os.MkdirAll(filepath.Join(dir, "custom-hooks"), 0o755); err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command("git", "-C", dir, "config", "core.hooksPath", "custom-hooks")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git config: %v\n%s", err, out)
	}

	res := &Result{}
	if err := InstallGitHook(dir, res); err != nil {
		t.Fatal(err)
	}
	if len(res.Created) != 0 || len(res.Merged) != 0 {
		t.Fatalf("expected no writes under core.hooksPath, got %v", res)
	}
	if len(res.Skipped) != 1 || !strings.Contains(res.Skipped[0], "core.hooksPath") {
		t.Fatalf("expected a core.hooksPath report, got %v", res.Skipped)
	}
	if entries, _ := os.ReadDir(filepath.Join(dir, "custom-hooks")); len(entries) != 0 {
		t.Error("wrote into the custom hooks path — must only report")
	}
	if _, err := os.Stat(filepath.Join(dir, ".git", "hooks", "pre-commit")); !os.IsNotExist(err) {
		t.Error("wrote a fallback into .git/hooks despite core.hooksPath being set")
	}
}

func TestInstallGitHookNonGitDirNoOp(t *testing.T) {
	dir := t.TempDir()
	res := &Result{}
	if err := InstallGitHook(dir, res); err != nil {
		t.Fatal(err)
	}
	if len(res.Created) != 0 || len(res.Merged) != 0 || len(res.Skipped) != 0 {
		t.Errorf("expected a silent no-op for a non-git directory, got %v", res)
	}
}
