# Always overwrite scaffolding logic files on init/update — E2E Test Playbook

**Test ID:** 010
**Title:** Always overwrite scaffolding logic files on init/update
**Date:** 2026-08-19
**Tester:** User Testing Agent
**Task:** 010

## Objectives

Validate the task's core claim against a real, installed project: on a
second `agentic-cms init`/`update` run, framework-owned scaffolding logic
(`.claude/skills/`, `.claude/agents/`, `.agentic-cms/templates/`,
`.agentic-cms/scripts/`, `.agentic-cms/hooks/`, `.codex/`) is always
overwritten with the shipped version — even if it was locally edited —
while user-owned content (`wiki/`, `CLAUDE.md`'s merge block) is still
never touched. This exercises the actual installed artifacts through two
real `init` invocations against a **brownfield** target (a `git clone` of
this repository, matching this project's established playbook convention —
see `.smaqit/user-testing/tests/005_first-class-content-classification.md`),
not just the Go unit tests' synthetic `t.TempDir()` fixtures.

## Prerequisites

- Go toolchain (to build the CLI) and `git` on `PATH`
- This repository's primary checkout in a normal, clonable state (any
  uncommitted local changes are irrelevant — `git clone` only copies
  committed history)

## Test Steps

### Step 1 — Build & Unit Test Gate

- [x] `make build` exits 0 (produces `installer/agentic-cms`)
- [x] `make test` exits 0 — runs `go vet ./...` and `go test ./...`, zero failures, including `TestInstallOverwritesFrameworkFiles`

### Step 2 — Install & Verify (first init — this CLI's "deploy" is `init` into a target project)

> **Deviation from script (stronger, at user's request):** instead of the
> locally-built `$(pwd)/installer/agentic-cms`, the actual released
> `v0.6.2` binary was installed via the project's public bootstrap —
> `curl -fsSL https://raw.githubusercontent.com/ruifrvaz/agentic-cms/main/install.sh | bash`
> — landing at `~/.local/bin/agentic-cms`. This validates the real
> artifact a user would get, not just a local build.

```sh
BIN="$HOME/.local/bin/agentic-cms"   # v0.6.2, installed via install.sh
D=$(mktemp -d)
git clone --quiet "$(pwd)" "$D"
git -C "$D" config user.email test@test.local
git -C "$D" config user.name test
"$BIN" init "$D" | tee /tmp/first-init-out.txt
```

- [x] Output reports framework-owned paths (`.claude/skills/*/SKILL.md`, `.claude/agents/*.md`, `.agentic-cms/templates/*`, `.agentic-cms/scripts/*`, `.agentic-cms/hooks/pre-commit`, `.codex/hooks.json`) as `created` — this repo's own root carries none of these (only `scaffold/tree/` has the embedded copies), so this is a genuine fresh install
- [x] Output reports `merged   CLAUDE.md (block appended)` — this repo's real root `CLAUDE.md` (`@AGENTS.md`) exists but has no managed block yet, exercising the real brownfield merge path against real content
- [x] Summary line reads `Done: 39 created, 0 updated, 1 merged, 0 skipped.` — nothing to update or skip yet on a target this fresh

### Step 3 — Additional Validation (re-init: overwrite vs. skip split)

Mutate one framework-owned file and one user-content file, then re-run `init` against the same target:

```sh
echo "LOCAL EDIT — should be overwritten" >> "$D/.claude/skills/content-lint/SKILL.md"
echo "LOCAL EDIT — must survive" > "$D/wiki/index.md"
"$BIN" init "$D" | tee /tmp/second-init-out.txt
```

- [x] Output reports `updated  .claude/skills/content-lint/SKILL.md` — **not** `skipped`
- [x] `grep -c "LOCAL EDIT — should be overwritten" "$D/.claude/skills/content-lint/SKILL.md"` returns `0` — the local edit was actually reset to the shipped content, not merely relabeled in the output
- [x] Output reports every other framework-owned path from Step 2 as `updated` too (skills, agents, templates, scripts, hooks, `.codex/`) — 27/27
- [x] Output reports `skipped  wiki/index.md (exists)`
- [x] `cat "$D/wiki/index.md"` still shows exactly `LOCAL EDIT — must survive`
- [x] Output reports `skipped  CLAUDE.md (block already present)` — the merge-block special case still correctly no-ops on a second run, unaffected by this task's change; begin-marker count stayed at 1 (not duplicated)
- [x] Summary line's `updated` count equals the number of `updated` lines printed (27 = 27), and its `created` count is `0` — `Done: 0 created, 27 updated, 0 merged, 13 skipped.`
- [x] Nothing under this repo's own working tree or `$(pwd)/.git/hooks/` changed: `git status --short` here showed no `.claude`/`.agentic-cms`/`.codex`/`CLAUDE.md` entries, and `.git/hooks/pre-commit` still does not exist here

## Pass/Fail Criteria

**PASS** — All checkboxes checked. Steps 1–2 commands exit 0. Step 3's
mutated skill file is genuinely restored (verified via `grep -c`, not just
the printed `updated` line) while `wiki/index.md`'s local edit survives
verbatim and `CLAUDE.md`'s merge block is not duplicated.

**FAIL** — Any checkbox unchecked, any unexpected exit code, the mutated
skill file still contains the local edit after re-init, `wiki/index.md`'s
content changed, or any framework-owned path is reported `skipped` instead
of `updated` on the second run.

## Evidence to Capture

- Full output of `make test`
- Full `agentic-cms init` output from both Step 2 and Step 3 (`/tmp/first-init-out.txt`, `/tmp/second-init-out.txt`)
- `grep -c` result confirming the skill file's local edit was removed
- Final contents of `$D/wiki/index.md`
