# Phase 5 — Submission: Requirements

**Goal:** turn a working product into a *submitted* one — polished,
documented, demo-proof — before **5 December 2026, 05:00 WIB**.
Nothing new gets built in this phase except fixes; this phase is
about what the judges see.

Depends on: Phase 4 complete and approved.

## User stories

1. As a judge, I want to clone the repo, run one command, and be
   inside a living demo within two minutes, so that nothing about
   setup stands between me and the idea.
2. As a judge, I want the README to explain what this is, how Kiro was
   used, and where the `/.kiro` folder's steering and specs drove the
   build, so that the Kiroween criteria are easy to score.
3. As the builder, I want a final verification checklist that proves
   the submission is complete *before* the deadline, so that nothing
   is discovered missing at 04:00 on deadline day.

## Acceptance criteria (EARS style)

### Polish

- WHEN a fresh user runs the lab in demo mode, THE SYSTEM SHALL show
  no placeholder text, no broken renders, and no dead links on any
  page reachable from the home screen.
- WHEN any page is loaded, THE SYSTEM SHALL present consistent
  laboratory theming, readable contrast, and working status filters,
  Night Log, detail pages, graveyard, and terrarium.
- IF an LLM or Telegram is unconfigured, THEN THE SYSTEM SHALL present
  the affected features in their designed fallback states (templates /
  lab-only approvals), never as errors.

### README and documentation

- THE README SHALL contain: the pitch and closing line from
  PROJECT_PLAN.md, screenshots or a GIF of the lab, quick start
  (demo mode first, real enrolment second), architecture summary,
  tech stack table, and the feature list.
- THE README SHALL contain a section explaining the `/.kiro` folder:
  what the steering files enforce and how the five phase specs were
  executed with Kiro, including how credits were allocated
  (40 / 30 / 20 / 10 across foundation, features, polish, reserve).
- THE REPOSITORY SHALL keep `/.kiro` at the root, intact, including
  the ticked-off `tasks.md` files as a record of the build.

### Demo-mode hardening

- WHEN the lab runs in demo mode on a clean machine with only the
  built binaries, THE SYSTEM SHALL run the complete 3-minute demo
  script (PROJECT_PLAN.md §10) without manual intervention, seeded
  data repair, or network access.
- WHEN the demo reaches the migration ceremony, guard event, run, and
  MCP beats, THE SYSTEM SHALL perform them in demo mode exactly as
  specified in Phases 3–4.

### Submission assets

- THE SUBMISSION SHALL include a Devpost text draft covering: problem,
  solution, how it was built (with Kiro), the creature system, the
  nine living features, and what is next.
- THE SUBMISSION SHALL include a demo video shot list following the
  3-minute script in PROJECT_PLAN.md §10, with each beat's screen,
  action, and narration line.
- THE FINAL VERIFICATION CHECKLIST (tasks.md §5) SHALL be fully ticked
  and committed before the deadline.
