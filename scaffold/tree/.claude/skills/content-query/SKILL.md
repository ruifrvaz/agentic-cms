---
name: content-query
description: Answer a directed question over the entire knowledge base — search wiki, docs, and raw for all relevant information and synthesize a cited answer (RAG-style, but against the compiled wiki). Use for "what do we know about X", "compare X and Y", "what does the content base say about Z".
license: MIT
compatibility: Requires bash and python3 (uses the .agentic-cms/scripts toolkit)
allowed-tools: Read Grep Glob Bash(.agentic-cms/scripts/*)
---

# content-query — directed query over the knowledge base

Read `CONTENT.md` at the project root first if you haven't this session. The wiki
is a compiled, synthesized artifact — query it instead of re-deriving answers from
raw sources. Toolkit contract: `.agentic-cms/scripts/README.md`.

## Retrieval procedure (deterministic first, judgment second)

1. **Candidate set** — run both, union the results:
   ```sh
   .agentic-cms/scripts/ac-index list
   .agentic-cms/scripts/ac-search <term1> <term2> <synonyms> <abbreviations>
   ```
   Take every index entry plausibly relevant plus the top `ac-search` files.
2. **Drill in**: read the candidate pages; follow their `refs:` frontmatter
   (`ac-page meta <file>`) and inline links one hop when they bear on the question.
3. **Recent material**: `.agentic-cms/scripts/ac-log tail 10` — flag anything recently
   ingested but not yet well-integrated.
4. **Raw as last resort**: only descend into `raw/` when wiki/docs can't answer
   and a listed source plausibly can. If you find substantial un-ingested
   information (`ac-inventory` → `raw_uningested`), suggest `content-import`.

## Answering (judgment)

- Synthesize across pages; cite every claim with a relative link to its page.
- Surface disagreements between pages instead of silently picking a side.
- State what the knowledge base does NOT cover; offer `content-research`.
- Match output form to the question (prose, table, timeline); hand deck/document
  deliverables to `content-export`.

## File the answer back

If the answer involved real synthesis (comparison, analysis, new connection),
offer to file it:

```sh
.agentic-cms/scripts/ac-page new concept wiki/concepts/<slug>.md --title "<Title>"
# write the synthesis into the page, list refs: of every page used
.agentic-cms/scripts/ac-index add concepts wiki/concepts/<slug>.md "<one-line summary>"
.agentic-cms/scripts/ac-log append query "<question>" "filed wiki/concepts/<slug>.md"
.agentic-cms/scripts/ac-index check
```

For simple lookups: `ac-log append query "<question>" "no page filed"`.

## Rules

- Read-only during retrieval; writes only when filing an answer back.
- Carry source confidence and dates through to the answer.
