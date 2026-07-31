---
name: content-export
description: Export content from the CMS into a PPTX presentation or DOCX document. Use for "make a deck from X", "export Y as a Word doc", "summarize topic Z into a presentation".
---

# content-export — export to PPTX / DOCX

Read `CONTENT.md` at the project root first if you haven't this session.

## Steps

1. **Scope**: confirm what to export (a topic, a set of pages, a question-driven
   summary), the format (pptx/docx), and the audience/length if not stated.
2. **Gather**: via `wiki/index.md`, collect the relevant pages and read them. The
   wiki is the synthesis — export from it, not from raw sources.
3. **Delegate** to the `content-exporter` subagent with: the gathered page paths,
   format, audience, and output path `exports/<name>.<ext>`. If subagents are
   unavailable, build the file yourself (`python-pptx` for decks, pandoc
   `-t docx` for documents; `pip install python-pptx` if needed).
4. **Verify** the output file exists and opens (e.g. re-extract its text) and tell
   the user where it is.
5. **Append** to `wiki/log.md`: `## [YYYY-MM-DD] export | <name>.<ext>`.

## Rules

- Exports are derived artifacts: they live in `exports/` (create it if missing),
  never inside `raw/`, `docs/`, or `wiki/`.
- Decks: one idea per slide, titles that assert the point, details in speaker notes.
