---
name: candidate-interview
version: 0.1.0
base: "{{VERSION}}"
templates:
  - topic-company.md
  - topic-candidate.md
  - topic-interview-prep.md
  - topic-interview-practice.md
  - topic-offer.md
  - topic-poc.md
  - company-overview.md
  - engineering-culture.md
  - culture-values.md
  - candidate-cv.md
  - public-footprint.md
  - role-fit-analysis.md
  - hiring-process.md
  - role-research.md
  - candidate-preparation-guide.md
  - round-plan.md
  - round-debrief.md
  - offer-outcome.md
  - poc-plan.md
  - poc-learning-log.md
  - rehearsal-interview.md
  - rehearsal-debrief.md
  - rehearsal-refresher.md
  - rehearsal-exam.md
  - entity-company.md
  - entity-candidate.md
  - entity-flagship-project.md
  - entity-poc.md
  - entity-panelist.md
  - concept-culture-values.md
  - concept-role-family.md
  - source-candidate-cv.md
  - source-preparation-guide.md
  - source-recruiter-email.md
skills:
  - interview-setup
  - interview-round
  - interview-rehearsal
  - interview-refresher
paths:
  - exercises/
---

# Installed content type

This project runs the **candidate-interview** content type, installed by
`agentic-cms init --type candidate-interview`. It layers a domain on top of
the base agentic-cms scaffold: a personal knowledge base supporting one
job-interview recruitment process end to end, from the candidate's side. The
type is company-, role-, and stack-agnostic by design — nothing it ships
names a company, a role, or a technology; those are per-engagement content
created by the type's skills. See `CONTENT.md`'s `## Type: candidate-interview`
section for the full schema (topics, the Rounds table contract, the practice
loop, coding exercises, classification defaults).

The frontmatter block above is this type's machine-readable manifest —
`name`, `version`, the `agentic-cms` `base` version this install was stamped
from, and the paths it owns beyond the base scaffold (`templates` under
`.agentic-cms/templates/`, `skills` under `.claude/skills/`, plus any other
owned `paths`). `init`/`update` always refresh these on every run, the same
framework-owned treatment the base scaffold's own skills and templates get.

## What the type adds on top of the base scaffold

Base scaffold paths (`.agentic-cms/scripts/`, the five base templates, the
nine `content-*` skills, the three subagents) are unchanged and
byte-identical to a typeless install. On top, this type adds:

- **Schema** (`CONTENT.md`): the `rehearsal` page type in the frontmatter
  enum, the `interview-practice` numbered-chain filename exception, the
  `exercises/` and `TYPE.md` lines in the directory map, and this file's
  companion `## Type: candidate-interview` section.
- **Templates** (`.agentic-cms/templates/`, listed in this file's
  frontmatter): topic READMEs, docs items, the `rehearsal-*` chain,
  entities, concepts, and sources for the six type topics (`company`,
  `candidate`, `interview-prep`, `interview-practice`, `offer`, `poc`).
  Every template uses only the placeholders `ac-page new` substitutes
  (`{{TITLE}}`, `{{TOPIC}}`, `{{RAW_PATH}}`, `{{STATUS}}`,
  `{{CLASSIFICATION}}`, `{{CLASSIFIED_HASH}}`, `{{DATE}}`).
- **Skills** (`.claude/skills/`): `interview-setup` (bootstrap an
  engagement), `interview-round` (plan a round; debrief it afterwards),
  `interview-rehearsal` (mock round), `interview-refresher` (gap course +
  exam).
- **Agent instructions** (`CLAUDE.md`): a second managed block,
  `<!-- agentic-cms:type:candidate-interview:begin/end -->`, appended after
  the base block.
- **Code** (`exercises/001-rate-limiter/`): one illustrative coding drill —
  runnable code outside the CMS layers (see the Type section for why).
- **Content**: none. `docs/` and `wiki/` ship empty (base `index.md` and
  `log.md` only) — this type's skills and templates create every page per
  engagement, not the installer.

## Acceptance shape

`interview-setup`, then `interview-round`, `interview-rehearsal`, and
`interview-refresher` run against a real engagement, must be able to produce
a tree with: the six topics populated, one `round-plan`/`round-debrief` pair
per interview round linked from the `interview-prep` Rounds table, a
`interview-practice` rehearsal chain (`NNN-interview.md` /
`NNN-debrief.md`, optionally `NNN-refresher.md` / `NNN-exam.md`), panelist
entities backing each round's panel, and (if a proof-of-concept was built)
`poc-plan.md` / `poc-learning-log.md`. Shape, not byte-equality, is the
acceptance bar — every engagement's actual content differs.
