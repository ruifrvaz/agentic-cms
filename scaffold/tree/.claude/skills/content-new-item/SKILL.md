---
name: content-new-item
description: Create a single new content item (markdown page) inside an existing topic in docs/. Use when the user wants to add a document, article, spec, or page about something specific — "add a page on X", "write up Y", "document Z".
license: MIT
compatibility: Requires bash and python3 (uses the .agentic-cms/bin toolkit)
allowed-tools: Read Write Edit Grep Glob Bash(.agentic-cms/bin/*)
---

# content-new-item — create a content item

Read `CONTENT.md` at the project root first if you haven't this session. All
`ac-*` commands live in `.agentic-cms/bin/`, return JSON — check `"ok"` after
every call (contract: `.agentic-cms/bin/README.md`).

## Steps

1. **Locate the topic**:
   ```sh
   .agentic-cms/bin/ac-inventory
   ```
   Pick the topic from `topics`. If none fits, run `content-new` first (confirm
   with the user).
2. **Check for duplicates**:
   ```sh
   .agentic-cms/bin/ac-search <subject> <synonyms>
   ```
   If a page for this subject exists, extend it (or use `content-add-notes`).
3. **Create from template** (kebab-case filename, no dates in the name):
   ```sh
   .agentic-cms/bin/ac-page new doc docs/<topic>/<item>.md --title "<Title>" --topic <topic>
   ```
   Then write the real content (judgment step): 2-3 sentence Summary, the body,
   frontmatter `tags:` and `sources:` (raw paths it derives from).
4. **Integrate with the wiki**: if the item introduces or substantially adds to
   an entity/concept, create the wiki page (`ac-page new entity|concept ...` +
   `ac-index add ...`) or update the existing one, and cross-link both ways
   (inline links + `refs:` frontmatter). After editing any existing page:
   ```sh
   .agentic-cms/bin/ac-page touch <edited-file>
   ```
5. **Register, log, verify**:
   ```sh
   .agentic-cms/bin/ac-index add topics docs/<topic>/<item>.md "<one-line summary>"
   .agentic-cms/bin/ac-log append new-item "<topic>/<item>"
   .agentic-cms/bin/ac-index check && .agentic-cms/bin/ac-links check
   ```
   Both must report `"clean": true` before finishing. Also add the item to the
   topic's `README.md` item list.

## Rules

- One subject per file; two subjects → two items, cross-linked.
- Substance over scaffolding: a new item must contain real content, not just headers.
