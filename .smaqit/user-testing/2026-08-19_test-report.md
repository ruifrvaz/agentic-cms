# User Testing Report

**Date:** 2026-08-19
**Repository:** ruifrvaz/agentic-cms
**Branch:** main
**Commit:** 6092345
**OS/Arch:** linux/x86_64
**Duration:** ~10 minutes

## Scope
- Test file: `.smaqit/user-testing/tests/010_always-overwrite-scaffolding-logic-files-on-init-update.md`
- Commands executed:
   - `make build`
   - `make test`
   - `curl -fsSL https://raw.githubusercontent.com/ruifrvaz/agentic-cms/main/install.sh | bash` (installs the real released binary — a deliberate upgrade over the playbook's local-build step, at the user's request)
   - `git clone` of this repo into a temp dir, then `agentic-cms init` (twice, with a mutation between runs)
   - `grep -c`, `cat`, `git status --short` verification commands

## Checklist
- [x] Test command discovered and confirmed (`Makefile`: `build`, `test`)
- [x] Dependencies installed (Go toolchain, git — already present)
- [x] Test suite executed (`make test`, plus two live `agentic-cms init` runs)
- [x] Results captured (pass/fail + key errors) — none; all checks passed
- [x] Evidence collected (per test file 010, all 12 checkboxes verified and marked)

## Execution Log (Timestamped)
- `make build` — built `installer/agentic-cms` from source, exit 0
- `make test` — `go vet ./...` + `go test ./...`, both packages `ok`, including the new `TestInstallOverwritesFrameworkFiles`
- Installed the actual **released v0.6.2** binary via `install.sh` to `~/.local/bin/agentic-cms`, replacing the plan to use the local build — verified `agentic-cms --version` reports `v0.6.2`
- Cloned this repo into a fresh temp dir, ran `agentic-cms init` — `Done: 39 created, 0 updated, 1 merged, 0 skipped.`
- Hand-edited `.claude/skills/content-lint/SKILL.md` (appended a marker line) and overwrote `wiki/index.md` with local content
- Re-ran `agentic-cms init` against the same target — `Done: 0 created, 27 updated, 0 merged, 13 skipped.`
- Verified `content-lint/SKILL.md`'s marker line is gone (`grep -c` → `0`) — genuinely reset, not just relabeled
- Verified `wiki/index.md` still reads the hand-written content
- Verified `CLAUDE.md`'s managed block was not duplicated (1 begin-marker)
- Verified this repo's own working tree and `.git/hooks/` were untouched by either run
- Cleaned up the temp target directory

## Results
- Overall: **PASS**
- Summary:
   - 12/12 playbook checkboxes passed, 0 failed
   - Confirmed against the actual GitHub release artifact (v0.6.2), not just a local build — closes the loop on the original user report (skills silently going stale on `agentic-cms update`/`init`)
   - Framework-owned scaffolding (skills, agents, templates, scripts, hooks, `.codex/`) is now unconditionally refreshed on re-init; user-owned content (`wiki/`, `CLAUDE.md`'s merge state) remains untouched
   - Reported counts are internally consistent: `updated` line count (27) matches the summary's `updated` field exactly

## Pain Points
- Blockers:
   - None
- Issues:
   - None
- UX Friction:
   - None
- Performance:
   - None

## Recommendations
- None — the fix behaves exactly as specified and matches the release artifact. No further action needed for this task.
