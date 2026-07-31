---
name: content-research
description: Research a topic on the web and file the findings into the CMS. Use when the user wants to gather external information — "research X", "find out about Y", "what's the latest on Z" — and keep the results in the content base.
---

# content-research — web research into the content base

Read `CONTENT.md` at the project root first if you haven't this session.

## Steps

1. **Scope**: confirm the research question and which topic the findings belong to.
   Check `wiki/index.md` for what is already known — do not re-research settled ground;
   note what exists and research the gaps.
2. **Delegate** to the `content-researcher` subagent with: the research question, the
   relevant existing pages (paths), and instructions to return findings with source
   URLs. If subagents are unavailable, do the web research yourself.
3. **File the findings**:
   - New material → a content item in `docs/<topic>/` (via the content-new-item
     conventions), with URLs listed in `sources:`.
   - Updates to existing pages → edit in place; flag contradictions with existing
     claims explicitly rather than silently overwriting.
   - Update affected wiki entity/concept pages and cross-links.
4. **Update** `wiki/index.md` and append to `wiki/log.md`:
   `## [YYYY-MM-DD] research | <question>` with the pages touched.

## Rules

- Every researched claim carries its source URL.
- Distinguish clearly between what sources say and your own synthesis.
- Findings the user asked for but that answer a one-off question still get filed —
  explorations compound in the wiki.
