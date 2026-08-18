# Task Planning

## Active

| ID | Title | Status | Notes |
|----|-------|--------|-------|
| 003 | Register mode for reference files | Not Started | Lightweight `content-import` Register mode: wiki/sources sidecar + index + log for archival/reference files (e.g. invoices), no markdown conversion; start after 002 merges |
| 005 | First-class content classification (auto CIA rating) | In Progress | `classification: C0–C3` frontmatter auto-assigned by agent at write time against a CONTENT.md rubric; toolkit validation (`ac-page --classification`/`classify`, `ac-inventory` reporting); lint bleed/ratchet enforcement. Refined 2026-08-18, then reviewed for SRP/KISS same day: single `ac-classify` engine (`classified-hash` staleness + heuristic floors, raise-only) with four thin callers — post-tool-use hooks (Claude Code + Codex only; Copilot deferred, no blocking capability), a git pre-commit gate on staged blobs, skill verify tails as the guaranteed fallback, and content-lint's existing audit sweep gaining a category. No caller may duplicate detection logic (explicit acceptance criterion). 004 dependency satisfied (merged, v0.3.0). Segregation/vault out of scope — follow-up task once downstream pattern settles |
| 006 | Shell installer bootstrap (no Go dependency) | In Progress | `curl \| bash` `install.sh` at repo root, modeled selectively on `smaqit-extensions/install.sh`; adopts ancestor-dir resolution into `agentic-cms update` only — no global-install, `--scope`, `uninstall`, or multi-agent global directory surface |
| 007 | Rename .agentic-cms/bin to .agentic-cms/scripts | In Progress | `bin/` collides with the generic IDE build-output-exclude convention but holds committed source; Tier 1 rename only (scaffold source + all living references), no update-time migration for already-installed projects. Started before 005 completes (user decision); 005 reconciles its own ~150 bin/ references against the renamed baseline afterward |

## Completed

| ID | Title | Completed | Notes |
|----|-------|-----------|-------|
| 001 | Installer smoke test script | 2026-07-31 | Bash end-to-end smoke test of `agentic-cms init` via `make smoke-test` |
| 002 | Draft content state | 2026-08-15 | PR #2; `status: draft` frontmatter + `docs/<topic>/drafts/` convention, `ac-page --status`/`promote` toolkit support; shipped as part of v0.2.0 |
| 004 | Archive content lifecycle | 2026-08-16 | PR #3; `status: archived` + `docs/<topic>/archive/` mirroring drafts, `ac-page archive`, archived items stay indexed under `## Archived`; renamed content-new-item → content-manage-item; released as v0.3.0 |

## Abandoned

| ID | Title | Date | Reason |
|----|-------|------|--------|
