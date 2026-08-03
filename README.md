# agentic-cms

An agentic content management system built purely on markdown files, templates,
skills, and subagents — no database, no server, no license. A single Go binary
installs a thin scaffolding on top of any project folder; from then on, your coding
agent (Claude Code today; Codex/Copilot planned) does the content work: importing
raw documents, organizing them into a structured knowledge base, researching gaps,
and exporting deliverables.

Heavily inspired by [Andrej Karpathy's LLM Wiki](https://gist.github.com/karpathy/442a6bf555914893e9891c11519de94f):
an agent-maintained, compounding wiki that sits between you and your raw sources.
Knowledge is compiled once and kept current — not re-derived on every question.

## How it works

Three layers, installed on top of your existing (brownfield) or empty (greenfield)
project:

| Layer | Folder | Ownership |
|---|---|---|
| Immutable sources | `raw/` | You add; agents only read |
| Organized docs | `docs/<topic>/` | Agents write via skills; you review |
| Synthesis wiki | `wiki/` (index, log, entities, concepts, sources) | Agents own entirely |

The contract between you and the agent lives in **`CONTENT.md`** — the schema file
that turns a generic chatbot into a disciplined content maintainer. `CLAUDE.md`
points agents at it automatically.

You work from any GUI that hosts a coding agent (VS Code, Cursor, a terminal) —
the agent runs the skills, you browse the markdown.

## Install

Linux only for now.

```sh
go install github.com/ruifrvaz/agentic-cms@latest
# or from a clone:
make build && sudo make install
```

## Usage

```sh
cd your-project
agentic-cms init
```

`init` is non-destructive and idempotent: existing files are never overwritten, and
an existing `CLAUDE.md` gets a managed `<!-- agentic-cms -->` block appended rather
than being replaced. Run it again after upgrading to pick up new scaffold files —
or run `agentic-cms update`, which does both: it fetches the latest release,
replaces the binary, and re-runs `init` in the current directory if it looks like
an installed project (a `.agentic-cms/` directory is present). `agentic-cms version`
prints what's currently installed.

Then open your agent in the project and use the skills:

| Skill | Purpose |
|---|---|
| `content-new` | Start a new topic in `docs/` |
| `content-new-item` | Add a content item to a topic |
| `content-import` | Ingest PPTX/DOCX/PDF/text into `raw/` → `docs/` → `wiki/`; brownfield sweeps of whole folders |
| `content-research` | Web-research a question and file the findings |
| `content-query` | Answer a directed question over the whole knowledge base, cited |
| `content-add-notes` | Annotate an existing item |
| `content-list` | Inventory, recent activity, drift report |
| `content-lint` | Health check: contradictions, orphans, stale claims, broken links |
| `content-export` | Build a PPTX or DOCX deliverable from the content base |

Subagents (in `.claude/agents/`): **content-researcher** (web research),
**content-importer** (bulk brownfield import), **content-exporter** (PPTX/DOCX
generation). Skills delegate to them for heavy work.

## What gets installed

```
your-project/
├── CONTENT.md              ← the schema (read this first)
├── CLAUDE.md               ← created, or extended with a managed block
├── raw/                    ← immutable sources (+ assets/)
├── docs/                   ← organized topical markdown
├── wiki/                   ← index.md, log.md, entities/, concepts/, sources/
├── .agentic-cms/           ← templates, installed version, deterministic bin/ toolkit
└── .claude/
    ├── skills/content-*/   ← the nine skills above
    └── agents/             ← researcher, importer, exporter
```

Conversion tooling (pandoc, `markitdown`, `python-pptx`) is installed on demand by
the agent when importing/exporting — the scaffold itself is pure markdown.

## Repository layout

```
main.go            CLI (init, update, version, help)
update.go          self-update: fetch latest GitHub release, replace the
                   binary, re-run init in the current directory
scaffold/embed.go  go:embed + non-destructive installer
scaffold/tree/     the scaffolding that gets installed (edit this to change
                   what ships; `all:` embed includes the dot-directories)
```

## Development

```sh
make build       # build the binary (ldflags-versioned from git describe)
make build-all   # cross-compile linux/amd64 + linux/arm64 release assets
make test        # go vet + go unit tests
make smoke-test  # build, then run the installer end-to-end against a temp sandbox
make version     # print the version that would be built
```

`smoke-test` builds the binary, runs the Go test suite, sanity-checks `--version`/
`--help`, then drives `scripts/smoke-test-installer.sh`: a real `agentic-cms init`
run against a `mktemp` sandbox, diffed against `scaffold/tree/` (after resolving
`{{DATE}}`), plus idempotency and brownfield `CLAUDE.md`-merge checks. Set
`KEEP_SMOKE_DIR=1` to keep the sandbox for inspection.

## Design notes

- **Markdown is the database.** Everything is greppable, diffable, and versionable
  with git; Obsidian/VS Code are the UI.
- **`wiki/index.md` is the retrieval system.** At moderate scale (~hundreds of
  pages) an agent reading the index beats embedding RAG infrastructure.
- **`wiki/log.md` is the audit trail.** `grep "^## \[" wiki/log.md | tail -5`.
- **Immutability by convention, enforced by schema.** Skills and agents are
  instructed never to touch `raw/`; lint catches drift.
- **Deterministic core, judgment at the edges.** `.agentic-cms/bin/` (the `ac-*`
  toolkit) owns every mechanical operation — page creation, frontmatter,
  `wiki/index.md`, `wiki/log.md`, link checking, search — as JSON-in/JSON-out
  bash+python3 scripts. Skills call it for structure and use their own judgment
  only for summarizing, synthesis, and wording.
- **Agent-agnostic core.** Skills/agents are markdown with YAML frontmatter;
  Codex/Copilot support means adding their config surface (e.g. `AGENTS.md`)
  around the same skill bodies.

## Roadmap

- Codex and Copilot compatibility (AGENTS.md generation)
- Scaffold diffing for `agentic-cms update` — it currently re-runs the
  non-destructive `init` (adds new files, never touches existing ones); it
  can't yet offer to update scaffold files whose *content* changed upstream
- Optional MCP: local search over the wiki (qmd-style) for large content bases
- macOS/Windows targets (the update mechanism is already platform-general;
  only the release build matrix needs to grow)

## License

MIT
