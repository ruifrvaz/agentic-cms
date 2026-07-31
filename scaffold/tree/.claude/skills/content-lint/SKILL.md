---
name: content-lint
description: Health-check the content base — contradictions, stale claims, orphan pages, missing cross-references, broken links, index/reality drift. Use for "lint the wiki", "check content health", or periodically after several ingests.
---

# content-lint — health check

Read `CONTENT.md` at the project root first if you haven't this session.

## Checks

1. **Mechanical** (fix directly):
   - Broken relative links in `docs/` and `wiki/`
   - Pages missing from `wiki/index.md`; index entries pointing nowhere
   - Missing/malformed frontmatter; stale `updated:` dates on files you fix
2. **Structural** (fix directly, note in report):
   - Orphan pages with no inbound links — link them from the right hub or index
   - Concepts/entities mentioned across ≥2 pages that lack their own page — list
     as candidates, create only with user approval
3. **Content-level** (report, do not silently fix):
   - Contradictions between pages
   - Claims likely superseded by newer sources (compare source dates)
   - Data gaps a `content-research` run could fill

## Steps

1. Sweep `docs/` and `wiki/` (grep for links, parse frontmatter, cross-check index).
2. Apply mechanical + structural fixes.
3. Report: fixed / candidates / content-level findings, each with file paths.
4. **Append** to `wiki/log.md`: `## [YYYY-MM-DD] lint | <n> fixed, <m> flagged`.

## Rules

- Never alter the meaning of a page during lint — meaning changes require the user.
- `raw/` is out of scope except for detecting never-ingested sources.
