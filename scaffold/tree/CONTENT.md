# CONTENT.md — Content Management Schema

This file is the schema for the content management system in this project. It tells
any coding agent (Claude Code, Codex, Copilot, ...) how content is structured, what
the conventions are, and which workflows to follow. Read it before touching anything
under `raw/`, `docs/`, or `wiki/`. It is the single most important file of the CMS:
it is what makes an agent a disciplined content maintainer rather than a generic
chatbot. Co-evolve it with the user as conventions solidify, but never change it
silently — propose schema changes explicitly.

## The three layers

**`raw/` — immutable sources.** The user's curated collection of source material:
articles, papers, exports, PPTX/DOCX files, transcripts, images (in `raw/assets/`).
Agents READ from this layer but NEVER modify or delete anything in it. It is the
source of truth. New material always lands here first.

**`docs/` — organized documentation.** Structured markdown derived from raw sources
and authored content, organized by topic. This is the working layer used during
agentic sessions: clean, readable, one topic per folder, one item per file. Agents
create and update files here through the content skills.

**`wiki/` — the synthesis layer.** An agent-maintained, interlinked knowledge base
that sits on top of `docs/`: entity pages, concept pages, source summaries, plus
`index.md` (catalog) and `log.md` (journal). The agent owns this layer entirely.
The user reads it; the agent writes it. Knowledge is compiled once and kept
current — not re-derived on every question.

## Directory map

```
raw/                  immutable sources (never modify)
  assets/             images and binary attachments
docs/                 organized topical documentation
  <topic>/            one folder per topic
    <item>.md         one file per content item
    drafts/           optional: work-in-progress items (status: draft)
wiki/
  index.md            catalog of every page: link + one-line summary, by category
  log.md              append-only journal of operations
  entities/           pages about specific things (people, products, systems, orgs)
  concepts/           pages about ideas, patterns, themes
  sources/            one summary page per ingested raw source
.agentic-cms/
  bin/                deterministic CLI toolkit (ac-page, ac-index, ac-log,
                      ac-links, ac-inventory, ac-search) — JSON on stdout;
                      contract in bin/README.md
  templates/          markdown templates for every page type
  VERSION             scaffold version installed by agentic-cms
.claude/
  skills/             content-* skills (workflows)
  agents/             researcher, importer, exporter subagents
```

## The toolkit

Mechanical operations MUST go through `.agentic-cms/bin/` instead of hand-editing:
creating pages from templates (`ac-page new`), reading/updating frontmatter
(`ac-page meta|touch`), index maintenance (`ac-index add|remove|check`), logging
(`ac-log append`), link checking (`ac-links check`), inventory (`ac-inventory`),
and term search (`ac-search`). Every command prints one JSON object; check `"ok"`
before proceeding, and finish write operations with `ac-index check` /
`ac-links check` returning `"clean": true`. Judgment work — summarizing,
synthesis, wording — stays with the agent. Full contract:
`.agentic-cms/bin/README.md`.

## Conventions

- **Filenames**: kebab-case, descriptive, no dates in names (`transformer-architecture.md`).
- **Frontmatter**: every markdown page in `docs/` and `wiki/` carries YAML frontmatter:

  ```yaml
  ---
  title: Human Readable Title
  type: doc | entity | concept | source | note
  topic: <topic>            # docs/ pages only
  tags: [tag1, tag2]
  created: YYYY-MM-DD
  updated: YYYY-MM-DD
  sources: [raw/file.pdf]   # raw sources this page draws on
  refs: [wiki/concepts/x.md]# pages this page depends on
  status: final             # `draft` while work-in-progress; `final` once ready
  ---
  ```

- **Drafts**: `status: draft` marks a page as not yet first-class content — a
  brainstorm or half-formed note that will be iterated on before it graduates.
  Set via `ac-page new ... --status draft`, writing to
  `docs/<topic>/drafts/<item>.md` (near its eventual home). Drafts are NOT
  wired into `wiki/index.md` or `wiki/log.md` while draft — no index entry,
  no log entry, no orphan-page flag from `content-lint`. Only `status: final`
  (the template default) or an absent `status` field counts as first-class
  content subject to the normal wiki bookkeeping below. **Promotion**:
  `ac-page promote docs/<topic>/drafts/<item>.md docs/<topic>/<item>.md`
  moves the file, sets `status: final`, and touches `updated:` in one call;
  then run the normal `content-new-item` register+log step (index entry, log
  entry) against the new path.
- **Cross-links**: use relative markdown links (`[X](../concepts/x.md)`). Obsidian-style
  `[[wikilinks]]` are acceptable if the user works in Obsidian; pick one per project
  and record the choice here.
- **Templates**: always start new pages from `.agentic-cms/templates/`. Do not invent
  ad-hoc page shapes.
- **index.md**: updated on EVERY operation that adds, renames, or removes a page.
  Format: one line per page — `- [Title](path) — one-line summary`.
- **log.md**: append-only. Every operation appends one entry with the header format
  `## [YYYY-MM-DD] <operation> | <subject>` so the log is greppable
  (`grep "^## \[" wiki/log.md | tail -5`). Operations: `init`, `new`, `new-item`,
  `import`, `research`, `notes`, `lint`, `export`.
- **Contradictions**: when new material contradicts an existing page, do not silently
  overwrite. Note the contradiction inline (`> ⚠ Contradicts [X](...) — newer source
  says ...`), update if the newer source is authoritative, and log it.

## Operations

**Ingest** (skills: `content-import`, `content-new-item`, `content-add-notes`):
new material lands in `raw/`, is converted/summarized into `docs/`, and its key
information is integrated into the wiki — source summary page created, relevant
entity/concept pages updated, index updated, log appended. One source may touch
many wiki pages; that is expected and desired.

**Query** (skill: `content-query`): to answer questions about the content, read
`wiki/index.md` first, drill into relevant pages, then `docs/` and finally `raw/`
only if needed. Cite pages in answers. Good answers worth keeping should be filed
back into the wiki as new pages — explorations compound.

**Lint**: periodically health-check the wiki: contradictions between pages, stale
claims, orphan pages with no inbound links, concepts mentioned but lacking a page,
missing cross-references, broken links. Fix mechanical issues directly; raise
content-level issues to the user. Log the pass.

## Greenfield vs brownfield

- **Greenfield**: empty project. Use `content-new` to define the first topics; grow
  `docs/` and `wiki/` from scratch.
- **Brownfield**: project with existing unstructured content. Use `content-import`
  with the importer subagent to sweep existing files into the structure: originals
  are copied (never moved) into `raw/`, converted to markdown in `docs/`, and
  summarized into the wiki.

## Division of labor

The user curates sources, directs the analysis, and asks the questions. The agent
does everything else: summarizing, cross-referencing, filing, index and log
bookkeeping, and lint. The maintenance burden is the agent's job — that is the
entire point of this system.
