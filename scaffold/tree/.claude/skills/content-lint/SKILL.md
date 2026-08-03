---
name: content-lint
description: Health-check the content base — contradictions, stale claims, orphan pages, missing cross-references, broken links, index/reality drift. Use for "lint the wiki", "check content health", or periodically after several ingests.
license: MIT
compatibility: Requires bash and python3 (uses the .agentic-cms/bin toolkit)
allowed-tools: Read Edit Grep Glob Bash(.agentic-cms/bin/*)
---

# content-lint — health check

Read `CONTENT.md` at the project root first if you haven't this session. Toolkit
contract: `.agentic-cms/bin/README.md`.

## Phase 1 — mechanical (deterministic, fix directly)

```sh
.agentic-cms/bin/ac-links check     # broken links → fix each, ac-page touch the file
.agentic-cms/bin/ac-index check     # drift → ac-index remove / ac-index add per finding
.agentic-cms/bin/ac-inventory       # raw_uningested → report as import candidates
```

Also verify frontmatter integrity: `ac-page meta` on any page you touch; fix
missing/malformed frontmatter fields against the templates.

Re-run `ac-links check` and `ac-index check` after fixing — both must report
`"clean": true` before moving on.

## Phase 2 — structural (fix directly, note in report)

- **Orphans**: for each wiki page in `ac-index list`, run
  `ac-search "<page-slug>"` — pages with no inbound links outside the index get
  linked from the right hub page.
- **Missing pages**: concepts/entities recurring across ≥2 pages (spot via
  `ac-search`) that lack their own page — list as candidates, create only with
  user approval.

## Phase 3 — content-level (judgment: report, do not silently fix)

- Contradictions between pages
- Claims likely superseded by newer sources (compare `sources:` dates via `ac-page meta`)
- Data gaps a `content-research` run could fill

## Wrap up

```sh
.agentic-cms/bin/ac-log append lint "<n> fixed, <m> flagged" "<summary of findings>"
```

Report: fixed / candidates / content-level findings, each with file paths.

## Rules

- Never alter the meaning of a page during lint — meaning changes require the user.
- `raw/` is out of scope except for the un-ingested check.
