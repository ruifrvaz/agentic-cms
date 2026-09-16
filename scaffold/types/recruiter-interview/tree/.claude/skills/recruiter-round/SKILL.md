---
name: recruiter-round
description: Design one round of the loop — the interviewer guide (question bank per competency, rubric, timing, interviewer parts), one wiki entity page per interviewer, and the row in the Rounds table; author a take-home assignment as runnable code plus its CMS brief; or, once two or more candidates are scored on a round, run a calibration pass and file a calibration note. Use when the user says "design the <round>", "add a round", "the <round> panel is X and Y", "write the take-home", "calibrate <round>", or invokes /recruiter-round <key> | take-home | calibrate <key>.
license: MIT
compatibility: recruiter-interview type (CONTENT.md, Type section); requires docs/interview-design/README.md with a Rounds table (created by recruiter-setup); uses the .agentic-cms/scripts toolkit and the interview-guide, take-home-assignment, entity-interviewer, and calibration-note templates
allowed-tools: Read Write Edit Grep Glob Bash(.agentic-cms/scripts/*)
---

# recruiter-round — design, take-home, calibrate

Read `CONTENT.md`'s Type section first if you haven't this session. All
`ac-*` commands live in `.agentic-cms/scripts/`, return JSON — check `"ok"`
after every call. Three modes: **design** (default), **take-home**, and
**calibrate**.

## Design mode

1. **Load grounding:** `docs/interview-design/README.md` (the Rounds
   table), `docs/role/competency-framework.md` (scale, anchors, coverage
   map — the round must assess the competencies the map assigns to it),
   `docs/role/hiring-process.md`, `docs/company/culture-values.md` for
   behavioral rounds, and every existing interview guide (so rounds don't
   duplicate questions).

2. **Resolve the round:** the `key` argument, else ask. Keys are kebab-case
   round names, never stage numbers. If the row already has a Guide link,
   stop — extend the guide with `content-add-notes` or run calibrate
   mode. If there is no row yet, add one (key, name, duration, format,
   `_pending_`, `_pending_`) and add the round to the framework's coverage
   map.

3. **Interviewers:** names and roles come from the user. For each colleague
   without a page:
   ```sh
   .agentic-cms/scripts/ac-page new entity-interviewer wiki/entities/<slug>.md --title "<Name>"
   .agentic-cms/scripts/ac-index add entities wiki/entities/<slug>.md "<seat> — <round keys>"
   ```
   Fill the seat and the competencies they score; add the round key to
   `tags`. Interviewer pages are C1: no assessment of the person, ever.

4. **Interview guide:**
   ```sh
   .agentic-cms/scripts/ac-page new interview-guide docs/interview-design/<key>-guide.md --title "<Round name> — Interviewer Guide" --topic interview-design
   ```
   Keep only the sections that fit the format. `## Panel` links every
   interviewer with their part and the competencies they own. For each
   competency the coverage map assigns: a bank of 3-5 questions, probes,
   the level-separating follow-up, anchored descriptors made concrete for
   these questions (from the framework — never invent a new scale), false
   positives. Coding rounds point at `exercises/NNN-*/` and name the
   process signals to watch. Set `refs:` (framework, competency pages,
   interviewer pages); add the round key to `tags`. C1 — but never shared
   with candidates.

5. **Register, cross-link, verify:** fill the Guide and Interviewers cells
   in the Rounds table; `ac-index add topics docs/interview-design/<key>-guide.md
   "<one-liner>"`; `ac-log append new-item "interview-design/<key>-guide"`;
   add the guide to each competency page's Related and each interviewer's
   Relations (`ac-page touch` each edited page); update the
   candidate-preparation guide's roadmap row if it exists; then
   ```sh
   .agentic-cms/scripts/ac-index check && .agentic-cms/scripts/ac-links check && .agentic-cms/scripts/ac-classify check docs/interview-design/<key>-guide.md
   ```
   All must report `"clean": true`.

## Take-home mode

1. **Author the code** under `exercises/NNN-<slug>/` (next free NNN): a
   small, self-contained, stdlib-only module in the role's language with
   an intentionally incomplete/buggy implementation, a full test suite
   that only passes once fixed, a candidate-facing `README.md` (the task,
   the rules, how to run it — nothing about grading), and an
   interviewer-only `GRADING.md` (the fixes, the discriminating item, what
   a strong submission looks like, time expectations). Verify the tests
   fail as intended and pass against a reference solution kept under
   `exercises/NNN-<slug>/reference/` — never in the hand-over bundle.
   Exercises are code, not CMS pages: no frontmatter, no index/log entry.

2. **CMS brief:**
   ```sh
   .agentic-cms/scripts/ac-page new take-home-assignment docs/interview-design/take-home-<slug>.md --title "Take-Home — <name>" --topic interview-design
   ```
   Purpose and competencies (linked), the assignment, the exact hand-over
   bundle contents (README, module, tests — never GRADING.md or
   reference/), the grading rubric per competency. Link the folder by
   path. Register in the topic README's Items, index, log, and verify as
   in design mode. Add a Rounds row (`take-home` key, async format) if the
   loop treats it as a stage.

## Calibrate mode

1. **Gather the evidence:** every scorecard for the round
   (`docs/candidates/*-<key>-scorecard.md`, two or more candidates or two
   or more interviewers on the same candidate). Compare scores given for
   the same competency against the evidence quoted.

2. **File the note:**
   ```sh
   .agentic-cms/scripts/ac-page new calibration-note docs/decisions/calibration-<key>-<n>.md --title "Calibration — <Round name>" --topic decisions --classification C2
   ```
   The divergence (by slug), the root cause, the change made — then make
   that change: edit the guide's anchors/questions and, if the descriptor
   itself was the problem, the competency framework and the competency
   page (`ac-page touch` each; re-rate if content became more sensitive).
   Record whether earlier scorecards get re-read. Register (opaque
   one-liner), log `notes` for every page edited, verify.

## Rules

- Guides derive their rubrics from the framework; a new score scale or an
  unmapped competency is a framework change to raise with the user, not
  something a guide invents.
- Interviewer pages are C1 and never assess the person; guides and briefs
  are C1 and never reach candidates; calibration notes are C2 and refer to
  candidates by slug only.
- GRADING.md and reference/ never leave the repository; the bundle is
  README + module + tests.
- No stage numbers anywhere — the Rounds table's keys are the names.
- Ratchet: raise freely, never lower; floor acks are the user's alone.
