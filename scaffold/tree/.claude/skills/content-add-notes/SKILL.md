---
name: content-add-notes
description: Add notes, observations, or updates to an existing content item. Use when the user wants to annotate a page — "add a note to X", "append this to Y", "note that Z changed".
license: MIT
compatibility: Requires bash and python3 (uses the .agentic-cms/bin toolkit)
allowed-tools: Read Edit Grep Glob Bash(.agentic-cms/bin/*)
---

# content-add-notes — annotate existing content

Read `CONTENT.md` at the project root first if you haven't this session. Toolkit
contract: `.agentic-cms/bin/README.md`.

## Steps

1. **Find the item** (deterministic):
   ```sh
   .agentic-cms/bin/ac-search <subject terms>
   .agentic-cms/bin/ac-index list
   ```
   If no page exists, offer `content-new-item` instead.
2. **Append** under the page's `## Notes` section (create the section if missing):

   ```markdown
   ### YYYY-MM-DD
   <the note>
   ```

3. **Assess impact** (judgment): does the note change any claim in the page body
   or in pages that reference it (`ac-page meta <file>` → `refs:`, plus
   `ac-search` for inbound links)? If yes, update those in place and flag
   superseded claims — a contradicting note must not just sit at the bottom.
4. **Touch and log** (deterministic):
   ```sh
   .agentic-cms/bin/ac-page touch <file>          # every file you edited
   .agentic-cms/bin/ac-log append notes "<page>"
   ```

## Rules

- Notes are the user's voice — record them faithfully, lightly edited for clarity.
- Small notes do not need index updates; only log them.
