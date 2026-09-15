---
name: interview-refresher
description: Close gaps found by interview rehearsals — reads the target debrief in docs/interview-practice/, authors a numbered refresher course targeting each gap, then (once the course is completed) generates and grades a multiple-choice exam from the same gaps. Use when the user says "refresher", "close my gaps", "run the refresher course", "generate the exam", or invokes /interview-refresher. Optional argument: a session number NNN (default: the latest debrief).
compatibility: candidate-interview type (CONTENT.md, Type section); requires docs/interview-practice/ with at least one NNN-debrief.md; uses the .agentic-cms/scripts toolkit and the rehearsal-refresher / rehearsal-exam templates in .agentic-cms/templates/
---

# interview-refresher — gap-targeted course + exam

Second half of the practice loop: `/interview-rehearsal` diagnoses, this skill
remediates. It extends the session's numbered chain in `docs/interview-practice/`
(per CONTENT.md's Type section):
`NNN-interview.md → NNN-debrief.md → NNN-refresher.md → NNN-exam.md`.

Classification: this file holds no C2 content; the gap material lives in the C2
pages it reads and writes. Rate every page written at write time per CONTENT.md.

## Steps

1. **Load grounding:**
   - `docs/interview-practice/README.md` and the target `NNN-debrief.md`
     (argument NNN, else the highest-numbered debrief). No debrief → stop and
     say a rehearsal must run first.
   - The relevant round's plan page (found via `docs/interview-prep/README.md`'s
     Rounds table) — the anchors the gaps were scored against.
   - The wiki concept pages relevant to the named gaps (whatever technology
     or domain each gap is about — see CONTENT.md's "Company stack in the
     wiki" section for where that content should live).
   - If `NNN-refresher.md` already exists, resume at its current state (unit
     checkboxes) instead of re-authoring; if `NNN-exam.md` exists ungraded,
     offer to administer/grade it.

2. **Extract and classify the gaps** from the debrief's scores, key moments,
   and drill list. Classify each as:
   - **knowledge** — missing facts or mental models (examinable by knowledge
     MCQs), or
   - **behavioral** — skills/habits under pressure (examinable by
     situational-judgment MCQs; also drilled live in `/interview-rehearsal`).
   State the gap list with classifications in one short message before writing.

3. **Phase 1 — author the course** as `docs/interview-practice/NNN-refresher.md`
   (`ac-page new rehearsal-refresher docs/interview-practice/NNN-refresher.md
   --title "NNN — Refresher Course" --topic interview-practice
   --classification C2` — C2 is the safe default since gaps derive from C2
   debriefs; add the round key to `tags`). One **numbered unit per
   gap**, each unit containing:
   - **Why** — the debrief moment that exposed it (one line, linked).
   - **Core material** — teach the missing model. Build on the wiki concept
     pages by linking them; if a concept page is thin or missing, extend or
     create it (normal CMS bookkeeping) rather than forking parallel notes
     into the course.
   - **Exercise** — something to actually do: for hands-on gaps prefer real
     drills against a local test environment for the technology in question;
     for behavioral gaps a rehearsable line or micro-drill.
   - **`- [ ] Unit N complete`** checkbox.
   Register, log, verify (index with opaque one-liner per the bleed rule,
   `ac-log append new-item`, all three checks clean), and add the item to the
   topic README's chain listing.

4. **Walk the course** if the user wants it interactive: teach unit by unit,
   answer questions, and check off each `- [ ] Unit N complete` (then
   `ac-page touch` + re-rate) as the user confirms it. Self-study is equally
   valid — the checkboxes are the source of truth either way.

5. **Completion gate:** Phase 2 runs only when every unit checkbox in
   `NNN-refresher.md` is checked. If asked to skip ahead, require an explicit
   user override and record it in the refresher page's Notes.

6. **Phase 2 — generate the exam** as `docs/interview-practice/NNN-exam.md`:
   create the page from the `rehearsal-exam` template (`ac-page new
   rehearsal-exam docs/interview-practice/NNN-exam.md --title "NNN — Exam"
   --topic interview-practice --classification C2`), then fill its Questions
   section and the answer-key comment **exactly in the structure the
   template prescribes** — the template is the only source for exam shape;
   no ad-hoc shapes. Content rules:
   - 3–5 questions per gap; **knowledge** gaps get knowledge MCQs,
     **behavioral** gaps get situational-judgment MCQs (a realistic interview
     moment, four candidate responses, one best).
   - Four options per question, one correct/best; distractors must be
     plausible (adjacent misconceptions, not jokes).
   - Every question maps to its gap; the answer key, per-question rationale,
     and gap map go inside the template's HTML-comment block (invisible on
     render, persisted for grading).
   Register, log, verify as in step 3.

7. **Administer and grade:** the user submits answers in any clear form
   (e.g. `1a 2c 3b ...`). Grade against the key; append a dated result entry
   to the exam page's Notes: overall score, per-gap score, and each gap marked
   **closed** (≥ the template's pass threshold) or **still open**. Still-open
   gaps are what the next `/interview-rehearsal` re-tests — its grounding
   reads this chain.

## Rules

- Numbering follows the sourced debrief: gaps from `NNN-debrief.md` produce
  `NNN-refresher.md` and `NNN-exam.md` — never an unrelated number.
- Don't duplicate wiki knowledge into course units — link and extend the
  concept pages; the course adds sequencing, gap-framing, and exercises.
- The exam tests the gaps, not general trivia: no question without a gap
  mapping.
- Grading is honest: report the score as computed; a failed gap stays "still
  open" — never soften it to closed.
- `raw/` stays untouched; every page write follows CMS bookkeeping (index with
  opaque C2 one-liners, log, checks clean); classification ratchet applies.
