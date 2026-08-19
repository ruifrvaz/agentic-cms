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
    archive/          optional: retired items (status: archived)
wiki/
  index.md            catalog of every page: link + one-line summary, by category
  log.md              append-only journal of operations
  entities/           pages about specific things (people, products, systems, orgs)
  concepts/           pages about ideas, patterns, themes
  sources/            one summary page per ingested raw source
.agentic-cms/
  scripts/            deterministic CLI toolkit (ac-page, ac-index, ac-log,
                      ac-links, ac-inventory, ac-search, ac-classify) —
                      JSON on stdout; contract in scripts/README.md
  hooks/pre-commit    git pre-commit gate (classification check on staged
                      content); wired into .git/hooks/ by `init`
  templates/          markdown templates for every page type
  VERSION             scaffold version installed by agentic-cms
.claude/
  skills/             content-* skills (workflows)
  agents/             researcher, importer, exporter subagents
  settings.json       Claude Code PostToolUse hook (classification check)
.codex/
  hooks.json          Codex PostToolUse hook (classification check)
```

## The toolkit

Mechanical operations MUST go through `.agentic-cms/scripts/` instead of hand-editing:
creating pages from templates (`ac-page new`), reading/updating frontmatter
(`ac-page meta|touch`), index maintenance (`ac-index add|remove|check`), logging
(`ac-log append`), link checking (`ac-links check`), inventory (`ac-inventory`),
and term search (`ac-search`). Every command prints one JSON object; check `"ok"`
before proceeding, and finish write operations with `ac-index check` /
`ac-links check` returning `"clean": true`. Judgment work — summarizing,
synthesis, wording — stays with the agent. Full contract:
`.agentic-cms/scripts/README.md`.

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
  status: final             # draft | final (default) | archived
  classification: C1        # C0 | C1 (default) | C2 | C3 — see Classification below
  classified-hash: a1b2c3d4e5f6  # stamped by ac-page at rating time; do not hand-edit
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
  then run the normal `content-manage-item` register+log step (index entry, log
  entry) against the new path.
- **Archive**: `status: archived` marks a first-class page as retired — no longer
  active content, but kept retrievable (archive ≠ delete). Archived items live at
  `docs/<topic>/archive/<item>.md`, the mirror of `drafts/`. Unlike a draft, an
  archived page STAYS in `wiki/index.md`, re-filed under the `## Archived`
  section, so `content-query` can still find it and knows it is retired.
  **Archiving**: `ac-page archive docs/<topic>/<item>.md
  docs/<topic>/archive/<item>.md` moves the file, sets `status: archived`, and
  touches `updated:`; then re-file the index entry (`ac-index remove` +
  `ac-index add archived`), append an `archive` log entry, fix any inbound
  links from active pages (`ac-links check` flags them), and drop the item
  from its topic `README.md` list. **Un-archiving** is `ac-page promote` from
  `archive/` back up plus the reverse re-filing and an `unarchive` log entry.
  `raw/` is never archived (it is immutable); `exports/` is never archived
  (derived artifacts are regenerable).
- **Classification**: every page carries a `classification: C0-C3`
  confidentiality rating — see the [Classification](#classification)
  section below for the full scale, the ratchet, floors, acks, the bleed
  rule, and enforcement.
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
  `import`, `research`, `notes`, `lint`, `export`, `archive`, `unarchive`, `classify`.
- **Contradictions**: when new material contradicts an existing page, do not silently
  overwrite. Note the contradiction inline (`> ⚠ Contradicts [X](...) — newer source
  says ...`), update if the newer source is authoritative, and log it.

## Classification

`classification: C0 | C1 | C2 | C3` rates a page's confidentiality on the
standard CIA confidentiality axis:

- **C0 Public** — no harm if published (marketing copy, public docs).
- **C1 Internal** (default; absent field means C1) — members-only, low
  harm (working notes, internal how-tos).
- **C2 Confidential** — need-to-know: financial figures, personal data
  (PII), private strategy.
- **C3 Restricted** — severe harm: verbatim private correspondence,
  legal instruments, anything credential-adjacent (API keys, tokens,
  passwords).

**Rating is the agent's job, at write time**, judged against the scale
above — not a mechanical classifier. Every write-path skill rates before
it writes and passes the rating to `ac-page new --classification`
(default C1 if the content is unremarkable). `ac-page` stamps a
`classified-hash:` alongside the rating — a short hash of the page body
— so later tooling can tell whether the content changed since it was
last rated.

**Ratchet**: an agent may *raise* a page's classification on update — a
new note can make a page more sensitive — but only the **user** may
*lower* one. Misclassifying downward is the dangerous direction;
over-protecting is always the recoverable one. No mechanical path in
this toolkit ever lowers a rating.

**Heuristic floors, not ratings**: `.agentic-cms/scripts/ac-classify`
pattern-matches page bodies for credential-shaped strings (→ implies at
least C3) and PII/financial content — emails, SSN-shaped numbers, and
currency amounts in any enumerated shape: symbol-prefixed (`$100`),
symbol-suffixed (`100€`), ISO-coded (`USD 100`, `100 EUR`), or word-form
(`1,000 dollars`) (→ implies at least C2). A pattern hit is a **floor**:
if the page's current rating is below it, tooling may auto-raise to the
floor, never lower toward it. The floor's contract is **recall** — zero
false negatives over the shapes it enumerates — so patterns are never
narrowed to reduce noise; a false positive is resolved by an explicit
user acknowledgment (below), never by a weaker pattern. The floor is
deliberately advisory-only for full C0–C3 judgment — a regex cannot rate
confidentiality — but authoritative as a floor, because raising is
always safe under the ratchet rule above.

**Floor acknowledgment — a user decision, never an agent's**: when the
**user** judges a floor hit a false positive (say, a public list-price
table tripping the currency floor), they — and only they — may keep the
lower rating with `ac-page classify <path> <level> --ack-floor`. This
stamps a `classification-ack:` bound to the same body hash as
`classified-hash:`; `ac-classify` then reports that floor hit as
non-blocking. **Any change to the page body invalidates the ack** and
the floor violation returns until the user re-reviews. Agents must never
pass `--ack-floor` on their own judgment — only relay an explicit user
instruction — and a plain `ac-page classify` without the flag withdraws
any standing ack. This is the ratchet's escape hatch made mechanical,
with its asymmetry intact: agents raise, only users lower or ack.

**Bleed rule**: `wiki/index.md` one-liners and `wiki/log.md` entries for
a **C2+** page may use only opaque summaries — no figures, no personal
details, no quoted content. The bookkeeping layer inherits the
classification of what it describes unless deliberately written to
avoid doing so.

**Enforcement fires at two moments plus one audit — see
`.agentic-cms/scripts/README.md`'s `ac-classify` contract for the exact
commands**: an agent-lifecycle hook after every `docs/`/`wiki/` write, a
git pre-commit gate on staged content (the last local checkpoint before
a push), and `content-lint`'s existing health sweep as the periodic
catch-all for anything the first two never saw (pre-existing drift,
edits committed with `--no-verify`, or edits made with no hook
installed). Every one of these is a thin caller of `ac-classify` — none
re-implements detection. The pre-commit gate is **delta-scoped**: it
blocks only for violations in the files staged in that commit (plus any
bleed in a staged `wiki/index.md`/`wiki/log.md`, regardless of which
page leaked), and summarizes pre-existing drift elsewhere in the tree as
one non-blocking warning — `content-lint` owns that backlog, so a
brownfield adoption can rate its legacy pages incrementally.

**Multi-instance segregation — where C2+ content physically lives (a
separate vault repo, a two-tier split) — is out of scope for this
schema.** Classification makes content classifiable and bleed-safe
within one CMS instance; segregation is deployment policy, layered on
top by the project if needed.

## Operations

**Ingest** (skills: `content-import`, `content-manage-item`, `content-add-notes`):
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
