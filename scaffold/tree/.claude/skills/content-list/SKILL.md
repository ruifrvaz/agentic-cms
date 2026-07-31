---
name: content-list
description: List and summarize the contents of the CMS — topics, items, wiki pages, recent activity. Use for "what's in the content base", "list topics", "what changed recently", or a health overview.
---

# content-list — inventory and status

Read `CONTENT.md` at the project root first if you haven't this session.

## Steps

1. **Read** `wiki/index.md` for the catalog and
   `grep "^## \[" wiki/log.md | tail -10` for recent activity.
2. **Verify against reality**: quick `ls`/glob of `docs/`, `wiki/`, and `raw/` —
   report any drift (pages missing from the index, index entries pointing at
   deleted files, raw sources never ingested).
3. **Report** concisely:
   - Topics with item counts
   - Wiki pages by category (entities / concepts / sources)
   - Raw sources not yet ingested
   - Last 5 log entries
   - Drift found, if any
4. If drift was found, offer to fix the index (mechanical fixes only — this is a
   light lint, not a full content lint).

## Rules

- Read-only except for offered index fixes. Do not modify content pages.
- Do not log pure listings; log only if the index was repaired
  (`## [YYYY-MM-DD] lint | index repair`).
