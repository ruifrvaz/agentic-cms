---
name: content-add-notes
description: Add notes, observations, or updates to an existing content item. Use when the user wants to annotate a page — "add a note to X", "append this to Y", "note that Z changed".
---

# content-add-notes — annotate existing content

Read `CONTENT.md` at the project root first if you haven't this session.

## Steps

1. **Find the item**: locate the target page via `wiki/index.md` or `docs/`. If it
   doesn't exist, offer to create it with `content-new-item` instead.
2. **Append** under the page's `## Notes` section (create the section if missing):

   ```markdown
   ### YYYY-MM-DD
   <the note>
   ```

3. **Assess impact**: does the note change any claim in the page body or in wiki
   pages that reference it? If yes, update those in place and flag superseded
   claims; a note that contradicts the body must not just sit at the bottom.
4. **Touch** the `updated:` frontmatter date.
5. **Append** to `wiki/log.md`: `## [YYYY-MM-DD] notes | <page>`.

## Rules

- Notes are the user's voice — record them faithfully, lightly edited for clarity.
- Small notes do not need index updates; only log them.
