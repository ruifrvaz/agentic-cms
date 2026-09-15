---
name: interview-rehearsal
description: Run a mock interview round — the assistant plays the panel in character, runs the round's clock, then breaks character for a scored debrief. Use when the user says "rehearse", "mock interview", "drill the interview", or invokes /interview-rehearsal. First argument selects the round by its key from docs/interview-prep/README.md's Rounds table (default: the next round without a debrief, or the last round if all are debriefed); a further argument selects a segment/axis named in that round's plan doc, or "full" (default).
compatibility: candidate-interview type (CONTENT.md, Type section); requires docs/interview-prep/README.md's Rounds table, the per-round plan/debrief pages it links, and wiki/entities/ persona pages (created by interview-setup and interview-round); filing uses the .agentic-cms/scripts toolkit and the rehearsal-interview / rehearsal-debrief templates
---

# interview-rehearsal — mock interview rounds

Rehearsal harness for whatever hiring-process rounds this KB is tracking. The
assistant conducts the round **in character as that round's panel**, one
question per turn, then breaks character and debriefs against the KB's
prepared anchors. This skill is a generic mechanism — it does not know the
rounds, personas, or rubric itself. All of that lives in content the skill
reads at run time:

- **Rounds, order, and format** — `docs/interview-prep/README.md`'s Rounds
  table: one row per round with at minimum a `key` (the argument value),
  `name`, `duration`, `format` summary, and links to that round's `plan` and
  `debrief` pages.
- **Panel personas** — each round's plan page has a `## Panel` section
  linking `wiki/entities/*.md` pages, one per panelist, each with a one-line
  role/style descriptor. Play those personas as written; do not invent traits
  the entity page doesn't support.
- **Anchors/rubric** — each round's plan page has the scoring anchors the
  debrief step grades against (structure expected, what "hit" looks like).
- **Company culture rubric** — if the company publishes a values/culture
  statement, it lives at `docs/company/culture-values.md`; load it for
  behavioral grading if that round's anchors reference it.

If any of this content is missing or a round's plan page has no `## Panel`
section, **stop and say so** rather than improvising rounds or personas.

Classification: this file contains no C2 content — the personas and question
material live in the C2 pages it loads at run time. Rehearsal transcripts are
C2-sensitive; they stay in-session unless explicitly filed via step 5.

## Steps

1. **Load grounding (read in full, every invocation).**
   - `docs/interview-prep/README.md` — the Rounds table and any general
     process notes.
   - `docs/company/culture-values.md`, if it exists and the target round's
     anchors reference it.
   - `docs/interview-practice/` — the topic README and the **latest**
     `NNN-debrief.md`; carry forward its flagged weaknesses and re-test them.
     If the latest session also has a graded `NNN-exam.md` (see the
     `interview-refresher` skill), prioritize its **still-open** gaps and
     spot-check the "closed" ones with one question each.
   - The target round's **plan page** (linked from the Rounds table): panel,
     segments/axes, predicted questions, anchors.
   - Any prior round's **debrief page** the plan page cross-links as
     "lessons to simulate regardless of round" (format pivots, panel
     dynamics that generalize).
   - Any artifacts the round's plan page points at as the candidate's own
     work-in-progress (e.g. a PoC under `docs/poc/` and its `wiki/entities/`
     page) — these get probed hardest on honest state vs. built-vs-tasked.

2. **Resolve the round and mode** from the arguments:
   - No argument → the first Rounds-table row with no `debrief` page yet, or
     the last row if every round already has one.
   - A round `key` → that row. Unknown key → stop and list the valid keys
     from the table.
   - `full` (default mode): the whole round on a clock proportional to its
     `duration`, covering every segment/axis the plan page defines.
   - A segment/axis name from the plan page → that part only, at full depth.
   - A quoted topic (e.g. `"caching strategy"`): single-question deep drill
     with immediate debrief after the answer.
   - **Coding-round mechanics** (when the round's format includes a live
     coding drill): it runs against a real `exercises/NNN-*/` module (per
     CONTENT.md's convention). Pick the next unattempted exercise — or, if
     none remains, author a fresh one first (broken implementation + full
     test suite, verified against a reference solution before handing over,
     reference never shown). Set the time budget the plan page specifies.
     The user works in their own terminal and narrates; the panel personas
     react to what the user reports — probing methodology, pressing on clock
     awareness, never handing over fixes. The debrief afterwards reviews the
     actual diff. A per-attempt debrief page also belongs in the CMS per
     CONTENT.md (exercise folders themselves stay outside the CMS layers).
   - **Format-pivot rule**: in `full` mode, deliberately deviate from the
     plan's predicted structure at least once — an unexpected framing, a
     metaphor hiding a technical model, or an axis entered through an
     anecdote. Overfitting to the plan is the failure mode being drilled
     against. If a prior round's debrief recorded a real format pivot, echo
     that kind of deviation here too.

3. **Conduct in character.** Rules of play:
   - One voice per panelist from the plan page's `## Panel` section, clearly
     labeled by the name/label the entity page uses. Where a persona's page
     is thin, default to a plausible generic seat for that segment (e.g. one
     seat carrying depth, one observing and landing few but pointed
     questions) rather than inventing specifics the page doesn't support.
   - **One question per turn.** Wait for the user's answer before continuing.
   - **Anti-overfit mix**: draw roughly 70% of questions from the round's
     plan and 30% as fresh variations consistent with the personas — never
     run the plan's list verbatim in order, and never reveal mid-interview
     which questions were predicted.
   - **Behave like a time-boxed panel**: probe vague or unquantified answers
     once ("what was the measurable result?"), politely cut off rambling
     ("in the interest of time —"), and follow up on weak spots the plan
     page's anchors flag as likely.
   - Track the clock in turns (state segment/axis transitions).
   - **No coaching mid-interview.** Stay in character until the debrief.

4. **Debrief (break character explicitly).** For each answer, score against
   the round plan's anchors: **hit / partial / missed**, with one concrete
   improvement each. Always check: structured-answer discipline (e.g. STAR)
   with a quantified result where relevant; tightness (no rambling); the
   company's stated values connected to lived experience, not recited, where
   that rubric applies. Apply whatever round-specific checks the plan page
   defines (technical depth, design-conversation anchors, honest-state
   discipline on the candidate's own artifacts, etc.). End with: top 2
   strengths, top 2 gaps, and what to drill before the next rehearsal.

5. **File the session** (on user confirmation) into `docs/interview-practice/`
   as the next numbered pair, per CONTENT.md's Type section:
   - `NNN-interview.md` — the **verbatim** transcript (panel turns and candidate
     answers word for word, candidate wording unedited), created via
     `ac-page new rehearsal-interview docs/interview-practice/NNN-interview.md --title "NNN — Interview (<round>, <mode>)" --topic interview-practice --classification C2`.
   - `NNN-debrief.md` — the debrief **verbatim** (scores table, the moments
     that matter, strengths/gaps/drills), via
     `ac-page new rehearsal-debrief docs/interview-practice/NNN-debrief.md --title "NNN — Debrief" --topic interview-practice --classification C2`.
   - Add the round key to both pages' `tags`.
   - Cross-link the pair via `refs:`, list both in the topic README, register
     both in `wiki/index.md` (opaque one-liners — bleed rule), append log
     entries, and finish with `ac-index check` / `ac-links check` /
     `ac-classify check` all clean.

## Rules

- Never lower the difficulty because the user built the question plan — the
  value is in the 30% they didn't predict and in follow-up pressure.
- Match whatever language constraint the real round has (state it if the
  plan page specifies one, e.g. English-only).
- If the user answers with "skip", move on and mark that question **missed —
  skipped** in the debrief.
- If grounding pages are missing or the panel composition has changed since
  the pages were written, stop and say so rather than improvising personas.
