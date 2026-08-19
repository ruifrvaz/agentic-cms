---
name: content-lint
description: Health-check the content base — contradictions, stale claims, orphan pages, missing cross-references, broken links, index/reality drift. Use for "lint the wiki", "check content health", or periodically after several ingests.
license: MIT
compatibility: Requires bash and python3 (uses the .agentic-cms/scripts toolkit)
allowed-tools: Read Edit Grep Glob Bash(.agentic-cms/scripts/*)
---

# content-lint — health check

Read `CONTENT.md` at the project root first if you haven't this session. Toolkit
contract: `.agentic-cms/scripts/README.md`.

## Phase 1 — mechanical (deterministic, fix directly)

```sh
.agentic-cms/scripts/ac-links check     # broken links → fix each, ac-page touch the file
.agentic-cms/scripts/ac-index check     # drift → ac-index remove / ac-index add per finding
.agentic-cms/scripts/ac-inventory       # raw_uningested → report as import candidates
.agentic-cms/scripts/ac-classify sweep  # classification pass — see below
```

Also verify frontmatter integrity: `ac-page meta` on any page you touch; fix
missing/malformed frontmatter fields against the templates.

**Classification** (`ac-classify sweep` — this is content-lint's periodic
catch-all for whatever the write-time hooks/gates never saw: pre-existing
drift, edits committed with `--no-verify`, or edits made where no hook was
installed):
- `floor_violation` entries fix directly, no judgment needed — this
  mirrors exactly what the post-tool-use hook already does automatically,
  so treat it the same way: `ac-page classify <path> <implied_level>`.
  Exception: if the user has told you a specific page's floor hit is a
  false positive, don't re-raise it — report it and let the user decide
  whether to ack it (`ac-page classify <path> <level> --ack-floor`).
  Acking is always the user's call, never yours; pages the user already
  acked show `"acked": true` with no violation and need no action.
- `pages` with `"valid": false` (a hand-edited or corrupted enum value):
  fix directly to the nearest sensible level by reading the page.

Re-run `ac-links check`, `ac-index check`, and `ac-classify sweep` after
fixing — all three must report `"clean": true` before moving on.

## Phase 2 — structural (fix directly, note in report)

- **Orphans**: for each wiki page in `ac-index list`, run
  `ac-search "<page-slug>"` — pages with no inbound links outside the index get
  linked from the right hub page. Drafts (`status: draft`, under
  `docs/<topic>/drafts/`) never appear in `ac-index list` by design, so they
  never reach this check. Skip entries in the `archived` section: having no
  inbound links is an archived page's expected steady state, not a defect.
- **Missing pages**: concepts/entities recurring across ≥2 pages (spot via
  `ac-search`) that lack their own page — list as candidates, create only with
  user approval.
- **Stale ratings** (`ac-classify sweep`'s `stale: true` entries): re-read
  the page against `CONTENT.md`'s C0–C3 rubric and re-rate
  (`ac-page classify <path> <level>`). Only apply a raise directly — if the
  correct rating would *lower* what's currently there, flag it for the
  user instead of touching it (only the user may lower a rating).
- **Unrated pages** (`unrated: true` — pre-existing content from before
  this feature, absent field defaults to C1): rate against the rubric and
  set explicitly, same raise-only judgment as stale ratings above.
- **C2+ bleed** (`ac-classify sweep`'s `bleed` array): a `wiki/index.md` or
  `wiki/log.md` line is leaking figures, PII, or quoted content from the
  page it describes — this is the exact failure mode classification
  exists to prevent, flag it prominently. Rewrite the leaking line to an
  opaque summary (judgment call, not mechanical); if the leak itself
  suggests the summary page's own rating needs raising, do that too.

## Phase 3 — content-level (judgment: report, do not silently fix)

- Contradictions between pages
- Claims likely superseded by newer sources (compare `sources:` dates via
  `ac-page meta`) — skip `status: archived` pages; retired content is not
  expected to stay current
- Data gaps a `content-research` run could fill

## Wrap up

```sh
.agentic-cms/scripts/ac-log append lint "<n> fixed, <m> flagged" "<summary of findings>"
```

Report: fixed / candidates / content-level findings, each with file paths.

## Rules

- Never alter the meaning of a page during lint — meaning changes require the user.
- `raw/` is out of scope except for the un-ingested check.
