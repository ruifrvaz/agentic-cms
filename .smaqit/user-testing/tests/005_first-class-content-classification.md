# First-class content classification (auto CIA rating) — E2E Test Playbook

**Test ID:** 005
**Title:** First-class content classification (auto CIA rating)
**Date:** 2026-08-18
**Tester:** User Testing Agent
**Task:** 005

## Objectives

Validate that a real, installed project enforces confidentiality
classification the way the task's core claim promises: one detection
engine (`ac-classify`), triggered at every write path — agent write,
manual edit, and `git commit` — with no path silently letting sensitive
content through unrated. This exercises the actual installed artifacts
against a **brownfield** target (a real project with real files and real
git history), not the Go unit tests' synthetic greenfield `t.TempDir()`,
which check installer mechanics rather than classification behavior.

The install target is a `git clone` of this very repository into a
separate temp directory — not an empty `git init`. This is deliberately
more realistic than a blank sandbox (existing files, existing history,
exercising the brownfield/non-destructive install path a real user with
an existing project would hit) while staying fully isolated: a clone gets
its own independent `.git` at a separate filesystem path, so there is no
risk of git's upward repo-discovery reaching back into this project's own
`.git/hooks/` — the exact class of leak this task's own implementation
hit and fixed (`GIT_CEILING_DIRECTORIES` in `scripts/smoke-test-installer.sh`).
`git clone` (not `cp -r`) also means only committed history is copied —
no local build artifacts or uncommitted working-tree noise ride along.

## Prerequisites

- Go toolchain (to build the CLI) and `git` on `PATH`
- `python3` (the toolkit scripts require it; already a project dependency)
- This repository's primary checkout in a normal, clonable state (any
  uncommitted local changes are irrelevant — `git clone` only copies
  committed history)

## Test Steps

### Step 1 — Build & Unit Test Gate

- [x] `make build` exits 0 (produces `installer/agentic-cms`)
- [x] `make test` exits 0 — runs `go vet ./...` and `go test ./...`, zero failures

### Step 2 — Install & Verify (this CLI's "deploy" is `init` into a target project)

```sh
BIN="$(pwd)/installer/agentic-cms"
D=$(mktemp -d)
git clone --quiet "$(pwd)" "$D"
git -C "$D" config user.email test@test.local
git -C "$D" config user.name test
"$BIN" init "$D"
```

- [x] `git -C "$D" remote -v` shows the clone's origin pointing at this
      repo's local path, **not** the real `origin/agentic-cms.git` on
      GitHub — confirms `$D` is an isolated local clone, not a checkout
      that could accidentally push anywhere
- [x] Output reports `.agentic-cms/scripts/ac-classify` **created** and
      executable — this repo's own root has no `.agentic-cms/` of its own
      (only `scaffold/tree/` carries the embedded templates), so this is
      a fresh create, not a skip
- [x] Output reports `.claude/settings.json`, `.codex/hooks.json`
      **created** for the same reason (no root `.claude/`/`.codex/` exist
      in this repo either)
- [x] Output reports `merged   CLAUDE.md (block appended)` — this repo's
      real root `CLAUDE.md` (`@AGENTS.md`) exists but has no managed
      block yet, so this exercises the genuine brownfield merge path
      (`installClaudeMD` in `scaffold/embed.go`) against real content,
      not the Go test's synthetic "# My project" fixture
- [x] `cat "$D/CLAUDE.md"` shows the original `@AGENTS.md` line preserved,
      followed by the appended `<!-- agentic-cms:begin -->` block
- [x] Output reports `created  git hook: pre-commit (installed)` — a
      fresh clone's `.git/hooks/` only contains git's disabled `.sample`
      files, never an active `pre-commit`, so this is always the
      fresh-install path here, never the merge-append path
- [x] `[ -x "$D/.agentic-cms/scripts/ac-classify" ]` and `[ -x "$D/.git/hooks/pre-commit" ]` both true
- [x] `.git/hooks/pre-commit` in `$D` contains the `agentic-cms:begin` marker
- [x] Nothing under this repo's own working tree or `$(pwd)/.git/hooks/`
      changed: `git status --short` here still shows the same pre-test
      state, and `ls .git/hooks/pre-commit` here still fails (no such file)

### Step 3 — Live Trigger E2E (git pre-commit hook — the event-driven gate this task ships)

Run every command with `cd "$D"` (or `-C "$D"`) from Step 2's install.

- [x] **Turn 1 — clean commit passes.** Create a normal page and commit it:
  ```sh
  .agentic-cms/scripts/ac-page new doc docs/ai/notes.md --title Notes --topic ai --classification C1
  .agentic-cms/scripts/ac-index add topics docs/ai/notes.md "internal notes"
  git add -A && git commit -m "add notes"
  ```
  Expected: commit succeeds (exit 0); no `pre-commit: classification gate BLOCKED` output.
- [x] **Turn 2 — floor violation blocks.** Plant a credential-shaped string and try to commit:
  ```sh
  python3 -c "
  t = open('docs/ai/notes.md').read()
  open('docs/ai/notes.md','w').write(t.replace('## Content', '## Content\n\napi_key: sk-abcdef1234567890abcdef'))
  "
  git add -A && git commit -m "leak a secret"; echo "exit=$?"
  ```
  Expected: commit **fails** (`exit=1`); stderr shows `pre-commit: classification gate BLOCKED this commit` naming `docs/ai/notes.md` and an implied level of `C3`. Verify via `git log --oneline -1` that "leak a secret" is **not** the tip.
- [x] **Turn 3 — raising the rating unblocks it.**
  ```sh
  .agentic-cms/scripts/ac-page classify docs/ai/notes.md C3
  git add -A && git commit -m "leak a secret (rated)"; echo "exit=$?"
  ```
  Expected: `exit=0`; commit succeeds now that the rating matches the floor.
- [x] **Turn 4 — `--no-verify` bypass still works, and lint catches what
  slips through.** Use a **fresh, lower-rated page** for this turn, not
  `docs/ai/notes.md` from Turns 1–3 — that page is already at C3 by now,
  so *any* new content's floor is already covered and `clean` correctly
  stays `true` (staleness alone is advisory, not blocking, by design —
  see CONTENT.md's Classification section). Reusing it here doesn't
  demonstrate the catch; a page whose current rating doesn't already
  cover the new content's floor does:
  ```sh
  .agentic-cms/scripts/ac-page new doc docs/ai/expenses.md --title Expenses --topic ai --classification C1
  git add -A && git commit -m "add expenses page"
  python3 -c "
  t = open('docs/ai/expenses.md').read()
  open('docs/ai/expenses.md','w').write(t + '\ncontact jane@example.com re \$9,000\n')
  "
  git add -A && git commit --no-verify -m "bypass the gate"; echo "exit=$?"
  .agentic-cms/scripts/ac-classify sweep | python3 -c "import json,sys; d=json.load(sys.stdin); print('clean:', d['clean'])"
  ```
  Expected: `--no-verify` commit succeeds (`exit=0`) despite the new PII; the
  follow-up `ac-classify sweep` reports `clean: false`, with a
  `floor_violation: true` / `implied_level: "C2"` entry for
  `docs/ai/expenses.md`, proving the periodic catch-all still sees what
  the gate was bypassed on.

### Step 4 — Additional Validation (agent-hook payload contract)

The agent hooks (`.claude/settings.json`, `.codex/hooks.json`) can't be
triggered by a real Claude Code/Codex session inside this playbook, so
verify the contract directly against `ac-classify hook` using the exact
payload shape those platforms document (`tool_input.file_path`, `cwd`):

```sh
echo "{\"cwd\":\"$D\",\"tool_name\":\"Edit\",\"tool_input\":{\"file_path\":\"docs/ai/expenses.md\"}}" \
  | "$D/.agentic-cms/scripts/ac-classify" hook; echo "exit=$?"
```

- [x] Response is one of: `{"ok":true,"clean":true,...}` (exit 0), an
      auto-raise (`{"ok":true,"auto_raised":true,...}`, exit 0), or a
      block (`{"decision":"block","reason":"..."}` on stdout, same
      reason on stderr, exit 2) — never a silent skip for a real
      `docs/`-path `.md` file
- [x] A non-content path is ignored cleanly:
  ```sh
  echo "{\"cwd\":\"$D\",\"tool_name\":\"Bash\",\"tool_input\":{\"command\":\"ls\"}}" \
    | "$D/.agentic-cms/scripts/ac-classify" hook
  ```
  Expected: `{"ok":true,"skip":true}`, exit 0

## Pass/Fail Criteria

**PASS** — All checkboxes checked. Steps 1–2 commands exit 0. Turn 1 and
Turn 3 commits succeed; Turn 2 is genuinely blocked (verified via
`git log`, not just the printed message); Turn 4's bypass succeeds but
the subsequent sweep reports `clean: false`. Step 4's hook responses match
the documented contract with no silent skip on real content paths.

**FAIL** — Any checkbox unchecked, any unexpected exit code, Turn 2's
blocked commit actually lands in `git log`, or Step 4 silently skips a
real `docs/` path.

## Evidence to Capture

- Full output of `make test`
- Full `agentic-cms init` output from Step 2 (created/merged/skipped lines)
- Turn-by-turn command + output pairs for Step 3 (all four turns) and
  Step 4 (both hook invocations)
- `git log --oneline` from `$D` after Step 3, showing exactly which
  commits landed
- Final `ac-classify sweep` JSON output from Turn 4
