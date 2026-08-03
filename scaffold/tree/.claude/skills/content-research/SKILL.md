---
name: content-research
description: Research a topic on the web and file the findings into the CMS. Use when the user wants to gather external information — "research X", "find out about Y", "what's the latest on Z" — and keep the results in the content base.
license: MIT
compatibility: Requires bash, python3, and web access (uses the .agentic-cms/bin toolkit)
allowed-tools: Read Write Edit Grep Glob WebSearch WebFetch Bash(.agentic-cms/bin/*)
---

# content-research — web research into the content base

Read `CONTENT.md` at the project root first if you haven't this session. Toolkit
contract: `.agentic-cms/bin/README.md` — check `"ok"` after every `ac-*` call.

## Steps

1. **Scope against what's known** (deterministic):
   ```sh
   .agentic-cms/bin/ac-search <topic terms>
   .agentic-cms/bin/ac-index list
   ```
   Read the matching pages; research the gaps, not the settled ground.
2. **Delegate** to the `content-researcher` subagent with: the research question,
   the relevant existing page paths, and instructions to return findings with
   source URLs. If subagents are unavailable, do the web research yourself.
3. **File the findings**:
   - New material:
     ```sh
     .agentic-cms/bin/ac-page new doc docs/<topic>/<item>.md --title "<Title>" --topic <topic>
     .agentic-cms/bin/ac-index add topics docs/<topic>/<item>.md "<summary>"
     ```
     Write the findings in (judgment), URLs listed in `sources:`.
   - Updates to existing pages: edit in place, flag contradictions explicitly,
     then `.agentic-cms/bin/ac-page touch <file>`.
   - Update affected wiki entity/concept pages and cross-links the same way.
4. **Log and verify**:
   ```sh
   .agentic-cms/bin/ac-log append research "<question>" "<pages touched>"
   .agentic-cms/bin/ac-index check && .agentic-cms/bin/ac-links check
   ```
   Both must report `"clean": true`.

## Rules

- Every researched claim carries its source URL.
- Distinguish clearly between what sources say and your own synthesis.
- One-off answers still get filed — explorations compound in the wiki.
