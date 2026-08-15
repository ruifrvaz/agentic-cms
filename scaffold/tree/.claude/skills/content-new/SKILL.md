---
name: content-new
description: Create a new content topic in the CMS. Use when the user wants to start a new topic, area, or collection of documentation — "new topic", "start documenting X", "create a section for Y". Sets up docs/<topic>/, seeds the wiki, updates index and log.
license: MIT
compatibility: Requires bash and python3 (uses the .agentic-cms/bin toolkit)
allowed-tools: Read Grep Glob Bash(.agentic-cms/bin/*)
---

# content-new — create a new topic

Read `CONTENT.md` at the project root first if you haven't this session. All
`ac-*` commands live in `.agentic-cms/bin/`, run from the project root, and
return JSON — check `"ok"` after every call (contract: `.agentic-cms/bin/README.md`).

## Steps

1. **Clarify** the topic name (kebab-case) and its scope in one exchange if not
   obvious from the request.
2. **Check for duplicates** deterministically:
   ```sh
   .agentic-cms/bin/ac-index list
   .agentic-cms/bin/ac-search <topic> <synonyms>
   ```
   If a matching or near-duplicate topic exists, propose extending it instead.
3. **Create the topic page**:
   ```sh
   .agentic-cms/bin/ac-page new topic docs/<topic>/README.md --title "<Title>" --topic <topic>
   ```
   Then edit the created file: write the "About this topic" section from what the
   user told you (judgment step — the script only guarantees structure).
4. **Seed the wiki** only for entities/concepts the user explicitly named:
   ```sh
   .agentic-cms/bin/ac-page new entity wiki/entities/<slug>.md --title "<Name>"
   .agentic-cms/bin/ac-index add entities wiki/entities/<slug>.md "<one-line summary>"
   ```
   Do not invent speculative pages.
5. **Register and log**:
   ```sh
   .agentic-cms/bin/ac-index add topics docs/<topic>/README.md "<one-line summary>"
   .agentic-cms/bin/ac-log append new "<topic>" "<what was created>"
   .agentic-cms/bin/ac-index check
   ```
   `check` must return `"clean": true`; fix anything it reports before finishing.

## Rules

- Greenfield-friendly: if `docs/` is empty this is likely the project's first topic.
- Never create empty placeholder items — items are added via `content-manage-item`.
