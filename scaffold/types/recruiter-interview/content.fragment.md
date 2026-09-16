<!-- ac-type-anchor: after "  sources/            one summary page per ingested raw source" -->
exercises/            type-specific: take-home assignments as runnable code (see Type section)
<!-- ac-type-anchor: replace "  templates/          markdown templates for every page type" -->
  templates/          markdown templates for every page type — the base
                      kinds (doc, entity, concept, source, topic) plus the
                      type-specific templates listed in the Type section
<!-- ac-type-anchor: after "  VERSION             scaffold version installed by agentic-cms" -->
  TYPE.md             the installed content type and what it owns
<!-- ac-type-anchor: replace "  skills/             content-* skills (workflows)" -->
  skills/             content-* skills (workflows) + type-specific skills
<!-- ac-type-anchor: after "- **Filenames**: kebab-case, descriptive, no dates in names (`transformer-architecture.md`)." -->
  **Exceptions (type-defined)**: items in the `candidates` topic are prefixed
  with the candidate's opaque slug — `cNNN-profile.md`, `cNNN-screening.md`,
  `cNNN-<round-key>-scorecard.md`, `cNNN-references.md`, `cNNN-offer.md` —
  so one candidate's pages sort together and the slug, never the name, is
  what bookkeeping refers to; and decision records in the `decisions` topic
  are dated, `YYYY-MM-DD-<stage>.md`, so a committee's records order
  chronologically. See the Type section.
<!-- ac-type-anchor: before "## Operations" -->
## Type: recruiter-interview

This CMS instance runs the **recruiter-interview** content type: the hiring
side's knowledge base for **one requisition** — define the role and the
competencies it is assessed on, design the interview loop, take candidates
through it with evidence-based scorecards, and reach decisions the whole
hiring team can stand behind. Everything in this section is type-defined on
top of the base schema above; the base schema still applies in full. The
type is company-, role-, and stack-agnostic — nothing below names a company,
a role, a technology, or a person; those are per-requisition content.

### Topics

A fresh instance is empty. `recruiter-setup` creates the five topics below
(each `docs/<topic>/README.md` from its `topic-<topic>` template); every
item lands in exactly one of them.

| Topic | Holds | Item templates |
|---|---|---|
| `role` | the requisition: what is being hired, how it is assessed, how the process runs, what candidates are told | `job-spec`, `competency-framework`, `hiring-process`, `candidate-preparation-guide` |
| `company` | the employer context interviewers and candidates need | `company-overview`, `engineering-culture`, `culture-values` |
| `interview-design` | the Rounds table, one interviewer guide per round, the take-home brief | `interview-guide`, `take-home-assignment` |
| `candidates` | the roster and every per-candidate page | `candidate-profile`, `screening-notes`, `scorecard`, `reference-check`, `candidate-offer` |
| `decisions` | hiring-committee decisions and interviewer calibration | `decision-record`, `calibration-note` |

Wiki templates: entities `entity-company`, `entity-role` (the requisition as
the hub page), `entity-candidate` (one per candidate — C2), and
`entity-interviewer` (one per colleague on a panel — C1); concepts
`concept-competency` (one per competency in the framework, the hub every
scorecard scores against), `concept-culture-values`, and
`concept-role-family`; sources `source-role-document` (a job description,
leveling guide, or loop definition received from the organization). Generic
`doc`/`entity`/`concept`/`source` templates remain available for anything
the type didn't anticipate.

### Candidates, slugs, and the bleed rule

Every candidate gets an **opaque slug** at intake — `c001`, `c002`, … (next
free number; never initials, never a name). The slug is the only candidate
identifier that appears in filenames, the roster, `wiki/index.md`, and
`wiki/log.md`. The candidate's name and everything personal live only inside
their C2 pages (the `entity-candidate` page body and the profile). Their
items in `docs/candidates/` are slug-prefixed (see Conventions) so they sort
together; `wiki/entities/cNNN.md` is the hub that links them all. The
roster is the `candidates` topic README: one row per candidate — slug,
current stage (a round key, `offer`, `rejected`, `withdrawn`), status, last
update, links — and nothing identifying.

### Retention

The base rules apply unchanged, and they shape what this type stores:
`raw/` is immutable and archive is the only lifecycle, so **candidate
documents — CVs, take-home submissions, correspondence — are never stored
in `raw/`.** They stay in the system of record (ATS, mailbox); the agent
reads them at intake and summarizes what the process needs into the
candidate's C2 pages. `raw/` holds the organization's own material only
(job descriptions, leveling guides, loop definitions). Rejected and
withdrawn candidates are archived (`status: archived`, `archive/`), never
deleted; an erasure request is handled out of band on the system of record
and on this repository's history, outside the CMS.

### Rounds and guides

`docs/interview-design/README.md` carries the **Rounds** table — one row per
round of the loop: `key` (kebab-case name, never a stage number), name,
duration, format, the round's interviewer-guide page, and the interviewer
entity pages who sit on it. It is the naming authority for rounds: skills
select a round by key, scorecard filenames embed it, tags use it. Each
`interview-guide` defines what the round assesses (which competencies),
the question bank per competency, the rubric (anchored descriptors per
score level, drawn from the competency framework), timing, and each
interviewer's part. `recruiter-round` creates guides and interviewer pages
and appends rows.

### Scoring and deciding

A `scorecard` is written per candidate per round by the interviewer(s),
against the guide's rubric: one score per competency **with the evidence**
(what the candidate said or did), an overall recommendation, and open
questions for later rounds. A `decision-record` compares candidates at a
stage per competency across their scorecards and records the committee's
call per candidate with its reasoning; a `calibration-note` records where
interviewers scored the same evidence differently and what the guide
should say to close the gap. The competency pages under `wiki/concepts/`
link every scorecard that scores them, so a competency's evidence across
candidates is one hop away.

### Take-home assignments

`exercises/NNN-<slug>/` (repo root) holds a take-home as runnable code —
"complete an unfinished codebase by fixing unit tests": a small,
self-contained, stdlib-only module with an intentionally incomplete/buggy
implementation and a full test suite. What candidates receive is the
module and its candidate-facing `README.md`; `GRADING.md` and, when
present, `reference/` (the reference solution) stay in this repository and
are never part of the hand-over bundle. Exercises are the type's one
deliberate exception to "markdown only": they are **not** CMS pages — no
frontmatter, no index/log entries, no classification — ordinary tracked
code. The `take-home-assignment` page in `interview-design` is the CMS-side
brief (purpose, competencies assessed, time box, grading rubric, what the
bundle contains) that links the folder by path. The language is
per-requisition; the type ships one illustrative exercise, not a library.

### Company stack in the wiki

The type ships **no technology-specific content**. Technologies the role
requires surface from the job spec, the competency framework, and the
interview guides; capture them the normal way — one `wiki/concepts/` page
per general technology or practice a round assesses, one `wiki/entities/`
page per specific product the team uses — linked from the guides and the
competency pages that need them. Skills and templates refer to "the
role's stack" generically and never hardcode a technology.

### Classification defaults for this type

Every per-candidate page — entity, profile, screening notes, scorecards,
reference check, offer — is **C2** (PII and assessment of a real person);
decision records and calibration notes are **C2**; a verbatim offer letter
or signed document is **C3**. Interviewer entities are **C1** (colleagues,
internal). Role and company material is C0/C1; a job spec or offer page
carrying compensation figures trips the C2 currency floor — rate it C2 or
have the user ack. The roster README is C1 because it carries slugs only.
Index and log one-liners for every C2 page name the slug and the page kind
and nothing else.

<!-- ac-type-anchor: before "## Greenfield vs brownfield" -->
**Type operations** (skills: `recruiter-setup`, `recruiter-round`,
`recruiter-candidate`, `recruiter-decision`): bootstrap a requisition, design
rounds and take-homes, take candidates through intake → scoring →
references → offer or archive, and record decisions and calibration — see
the Type section.

<!-- ac-type-anchor: replace "  `docs/` and `wiki/` from scratch." -->
  `docs/` and `wiki/` from scratch. For this type, run `recruiter-setup` first —
  it lays down the five type topics so `content-new` is only needed for
  anything beyond them.
<!-- ac-type-anchor: replace "  summarized into the wiki." -->
  summarized into the wiki. For this type, only the organization's own
  documents are import candidates — never candidate documents (see
  Retention).
