---
name: recruiter-candidate
description: Take one candidate through the loop — intake (assign the next opaque slug, create the hub entity, profile, and screening notes from the application read at intake, add the roster row), score a round (a scorecard against the round's guide, evidence per competency), file a reference check, or archive a rejected/withdrawn candidate. Use when the user says "add candidate", "new application", "score <slug> on <round>", "references for <slug>", "reject <slug>", "<slug> withdrew", or invokes /recruiter-candidate intake | score <slug> <key> | references <slug> | archive <slug> <reason>.
license: MIT
compatibility: recruiter-interview type (CONTENT.md, Type section); requires docs/candidates/README.md (the roster, created by recruiter-setup) and, for scoring, the round's interview guide; uses the .agentic-cms/scripts toolkit and the entity-candidate, candidate-profile, screening-notes, scorecard, and reference-check templates
allowed-tools: Read Write Edit Grep Glob Bash(.agentic-cms/scripts/*)
---

# recruiter-candidate — intake, score, references, archive

Read `CONTENT.md`'s Type section first if you haven't this session. All
`ac-*` commands live in `.agentic-cms/scripts/`, return JSON — check `"ok"`
after every call. Everything this skill writes about a candidate is **C2**;
the roster, index, and log refer to the candidate by **slug only**.

## Intake mode

1. **Assign the slug:** read `docs/candidates/README.md`'s Roster; the slug
   is the next free `cNNN`. Never derive it from the name.

2. **Read the application at intake** — CV, cover note, application form —
   from the system of record (ATS, mailbox). **Do not copy it into `raw/`**
   (Type section, Retention); summarize what the process needs.

3. **Create the hub and the pages:**
   ```sh
   .agentic-cms/scripts/ac-page new entity-candidate  wiki/entities/cNNN.md               --title "cNNN"                              --classification C2
   .agentic-cms/scripts/ac-page new candidate-profile docs/candidates/cNNN-profile.md     --title "cNNN — Profile"   --topic candidates --classification C2
   .agentic-cms/scripts/ac-page new screening-notes   docs/candidates/cNNN-screening.md   --title "cNNN — Screening" --topic candidates --classification C2
   ```
   The name goes in the hub's and profile's **bodies**, never in a title,
   filename, tag, or one-liner. The profile's evidence section is one line
   per framework competency, marked as claimed. Screening notes are filled
   after the screen; until then the outcome section says "pending". Add
   the slug to `tags` on all three; `refs:` the role entity and the
   competency pages.

4. **Roster and bookkeeping:** add the row (`cNNN | screening | active |
   <date> | [hub](../../wiki/entities/cNNN.md)`), list the two items under
   Items; then
   ```sh
   .agentic-cms/scripts/ac-index add entities wiki/entities/cNNN.md "candidate cNNN — hub (C2)"
   .agentic-cms/scripts/ac-index add topics docs/candidates/cNNN-profile.md "candidate cNNN — profile (C2)"
   .agentic-cms/scripts/ac-index add topics docs/candidates/cNNN-screening.md "candidate cNNN — screening (C2)"
   .agentic-cms/scripts/ac-log append new-item "candidates/cNNN — intake (opaque; pages are C2)"
   .agentic-cms/scripts/ac-page touch docs/candidates/README.md
   .agentic-cms/scripts/ac-index check && .agentic-cms/scripts/ac-links check && .agentic-cms/scripts/ac-classify check wiki/entities/cNNN.md docs/candidates/cNNN-profile.md docs/candidates/cNNN-screening.md
   ```
   After the screen: fill the outcome, set the roster stage to the first
   round key (or `rejected` → archive mode), `ac-page classify` the page
   C2 to re-stamp, log `notes`.

## Score mode (`score <slug> <key>`)

1. **Load:** the round's guide (`docs/interview-design/<key>-guide.md`),
   the candidate's hub and earlier scorecards (open questions carried
   forward), the competency framework's scale.

2. **Create the scorecard** — one per candidate per round:
   ```sh
   .agentic-cms/scripts/ac-page new scorecard docs/candidates/cNNN-<key>-scorecard.md --title "cNNN — <Round name> Scorecard" --topic candidates --classification C2
   ```
   Fill it from the interviewers' notes: one row per competency the guide
   assesses — score on the framework's scale or "no evidence", the
   evidence quoted or closely paraphrased, the interviewer; the moments
   that matter; per-interviewer and conferred recommendation; open
   questions for later rounds. Add the slug and round key to `tags`;
   `refs:` the guide, the hub, each competency page.

3. **Feed forward:** update the hub's strongest/weakest evidence and open
   questions; add the scorecard to each competency page's Related and the
   hub's Relations (`ac-page touch` each, re-rate the hub C2); advance the
   roster stage to the next round key (or `decision` if the stage ends in
   a committee); then
   ```sh
   .agentic-cms/scripts/ac-index add topics docs/candidates/cNNN-<key>-scorecard.md "candidate cNNN — <key> scorecard (C2)"
   .agentic-cms/scripts/ac-log append new-item "candidates/cNNN-<key>-scorecard (opaque; page is C2)"
   .agentic-cms/scripts/ac-index check && .agentic-cms/scripts/ac-links check && .agentic-cms/scripts/ac-classify check docs/candidates/cNNN-<key>-scorecard.md
   ```
   If two interviewers scored the same evidence differently, say so — that
   is `recruiter-round calibrate <key>`'s input.

## References mode (`references <slug>`)

`ac-page new reference-check docs/candidates/cNNN-references.md --title
"cNNN — References" --topic candidates --classification C2`; one subsection
per referee aligned to the competencies still open after the loop; the read
per competency; the outcome. Register (opaque), log, link from the hub,
verify as above. Referees are personal data too — nothing about them leaves
this page.

## Archive mode (`archive <slug> rejected|withdrawn`)

1. For every item of the candidate (`docs/candidates/cNNN-*.md`):
   ```sh
   .agentic-cms/scripts/ac-page archive docs/candidates/<item>.md docs/candidates/archive/<item>.md
   .agentic-cms/scripts/ac-index remove docs/candidates/<item>.md
   .agentic-cms/scripts/ac-index add archived docs/candidates/archive/<item>.md "<original one-liner>"
   ```
   The hub entity stays (entities are not archived) with its status
   updated to the reason and the date; its Relations point at the archived
   paths.
2. Roster: stage `rejected`/`withdrawn`, status `archived`; drop the items
   from the Items list. `ac-log append archive "candidates/cNNN (<reason>)"`.
3. `ac-links check` will flag inbound links from active pages (decision
   records, competency pages) — update them to the archived paths;
   `ac-page touch` each; all three checks clean.
   Archive is the only lifecycle: never delete a candidate's pages. An
   erasure request is handled out of band (system of record, repository
   history) — say so rather than deleting.

## Rules

- Slug everywhere outside the candidate's own pages: titles, filenames,
  tags, roster, index, log. A name in a one-liner is a bleed violation.
- Candidate documents are never stored in `raw/`; summarize at intake.
- Score competencies independently, with evidence, before conferring; "no
  evidence" is unscored, never a middle score.
- Ratchet: raise freely, never lower; floor acks are the user's alone.
