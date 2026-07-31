---
name: content-new
description: Create a new content topic in the CMS. Use when the user wants to start a new topic, area, or collection of documentation — "new topic", "start documenting X", "create a section for Y". Sets up docs/<topic>/, seeds the wiki, updates index and log.
---

# content-new — create a new topic

Read `CONTENT.md` at the project root first if you haven't this session.

## Steps

1. **Clarify** the topic name (kebab-case) and its scope in one exchange if not
   obvious from the request.
2. **Check** `wiki/index.md` — does this topic (or a near-duplicate) already exist?
   If yes, propose extending it instead of creating a new one.
3. **Create** `docs/<topic>/README.md` from `.agentic-cms/templates/topic.md`,
   filling `{{TITLE}}`, `{{TOPIC}}`, `{{DATE}}` (today, YYYY-MM-DD). Write the
   "About this topic" section from what the user told you.
4. **Seed the wiki** only if the topic implies obvious entities/concepts the user
   named — create those pages from templates. Do not invent speculative pages.
5. **Update** `wiki/index.md`: add the topic under "Topics" with a one-line summary.
6. **Append** to `wiki/log.md`: `## [YYYY-MM-DD] new | <topic>` plus 1-2 lines on
   what was created.

## Rules

- Greenfield-friendly: if `docs/` is empty this is likely the project's first topic;
  keep it simple.
- Never create empty placeholder items — items are added via `content-new-item`.
