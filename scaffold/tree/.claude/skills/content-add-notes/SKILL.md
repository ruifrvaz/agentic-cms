---
name: content-add-notes
description: Add notes, observations, or updates to an existing content item. Use when the user wants to annotate a page — "add a note to X", "append this to Y", "note that Z changed".
license: MIT
compatibility: Requires bash and python3 (uses the .agentic-cms/scripts toolkit)
allowed-tools: Read Edit Grep Glob Bash(.agentic-cms/scripts/*)
---

# content-add-notes — annotate existing content

Read `CONTENT.md` at the project root first if you haven't this session. Toolkit
contract: `.agentic-cms/scripts/README.md`.

## Steps

1. **Find the item** (deterministic):
   ```sh
   .agentic-cms/scripts/ac-search <subject terms>
   .agentic-cms/scripts/ac-index list
   ```
   If no page exists, offer `content-manage-item` instead.
2. **Append** under the page's `## Notes` section (create the section if missing):

   ```markdown
   ### YYYY-MM-DD
   <the note>
   ```

3. **Assess impact** (judgment): does the note change any claim in the page body
   or in pages that reference it (`ac-page meta <file>` → `refs:`, plus
   `ac-search` for inbound links)? If yes, update those in place and flag
   superseded claims — a contradicting note must not just sit at the bottom.
4. **Re-rate, touch, log, verify** (deterministic + one judgment call): a
   note can make a page more sensitive than it was — a "note that changed"
   check on classification, every time, not just when it looks obviously
   sensitive. Compare the note against `CONTENT.md`'s C0–C3 rubric:
   ```sh
   .agentic-cms/scripts/ac-page meta <file>            # see the current classification
   ```
   If the note pushes the page's sensitivity higher, raise it — **never
   lower** a rating here; only the user does that:
   ```sh
   .agentic-cms/scripts/ac-page classify <file> <higher-level>
   ```
   (`ac-page classify` also touches `updated:` and restamps the hash, so
   skip a separate `ac-page touch` for a re-rated file.) Then, for every
   edited file that was **not** re-rated:
   ```sh
   .agentic-cms/scripts/ac-page touch <file>
   ```
   Finally, log and verify:
   ```sh
   .agentic-cms/scripts/ac-log append notes "<page>"
   .agentic-cms/scripts/ac-classify check <every file touched>
   ```
   `"clean": true` confirms the rating (raised or unchanged) still matches
   the page as it now stands — this is the guaranteed check for setups
   where the post-tool-use hook isn't installed or can't block.

## Rules

- Notes are the user's voice — record them faithfully, lightly edited for clarity.
- Small notes do not need index updates; only log them.
- Classification only ever moves up here — raise on a more-sensitive note,
  otherwise leave it. Lowering a rating requires the user explicitly.
