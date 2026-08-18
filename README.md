# agentic-cms

An agentic content management system built purely on markdown files, templates,
skills, and subagents — no database, no server. A single Go binary installs a
thin scaffolding on top of any project folder; from then on, your coding agent
does the content work: importing raw documents, organizing them into a
structured knowledge base, researching gaps, and exporting deliverables.

Heavily inspired by [Andrej Karpathy's LLM Wiki](https://gist.github.com/karpathy/442a6bf555914893e9891c11519de94f):
an agent-maintained, compounding wiki that sits between you and your raw sources.
Knowledge is compiled once and kept current — not re-derived on every question.

## Features

- **One installer** — a single Go binary (or the `curl | bash` bootstrap, no Go
  toolchain required) installs a thin scaffold: schema, templates, skills, and
  subagents.
- **Markdown is the database** — no database, no server; everything lives in
  git-versioned markdown, diffable and greppable.
- **Three-layer content model** — immutable `raw/` sources, agent-organized
  `docs/`, and an agent-owned synthesis `wiki/`. See [How it works](#how-it-works).
- **Deterministic core, judgment at the edges** — the `.agentic-cms/scripts/`
  toolkit owns every mechanical operation (page creation, frontmatter,
  indexing, logging, link checking, search) as JSON-in/JSON-out scripts;
  skills use their own judgment only for summarizing and synthesis.
- **`wiki/index.md` is the retrieval system** — at moderate scale (hundreds of
  pages), an agent reading the index beats embedding RAG infrastructure.
- **`wiki/log.md` is the audit trail** — every content operation is appended;
  `grep "^## \[" wiki/log.md | tail -5`.
- **Immutability by convention, enforced by schema** — agents are instructed
  never to touch `raw/`; lint catches drift.
- **One detector, many thin callers** — confidentiality classification
  (`classification: C0-C3`) is enforced at four points (an agent hook after
  every write, a git pre-commit gate, each write-path skill's verify step,
  and `content-lint`'s health sweep), but all four are thin wrappers around
  a single `ac-classify` engine. None re-implements detection; they only
  differ in what they do with the result (block, warn, auto-fix, report).

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

## Compatibility

Currently supported:

| Platform | Status |
|----------|--------|
| Claude Code | ✅ Supported |
| GitHub Copilot | Planned |
| OpenAI Codex | Planned |

## Getting Started

**Prerequisite:** none — the `curl | bash` install needs no Go toolchain. Linux
only for now.

**Install:**

```sh
curl -fsSL https://raw.githubusercontent.com/ruifrvaz/agentic-cms/main/install.sh | bash
```

Downloads the latest release binary into `~/.local/bin`. Pin a version with
`AGENTIC_CMS_VERSION=vX.Y.Z`.

For Go users or contributors building from source:

```sh
go install github.com/ruifrvaz/agentic-cms@latest
# or from a clone:
make build && sudo make install
```

**Initialize:**

```sh
cd your-project
agentic-cms init
```

**Build something:**

1. Open your coding agent (Claude Code today; Codex/Copilot planned) in the project
2. Read `CONTENT.md` — the schema your agent follows
3. Greenfield: run the `content-new` skill to start a topic
4. Brownfield: drop sources into `raw/` (or point at an existing folder) and
   run `content-import`

## Commands

**CLI:**

| Command | Description |
|---------|-------------|
| `agentic-cms init [dir]` | Install the scaffolding into `dir` (default: `.`) |
| `agentic-cms update` | Fetch the latest release, replace the binary, and re-run `init` in the project directory |
| `agentic-cms version` | Print the installed version |
| `agentic-cms help` | Show usage |

**Skills** (invoke from your coding agent):

| Skill | Purpose |
|---|---|
| `content-new` | Start a new topic in `docs/` |
| `content-manage-item` | Add a content item to a topic; also owns the item lifecycle — drafts (`status: draft`), promotion, and archiving (`status: archived`) |
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

### Reinstallation and Updates

Running `agentic-cms init` on an existing installation is safe:

- **Non-destructive** — existing files are never overwritten
- **Idempotent** — re-running doesn't duplicate managed content
- **Merges `CLAUDE.md`** — appends a managed `<!-- agentic-cms -->` block
  instead of replacing an existing file
- **Picks up new scaffold files** — run again after upgrading, or use
  `agentic-cms update`, which fetches the latest release, replaces the
  binary, and re-runs `init` in the project directory

## What gets installed

```
your-project/
├── CONTENT.md              ← the schema (read this first)
├── CLAUDE.md               ← created, or extended with a managed block
├── raw/                    ← immutable sources (+ assets/)
├── docs/                   ← organized topical markdown
├── wiki/                   ← index.md, log.md, entities/, concepts/, sources/
├── .agentic-cms/
│   ├── scripts/            ← deterministic ac-* toolkit
│   ├── hooks/pre-commit    ← classification gate; init wires it into .git/hooks/
│   └── templates/          ← page templates
├── .claude/
│   ├── skills/content-*/   ← the nine skills above
│   ├── agents/             ← researcher, importer, exporter
│   └── settings.json       ← PostToolUse hook (classification check)
└── .codex/
    └── hooks.json          ← PostToolUse hook (classification check)
```

Conversion tooling (pandoc, `markitdown`, `python-pptx`) is installed on demand by
the agent when importing/exporting — the scaffold itself is pure markdown.

## Documentation

- **[CONTENT.md](scaffold/tree/CONTENT.md)** — the schema template installed
  into your project; read the installed copy first
- **[CHANGELOG.md](CHANGELOG.md)** — release history

## Repository layout

```
install.sh         curl | bash bootstrap: downloads the latest release binary
                   into ~/.local/bin, no Go toolchain required
main.go            CLI (init, update, version, help)
update.go          self-update: fetch latest GitHub release, replace the
                   binary, re-run init in the project directory
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

## Roadmap

- Codex and Copilot compatibility (AGENTS.md generation). Partial start:
  Claude Code and Codex both get a classification-check hook today;
  Copilot CLI's hook config is deferred (no documented blocking capability
  as of this writing) — its skill-tail fallback still enforces classification
  there, just without in-turn hook feedback.
- Scaffold diffing for `agentic-cms update` — it currently re-runs the
  non-destructive `init` (adds new files, never touches existing ones); it
  can't yet offer to update scaffold files whose *content* changed upstream,
  nor clean up renamed ones (upgraded installs keep the stale
  `content-new-item/` skill dir alongside its replacement `content-manage-item/`)
- Optional MCP: local search over the wiki (qmd-style) for large content bases

## License

MIT
