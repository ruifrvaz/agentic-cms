---
name: content-export
description: Export content from the CMS into a PPTX presentation or DOCX document. Use for "make a deck from X", "export Y as a Word doc", "summarize topic Z into a presentation".
license: MIT
compatibility: Requires bash and python3 (uses the .agentic-cms/bin toolkit); building files may install python-pptx or pandoc
allowed-tools: Read Write Grep Glob Bash
---

# content-export — export to PPTX / DOCX

Read `CONTENT.md` at the project root first if you haven't this session. Toolkit
contract: `.agentic-cms/bin/README.md`.

## Steps

1. **Scope**: confirm what to export (topic, page set, question-driven summary),
   format (pptx/docx), audience/length.
2. **Gather** (deterministic):
   ```sh
   .agentic-cms/bin/ac-index list
   .agentic-cms/bin/ac-search <topic terms>
   ```
   Collect and read the relevant pages. Export from the wiki/docs synthesis, not
   from raw sources.
3. **Delegate** to the `content-exporter` subagent with: the page paths, format,
   audience, and output path `exports/<name>.<ext>`. If subagents are
   unavailable, build the file yourself (`python-pptx` for decks, pandoc
   `-t docx` for documents; install if needed).
4. **Verify deterministically**: the file must exist and be non-trivial —
   ```sh
   ls -la exports/<name>.<ext>
   ```
   and re-extract its text programmatically to confirm nothing is empty or
   truncated. Tell the user where it is.
5. **Log**:
   ```sh
   .agentic-cms/bin/ac-log append export "<name>.<ext>" "<pages exported>"
   ```

## Rules

- Exports are derived artifacts: they live in `exports/` (create it if missing),
  never inside `raw/`, `docs/`, or `wiki/`.
- Decks: one idea per slide, titles that assert the point, details in speaker notes.
