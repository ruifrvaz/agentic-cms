# Task Planning

## Active

| ID | Title | Status | Notes |
|----|-------|--------|-------|
| 003 | Register mode for reference files | Not Started | Lightweight `content-import` Register mode: wiki/sources sidecar + index + log for archival/reference files (e.g. invoices), no markdown conversion; start after 002 merges |
| 005 | First-class content classification (auto CIA rating) | Not Started | `classification: C0–C3` frontmatter auto-assigned by agent at write time against a CONTENT.md rubric; toolkit validation (`ac-page --classification`/`classify`, `ac-inventory` reporting); lint bleed/ratchet enforcement. Start after 004 merges (content-manage-item rename). Segregation/vault out of scope — follow-up task once downstream pattern settles |

## Completed

| ID | Title | Completed | Notes |
|----|-------|-----------|-------|
| 001 | Installer smoke test script | 2026-07-31 | Bash end-to-end smoke test of `agentic-cms init` via `make smoke-test` |
| 002 | Draft content state | 2026-08-15 | PR #2; `status: draft` frontmatter + `docs/<topic>/drafts/` convention, `ac-page --status`/`promote` toolkit support; shipped as part of v0.2.0 |
| 004 | Archive content lifecycle | 2026-08-16 | PR #3; `status: archived` + `docs/<topic>/archive/` mirroring drafts, `ac-page archive`, archived items stay indexed under `## Archived`; renamed content-new-item → content-manage-item; released as v0.3.0 |

## Abandoned

| ID | Title | Date | Reason |
|----|-------|------|--------|
