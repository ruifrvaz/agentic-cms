# Task Planning

## Active

| ID | Title | Status | Notes |
|----|-------|--------|-------|
| 003 | Register mode for reference files | Not Started | Lightweight `content-import` Register mode: wiki/sources sidecar + index + log for archival/reference files (e.g. invoices), no markdown conversion; start after 002 merges |
| 006 | Shell installer bootstrap (no Go dependency) | PR Open (#6) | `curl \| bash` `install.sh` at repo root, modeled selectively on `smaqit-extensions/install.sh`; adopts ancestor-dir resolution into `agentic-cms update` only — no global-install, `--scope`, `uninstall`, or multi-agent global directory surface. NOTE: its "no global agent/skill directories to seed" design decision is superseded by task 008 — revisit that bullet |
| 008 | Global classification scanner | Not Started | Global Claude Code skill (`~/.claude/skills/agentic-cms-classify/`, installs on first binary run) + optional warning-gated `agentic-cms hooks enable-global` git pre-commit hook via `core.hooksPath`, chaining into per-project `.agentic-cms/hooks/pre-commit` when present. Motivated by a real incident: this repo (never `init`-ed on itself) had zero coverage when a direct git commit landed real PII, caught only by manual review. Depends on task 005 (Completed, v0.5.0) |

## Completed

| ID | Title | Completed | Notes |
|----|-------|-----------|-------|
| 001 | Installer smoke test script | 2026-07-31 | Bash end-to-end smoke test of `agentic-cms init` via `make smoke-test` |
| 002 | Draft content state | 2026-08-15 | PR #2; `status: draft` frontmatter + `docs/<topic>/drafts/` convention, `ac-page --status`/`promote` toolkit support; shipped as part of v0.2.0 |
| 004 | Archive content lifecycle | 2026-08-16 | PR #3; `status: archived` + `docs/<topic>/archive/` mirroring drafts, `ac-page archive`, archived items stay indexed under `## Archived`; renamed content-new-item → content-manage-item; released as v0.3.0 |
| 007 | Rename .agentic-cms/bin to .agentic-cms/scripts | 2026-08-18 | PR #4; `bin/` collided with the generic IDE build-output-exclude convention but held committed source; Tier 1 rename (scaffold source + all living references) via `git mv` + literal-string sweep; released as v0.4.0. Task 005 must reconcile its own ~150 bin/ references against this baseline before completing |
| 005 | First-class content classification (auto CIA rating) | 2026-08-18 | PR #5; single `ac-classify` engine (`classified-hash` staleness + heuristic floors, raise-only) with four thin callers — Claude Code/Codex agent hooks, a git pre-commit gate on staged blobs, skill verify tails as the guaranteed fallback, content-lint's sweep; `ac-page --classification`/`classify`, `ac-inventory` distribution tally; reconciled against task 007's bin→scripts rename; released as v0.5.0. Segregation/vault out of scope |

## Abandoned

| ID | Title | Date | Reason |
|----|-------|------|--------|
