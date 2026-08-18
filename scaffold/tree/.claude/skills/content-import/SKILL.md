---
name: content-import
description: Import raw sources (PPTX, DOCX, PDF, markdown, text, or a whole folder of unstructured content) into the CMS structure. Use for "import this file/deck/doc", "ingest X", or bringing an existing brownfield project's content into raw/, docs/ and the wiki.
license: MIT
compatibility: Requires bash and python3 (uses the .agentic-cms/scripts toolkit); conversion may install pandoc, markitdown, or python-pptx
allowed-tools: Read Write Edit Grep Glob Bash
---

# content-import — ingest raw sources

Read `CONTENT.md` at the project root first if you haven't this session. Toolkit
contract: `.agentic-cms/scripts/README.md` — check `"ok"` after every `ac-*` call.

## Two modes

**Single-source ingest**: one file — handle inline below.
**Brownfield sweep**: a folder of unstructured content — delegate to the
`content-importer` subagent (it follows this same procedure per file).

## Steps (per source)

1. **Secure the raw copy** (deterministic):
   ```sh
   cp -n <original> raw/            # copy, never move; never overwrite
   ```
   Extracted media goes to `raw/assets/`.
2. **Convert to markdown**:
   - DOCX: `pandoc -f docx -t gfm --extract-media=raw/assets <file>`
   - PPTX: `markitdown <file>` or a `python-pptx` script — slide titles →
     headings, bullets → lists, speaker notes → blockquotes
   - PDF: `pdftotext` or `markitdown`; OCR only if the text layer is missing
   Install what's missing (`pip install markitdown python-pptx --break-system-packages`).
3. **Rate, then file it** (deterministic shell, judgment content): rate the
   converted content against `CONTENT.md`'s C0–C3 scale before creating the
   page — you've just read the source, so you know whether it carries
   financials, PII, or correspondence.
   ```sh
   .agentic-cms/scripts/ac-inventory        # pick the topic; ask only if ambiguous
   .agentic-cms/scripts/ac-page new doc docs/<topic>/<item>.md --title "<Title>" --topic <topic> --raw-path raw/<file> --classification <C0|C1|C2|C3>
   ```
   Merge the converted markdown into the created page (clean heading levels,
   tables, export artifacts) and set frontmatter `sources: [raw/<file>]`. If
   the merged content turns out more sensitive than the initial rating,
   re-rate: `ac-page classify <path> <level>`.
4. **Wiki integration**: rate the source summary too — a summary of a
   confidential source is not automatically confidential itself (an opaque
   summary of a C2 source can be C1), but check against the rubric rather
   than assuming.
   ```sh
   .agentic-cms/scripts/ac-page new source wiki/sources/<source-slug>.md --title "<Source Title>" --raw-path raw/<file> --classification <C0|C1|C2|C3>
   .agentic-cms/scripts/ac-index add sources wiki/sources/<source-slug>.md "<one-line summary>"
   ```
   Write the summary/takeaways (judgment) — remember the bleed rule if the
   source is C2+: no figures, no personal details, no quoted content in the
   index summary. Then update every entity/concept page the source
   materially affects — new facts, contradictions (flag, don't overwrite),
   cross-references. `ac-page touch` every page you edit; `ac-index add`
   every page you create.
5. **Log and verify**:
   ```sh
   .agentic-cms/scripts/ac-log append import "<source>" "<pages created/updated>"
   .agentic-cms/scripts/ac-index check && .agentic-cms/scripts/ac-links check && .agentic-cms/scripts/ac-classify check <all pages touched this run>
   ```
   All three must report `"clean": true`.

## Rules

- `raw/` is append-only — cleanup happens in `docs/`, never by editing raw files.
- Sweeps: propose the file → topic mapping for approval before mass-processing.
- Conversion fidelity beats prettiness; never drop content silently — flag bad
  conversions.
- Rate every page a sweep creates, same as a single-source ingest — a
  brownfield folder is exactly the case most likely to contain financials or
  correspondence. Subagent-mediated writes are not guaranteed to trigger an
  agent-lifecycle hook on every platform, so step 5's `ac-classify check` is
  not optional here — it is the guaranteed backstop for sweep-created pages.
