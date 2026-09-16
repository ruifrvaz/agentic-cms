---
name: recruiter-interview
version: 0.1.0
base: "{{VERSION}}"
templates:
  - topic-role.md
  - topic-company.md
  - topic-interview-design.md
  - topic-candidates.md
  - topic-decisions.md
  - job-spec.md
  - competency-framework.md
  - hiring-process.md
  - candidate-preparation-guide.md
  - company-overview.md
  - engineering-culture.md
  - culture-values.md
  - interview-guide.md
  - take-home-assignment.md
  - candidate-profile.md
  - screening-notes.md
  - scorecard.md
  - reference-check.md
  - candidate-offer.md
  - decision-record.md
  - calibration-note.md
  - entity-company.md
  - entity-role.md
  - entity-candidate.md
  - entity-interviewer.md
  - concept-competency.md
  - concept-culture-values.md
  - concept-role-family.md
  - source-role-document.md
skills:
  - recruiter-setup
  - recruiter-round
  - recruiter-candidate
  - recruiter-decision
paths:
  - exercises/
---

# Installed content type

This project runs the **recruiter-interview** content type, installed by
`agentic-cms init --type recruiter-interview`. It layers a domain on top of
the base agentic-cms scaffold: the hiring side's knowledge base for **one
requisition** — define the role and the competencies it is assessed on,
design the interview loop, take candidates through it with evidence-based
scorecards, and reach decisions the whole hiring team can stand behind. The
type is company-, role-, and stack-agnostic by design — nothing it ships
names a company, a role, a technology, or a person; those are
per-requisition content created by the type's skills. It is the
hiring-side counterpart of `candidate-interview` (a sibling type); the two
are maintained independently. See `CONTENT.md`'s
`## Type: recruiter-interview` section for the full schema (topics,
candidate slugs and the bleed rule, retention, the Rounds table contract,
scoring and deciding, take-home assignments, classification defaults).

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

- **Schema** (`CONTENT.md`): a filename exception for candidate-slug-prefixed
  items and dated decision records, the `exercises/` and `TYPE.md` lines in
  the directory map, and this file's companion
  `## Type: recruiter-interview` section. No new page type — every page
  stays a base `doc`/`entity`/`concept`/`source`.
- **Templates** (`.agentic-cms/templates/`, listed in this file's
  frontmatter): topic READMEs, role/company/interview-design/candidates/
  decisions items, entities, concepts, and a source template, for the five
  type topics (`role`, `company`, `interview-design`, `candidates`,
  `decisions`). Every template uses only the placeholders `ac-page new`
  substitutes (`{{TITLE}}`, `{{TOPIC}}`, `{{RAW_PATH}}`, `{{STATUS}}`,
  `{{CLASSIFICATION}}`, `{{CLASSIFIED_HASH}}`, `{{DATE}}`).
- **Skills** (`.claude/skills/`): `recruiter-setup` (bootstrap a
  requisition), `recruiter-round` (design a round; author a take-home;
  calibrate interviewers), `recruiter-candidate` (intake; score a round;
  reference check; archive), `recruiter-decision` (compare candidates and
  decide; make an offer).
- **Agent instructions** (`CLAUDE.md`): a second managed block,
  `<!-- agentic-cms:type:recruiter-interview:begin/end -->`, appended after
  the base block.
- **Code** (`exercises/001-rate-limiter/`): one illustrative take-home —
  runnable code outside the CMS layers, with a candidate-facing `README.md`
  and an interviewer-only `GRADING.md` that never leaves this repository.
- **Content**: none. `docs/` and `wiki/` ship empty (base `index.md` and
  `log.md` only) — this type's skills and templates create every page per
  requisition, not the installer.

## Candidates are referred to by slug, never by name, outside their own pages

Every candidate gets an opaque slug at intake (`c001`, `c002`, …) — the
only candidate identifier that appears in filenames, the roster, and
`wiki/index.md`/`wiki/log.md`. The candidate's name and everything personal
live only inside their own C2 pages. Candidate documents (CVs, take-home
submissions, correspondence) are never stored in `raw/` — they stay in the
system of record; the agent reads them at intake and summarizes what the
process needs.

## Acceptance shape

`recruiter-setup`, then `recruiter-round`, `recruiter-candidate`, and
`recruiter-decision` run against a real requisition, must be able to
produce a tree with: the five topics populated, a Rounds table in
`interview-design` with one interviewer guide per round, one candidate
roster row plus slug-prefixed profile/screening/scorecard/reference pages
per candidate, and at least one decision record comparing candidates
against the competency framework. Shape, not byte-equality, is the
acceptance bar — every requisition's actual content differs. Every
roster/index/log line for a candidate must name only their slug, never
their name.
