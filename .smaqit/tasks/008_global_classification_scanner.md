---
status: Not Started
created: "2026-08-18"
---

# Global classification scanner

## Description

Task 005 built content classification (`ac-classify`) as a strictly
project-scoped system: every enforcement layer — agent hooks, the git
pre-commit gate, skill verify tails, content-lint's sweep — only exists in
a project after `agentic-cms init` has been explicitly run there. Nothing
protects a project that hasn't opted in, including **agentic-cms's own
source repository**, which has never run `init` on itself.

This gap is not hypothetical: during manual testing of task 005, a real
file containing a name, SSN, bank account number, and medical diagnosis
was committed directly to this repo's `main` branch (`scaffold/tree/docs/readme.md`,
commit `0c87864`) with nothing to stop it — no agent hook (none installed),
no pre-commit gate (none installed), nothing. It was caught only because
the user noticed by hand and asked why nothing stopped them. The commit
was unpushed and has since been removed via `git reset --hard`, so there
was no public exposure — but the near-miss is the direct motivation for
this task. Worse, the file's location (`scaffold/tree/`, which
`//go:embed all:tree` bakes into the compiled binary) meant that had it
shipped in a release, the PII would have been written into every project
that ran `agentic-cms init` from that point forward.

This task adds a **global** layer on top of task 005's project-scoped one,
with two genuinely different pieces because they solve different halves
of the problem:

1. **A global Claude Code skill**, installed to `~/.claude/skills/` the
   first time the `agentic-cms` binary is run (any subcommand) — available
   to an agent working in *any* project on the device, not just ones
   running `agentic-cms init`. This only helps when an agent is doing the
   writing.
2. **An optional, explicitly opt-in global git pre-commit hook** — the
   piece that actually would have caught the incident above, since that
   commit was made directly by the user in a terminal with no agent
   involved. Skills only fire inside an agent turn; only a git-level hook
   fires on every commit regardless of who or what made it.

## Issue Triage Context

**Mode:** Auto
**Technologies:** Claude Code, Git
**Platforms/Environments:** Linux
**Features/Integrations:** global skill installation, git core.hooksPath
**Versions/Constraints:** Depends on task 005 (Completed, v0.5.0) for `ac-classify`'s floor-pattern definitions and the `.agentic-cms/hooks/pre-commit` convention this task's global hook chains into

## Design Decisions

- **Two pieces, not one, because they cover different gaps.** The skill
  covers agent-mediated writes in non-agentic-cms projects. The hook
  covers *any* commit, agent or human, in *any* repo — including
  already-existing ones. Building only the skill would not have caught
  the motivating incident; that requires the hook.
- **Skill installs on first binary run, not a separate command** (user
  decision) — the same on-demand pattern already used for pandoc/
  markitdown/python-pptx (README's "installed on demand" design note):
  `main()` checks whether `~/.claude/skills/agentic-cms-classify/` exists
  before dispatching to any subcommand, and installs it if missing.
  `go install` alone can't do this (Go's installer has no post-install
  hook), so first-run detection is the only mechanism that works
  uniformly for `go install`, `make install`, and the task 006 shell
  bootstrap alike. Print one informational line the first time only
  ("installed the global classify skill to ~/.claude/skills/agentic-cms-classify/");
  silent on every run after.
- **The git hook is opt-in with an explicit warning, never automatic**
  (user decision). A new `agentic-cms hooks enable-global` subcommand
  prints the exact consequence before acting — "this sets git's
  core.hooksPath globally; EVERY repo on this machine will use the global
  hook instead of its own .git/hooks/, though repos with their own
  .agentic-cms/hooks/pre-commit keep running it via chaining (see below)"
  — and requires explicit confirmation (interactive y/N, or `--yes` for
  scripted use). `agentic-cms hooks disable-global` reverses it by
  unsetting `core.hooksPath`.
- **`core.hooksPath`, not `init.templateDir`** (user decision, after
  trade-off discussion). `core.hooksPath` is the only mechanism that
  retroactively covers already-existing repos — including this one, which
  is exactly the repo the incident happened in. `init.templateDir` only
  seeds *newly created or cloned* repos going forward and would not have
  helped here.
- **Explicit chaining, always built in, not optional.** `core.hooksPath`
  makes git ignore every repo's own `.git/hooks/pre-commit` entirely —
  without chaining, turning this on would silently downgrade every
  already-installed agentic-cms project from task 005's rating-aware gate
  to the simpler global scan. The global hook script therefore always
  checks for `<repo-root>/.agentic-cms/hooks/pre-commit` (the tracked,
  versioned script content, not the bypassed `.git/hooks/` entry point)
  and runs it if present; falls back to the generic global scan otherwise.
- **The global scan is floor-pattern-only — no rating, no
  `classification:` frontmatter, no `classified-hash`.** Most projects a
  global hook touches will never have run `agentic-cms init` and have no
  `CONTENT.md` schema at all. The global scanner (`ac-scan`) can only
  reasonably do what needs no per-file state: pattern-match staged content
  for credential-shaped strings (→ block) and PII/currency shapes (→ warn
  or block, same severity split as the project-scoped gate). No
  `classified-hash`, no auto-raise, no ratchet — those concepts require
  frontmatter that won't exist on most targets.
- **Pattern sync between `ac-classify` and `ac-scan`, enforced by test, not
  shared code.** `ac-classify` is deployed standalone into arbitrary
  project directories with no dependency guarantees (no PYTHONPATH, no
  package registry) — it can't import a shared module at runtime. Rather
  than build a templating system to generate two scripts from one source
  (real complexity for a two-script problem), keep `C2_PATTERNS`/
  `C3_PATTERNS` literally duplicated between `scaffold/tree/.agentic-cms/scripts/ac-classify`
  and the new global `ac-scan`, and add a Go test that extracts both
  pattern blocks and asserts byte-for-byte equality — CI fails the moment
  they drift, which is the same "no integration point re-implements
  detection" guardrail from task 005, adapted for a case where the two
  copies genuinely can't share a runtime import.
- **Global hook files live at `~/.agentic-cms/global-hooks/`** — mirrors
  the per-project `.agentic-cms/hooks/` convention, discoverable, and
  distinct from `~/.claude/skills/` (the skill and the hook are installed
  and activated independently; a user can have one without the other).
- **Cross-task consistency note, not a code change here:** task 006
  (Shell installer bootstrap, in progress as of this task's creation)
  currently states as a design decision: *"No `--install-global` step —
  agentic-cms's project-scoped `init` already covers that need; there are
  no global agent/skill directories to seed."* That was accurate when
  written, before this task existed. This task makes it stale. Not edited
  here since 006 is another session's active work in a different
  worktree — flagged so whoever completes 006 can revisit that bullet.

## Implementation Steps

1. **Global skill content** — new `scaffold/global/skills/agentic-cms-classify/SKILL.md`:
   agent-facing instructions describing when to run a classification scan
   (before writing sensitive-looking content, on request, or periodically)
   and how to invoke the bundled `ac-scan` script. Bundle
   `scaffold/global/skills/agentic-cms-classify/ac-scan` alongside it — a
   new script, structurally mirroring `ac-classify`'s `check`/`sweep`
   commands but stripped of frontmatter/rating logic: `ac-scan check
   <path...>` and `ac-scan sweep [dir]` report floor-pattern hits only
   (no `classification`, `unrated`, `stale`, `classified-hash` fields —
   just `path`, `floor`, `hit` boolean).
2. **First-run installer** — in `main.go`, before dispatching to any
   subcommand: check whether `~/.claude/skills/agentic-cms-classify/`
   exists; if not, copy the embedded `scaffold/global/skills/` content
   there (non-destructive — skip silently if already present, matching
   `Install()`'s existing skip-if-exists philosophy) and print one
   informational line. Requires a new `//go:embed` directive for
   `scaffold/global/` alongside the existing `scaffold/tree/` embed in
   `scaffold/embed.go` (or a sibling file), plus a `$HOME`-resolution
   helper (`os.UserHomeDir()`).
3. **Global git hook script** — new `scaffold/global/hooks/pre-commit`:
   collects staged files across the repo it's running in (same
   `git diff --cached` + `git checkout-index` snapshot approach as the
   per-project gate), checks for `<repo-root>/.agentic-cms/hooks/pre-commit`
   first and executes it (chaining) if present; otherwise runs `ac-scan
   sweep` (bundled alongside this script, same duplicated-but-tested
   pattern list as step 1) against the staged snapshot and blocks/warns
   using the same block-vs-warn split as the project-scoped gate.
4. **`agentic-cms hooks` subcommand** — new `hooks enable-global` and
   `hooks disable-global` in `main.go`: `enable-global` prints the exact
   warning text from the Design Decisions above, requires confirmation
   (interactive prompt, or `--yes`), then writes the embedded
   `scaffold/global/hooks/` content to `~/.agentic-cms/global-hooks/` and
   runs `git config --global core.hooksPath ~/.agentic-cms/global-hooks/`.
   `disable-global` runs `git config --global --unset core.hooksPath`
   (leaves the files in place, inert, so re-enabling is instant).
5. **Pattern-sync test** — new Go test (`scaffold/embed_test.go` or
   similar) that reads both `scaffold/tree/.agentic-cms/scripts/ac-classify`
   and `scaffold/global/hooks/ac-scan` (or wherever the patterns end up
   living), extracts the `C2_PATTERNS`/`C3_PATTERNS` blocks via regex, and
   fails if they differ.
6. **Documentation** — root `README.md`: new section explaining the global
   layer, how it differs from project-scoped `init`, and the explicit
   opt-in nature of the hook; `CONTENT.md` is NOT touched (it documents
   the project-scoped schema; the global layer has no frontmatter schema
   of its own).
7. **Verification**: `go vet ./...`, `go test ./...` (including the new
   pattern-sync test), `make smoke-test`. Manual pass: run the built
   binary once in a scratch `$HOME` (or with `HOME` overridden) and
   confirm the skill lands at `~/.claude/skills/agentic-cms-classify/`
   with the informational line printed exactly once; run it a second time
   and confirm silence. In a throwaway repo with no `.agentic-cms/`, run
   `agentic-cms hooks enable-global`, confirm the warning and confirmation
   gate, stage a credential-shaped file, confirm the global hook blocks
   the commit; in a *different* throwaway repo that *does* have task
   005's per-project gate installed, confirm the global hook chains into
   it (its richer stale/floor/bleed messages appear, not the generic
   scan's simpler output) rather than replacing it. Run `hooks
   disable-global` and confirm `core.hooksPath` is unset and an untouched
   repo's own local hook (if any) resumes working normally.

## Known Issues Triage

[Populated by smaqit.task-start via smaqit.utils.triage-issues. Do not edit manually.]

## Acceptance Criteria

- [ ] The global skill (`~/.claude/skills/agentic-cms-classify/`, with
      bundled `ac-scan`) installs automatically the first time the
      `agentic-cms` binary runs any subcommand, printing one informational
      line; subsequent runs are silent; installation is non-destructive
      (never overwrites an existing global skill directory)
- [ ] `ac-scan check`/`sweep` detects floor-pattern hits (credential-shaped
      → block-equivalent, PII/currency → warn-equivalent) with no
      dependency on `classification:`/`classified-hash` frontmatter
- [ ] `agentic-cms hooks enable-global` prints an explicit warning
      describing the `core.hooksPath` consequence, requires confirmation
      (interactive or `--yes`), then wires the global hook; `hooks
      disable-global` cleanly reverses it
- [ ] The global hook, once enabled, blocks a floor-violating commit in a
      repo with no project-scoped gate of its own
- [ ] The global hook, once enabled, chains into and preserves the richer
      per-project gate in a repo that already has
      `.agentic-cms/hooks/pre-commit` — verified by observing that
      project's specific (stale/floor/bleed) messages, not the generic
      scan's output
- [ ] `C2_PATTERNS`/`C3_PATTERNS` in `ac-classify` and `ac-scan` are
      identical, enforced by an automated test that fails on drift
- [ ] `go vet ./...`, `go test ./...`, and `make smoke-test` all pass
- [ ] Running this repo's own build with `hooks enable-global` active
      would have blocked the exact commit that motivated this task
      (verified by reconstructing an equivalent floor-violating commit in
      a scratch clone of this repo, per the manual verification in
      Implementation Step 7)

## Findings

[Populated by smaqit.task-complete. Do not fill in manually before task is complete.]

**Implementation approach:**
- TBD

**Decisions made:**
- TBD

**Blockers encountered:**
- TBD

**Follow-up identified:**
- TBD

## Files to Create / Modify

| File | Action |
|------|--------|
| scaffold/global/skills/agentic-cms-classify/SKILL.md | Create |
| scaffold/global/skills/agentic-cms-classify/ac-scan | Create |
| scaffold/global/hooks/pre-commit | Create |
| scaffold/embed.go (or a sibling file) | Modify — new embed + first-run install logic |
| main.go | Modify — first-run check, `hooks enable-global`/`disable-global` subcommands |
| scaffold/embed_test.go (or similar) | Create — pattern-sync test |
| README.md | Modify — document the global layer |

## Notes

Scope watch: this is a genuinely large task (new embed tree, new Go
subcommand surface, a second scanner script, a global-config-mutating
command with its own confirmation UX). If it balloons, the natural split
is the global skill (steps 1–2, low risk, no git config mutation) as this
task, with the global hook (steps 3–4) as a follow-up — decide at
task-start, not now, same as task 005's own scope-watch note.

Cross-reference: task 006 (in progress) asserts "there are no global
agent/skill directories to seed" — see Design Decisions above. Not
resolved here; flagged for that task's owner.

The incident that motivated this task: a highly sensitive file (name,
SSN, bank account, medical diagnosis) was committed directly to this
repo's `main` at `scaffold/tree/docs/readme.md` (commit `0c87864`, since
removed via `git reset --hard`, never pushed). This repo had zero
classification coverage because it has never run `agentic-cms init` on
itself — a separate, smaller possibility (not this task's scope) is
simply running `init` on this repo directly, which would give it task
005's project-scoped gate immediately, independent of whether the global
layer in this task also ships.

Child tasks inherit their active parent's branch, worktree, and workflow
mode. Only a standalone or parent task owns Git lifecycle cleanup.
