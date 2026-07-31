---
name: content-import
description: Import raw sources (PPTX, DOCX, PDF, markdown, text, or a whole folder of unstructured content) into the CMS structure. Use for "import this file/deck/doc", "ingest X", or bringing an existing brownfield project's content into raw/, docs/ and the wiki.
---

# content-import — ingest raw sources

Read `CONTENT.md` at the project root first if you haven't this session.

## Two modes

**Single-source ingest**: one file (PPTX, DOCX, PDF, md, txt, html).
**Brownfield sweep**: a folder of existing unstructured content, possibly from
another project. Delegate sweeps to the `content-importer` subagent; single files
you can handle inline.

## Steps (per source)

1. **Secure the raw copy**: if the file is not already under `raw/`, COPY it there
   (never move, never modify the original). Extracted images go to `raw/assets/`.
2. **Convert to markdown**:
   - DOCX: `pandoc -f docx -t gfm --extract-media=raw/assets <file>` (fallback:
     `python-docx` via a short script, or `markitdown`).
   - PPTX: `python-pptx` script or `markitdown <file>` — capture slide titles as
     headings, bullets as lists, and speaker notes as blockquotes.
   - PDF: `pdftotext` or `markitdown`; OCR only if the text layer is missing.
   - Install what you need (`pip install markitdown python-pptx`); prefer whatever
     is already available.
3. **File it**: place the cleaned markdown as `docs/<topic>/<item>.md` with full
   frontmatter (`sources:` pointing at the raw path). Choose the topic from
   `wiki/index.md`; ask the user only if genuinely ambiguous.
4. **Integrate into the wiki**:
   - Create `wiki/sources/<source>.md` from the source template: summary, key
     takeaways, pages updated.
   - Update every entity/concept page the source materially affects — new facts,
     contradictions (flag them), new cross-references. One source touching many
     pages is normal.
5. **Update** `wiki/index.md`; append to `wiki/log.md`:
   `## [YYYY-MM-DD] import | <source>` listing pages created/updated.

## Rules

- `raw/` is append-only. Conversion artifacts and cleanup happen in `docs/`, never
  by editing the raw file.
- For sweeps: propose the topic mapping (which files go where) to the user before
  executing, then process source by source.
