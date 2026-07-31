---
name: content-new-item
description: Create a single new content item (markdown page) inside an existing topic in docs/. Use when the user wants to add a document, article, spec, or page about something specific — "add a page on X", "write up Y", "document Z".
---

# content-new-item — create a content item

Read `CONTENT.md` at the project root first if you haven't this session.

## Steps

1. **Locate the topic**: check `wiki/index.md` and `docs/` for the topic this item
   belongs to. If no topic fits, run `content-new` first (confirm with the user).
2. **Check for duplicates** in the topic folder and the index. If a page for this
   subject exists, extend it instead (or use `content-add-notes`).
3. **Create** `docs/<topic>/<item>.md` from `.agentic-cms/templates/doc.md`.
   Kebab-case filename, no dates in the name. Fill frontmatter completely; write a
   2-3 sentence Summary and the Content body from what you know or what the user
   provided. Cite raw sources in `sources:` when the content derives from them.
4. **Integrate with the wiki**: if the item introduces or substantially adds to an
   entity or concept, create/update the corresponding `wiki/entities/` or
   `wiki/concepts/` page and cross-link both ways (`refs:` frontmatter + inline links).
5. **Update** the topic's `README.md` item list and `wiki/index.md`.
6. **Append** to `wiki/log.md`: `## [YYYY-MM-DD] new-item | <topic>/<item>`.

## Rules

- One subject per file. If the user asks for something spanning two subjects, make
  two items and cross-link them.
- Substance over scaffolding: a new item must contain real content, not just headers.
