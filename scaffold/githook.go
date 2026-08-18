package scaffold

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// markers delimiting the block agentic-cms manages inside .git/hooks/pre-commit.
const (
	gitHookMarkerBegin = "# >>> agentic-cms:begin >>>"
	gitHookMarkerEnd   = "# <<< agentic-cms:end <<<"
)

// gitHookBlock is the managed block that dispatches to the versioned gate
// script shipped at .agentic-cms/hooks/pre-commit. It intentionally runs
// relative to the repo root (a pre-commit hook's cwd), never an absolute
// path, so it keeps working across clones and worktrees.
const gitHookBlock = gitHookMarkerBegin + `
# agentic-cms classification gate. If a hook already existed above this
# block and it exits early, this block never runs — merge manually if that
# matters for your setup.
if [ -x .agentic-cms/hooks/pre-commit ]; then
  .agentic-cms/hooks/pre-commit "$@" || exit $?
fi
` + gitHookMarkerEnd + "\n"

// InstallGitHook wires the versioned .agentic-cms/hooks/pre-commit gate into
// dir's actual git hook path. Best-effort and non-destructive:
//   - not a git repository (or git unavailable) -> silent no-op
//   - core.hooksPath is set -> report only, never write into it
//   - no existing pre-commit hook -> install a small caller script
//   - an existing hook without our block -> append it (see caveat above)
//   - already installed -> skip
func InstallGitHook(dir string, res *Result) error {
	if _, err := exec.LookPath("git"); err != nil {
		return nil
	}
	if out, err := exec.Command("git", "-C", dir, "rev-parse", "--is-inside-work-tree").Output(); err != nil || strings.TrimSpace(string(out)) != "true" {
		return nil
	}

	if out, err := exec.Command("git", "-C", dir, "config", "--get", "core.hooksPath").Output(); err == nil {
		if hp := strings.TrimSpace(string(out)); hp != "" {
			res.Skipped = append(res.Skipped,
				"git hook: pre-commit (core.hooksPath="+hp+" — install .agentic-cms/hooks/pre-commit there manually)")
			return nil
		}
	}

	out, err := exec.Command("git", "-C", dir, "rev-parse", "--git-path", "hooks").Output()
	if err != nil {
		return nil // best-effort: an unexpected git state shouldn't fail init
	}
	hooksDir := strings.TrimSpace(string(out))
	if !filepath.IsAbs(hooksDir) {
		hooksDir = filepath.Join(dir, hooksDir)
	}
	if err := os.MkdirAll(hooksDir, 0o755); err != nil {
		return err
	}

	hookPath := filepath.Join(hooksDir, "pre-commit")
	existing, err := os.ReadFile(hookPath)
	if os.IsNotExist(err) {
		content := "#!/usr/bin/env bash\n" + gitHookBlock
		if werr := os.WriteFile(hookPath, []byte(content), 0o755); werr != nil {
			return werr
		}
		res.Created = append(res.Created, "git hook: pre-commit (installed)")
		return nil
	}
	if err != nil {
		return err
	}
	if strings.Contains(string(existing), gitHookMarkerBegin) {
		res.Skipped = append(res.Skipped, "git hook: pre-commit (block already present)")
		return nil
	}
	out2 := strings.TrimRight(string(existing), "\n") + "\n\n" + gitHookBlock
	if werr := os.WriteFile(hookPath, []byte(out2), 0o755); werr != nil {
		return werr
	}
	res.Merged = append(res.Merged, "git hook: pre-commit (block appended)")
	return nil
}
