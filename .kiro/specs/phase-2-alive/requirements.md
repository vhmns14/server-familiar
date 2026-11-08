# Phase 2 — Alive: Requirements

**Goal:** familiars stop being pictures with numbers and start being
*alive* — mood, energy, and hunger derived from real metrics; growth
stages; incidents that leave permanent scars; resurrections that add
stitches; a grounded diary; and the Night Log.

Depends on: Phase 1 complete and approved.

## User stories

1. As the operator, I want each familiar's face to show how it feels
   (feverish, sleepy, hungry), so that I read the fleet at a glance.
2. As the operator, I want incidents to leave permanent marks, so that
   a familiar's history is visible on its body.
3. As the operator, I want a familiar that comes back from downtime to
   gain a stitch, so that survival is celebrated, not just logged.
4. As the operator, I want to read each familiar's diary and a fleet
   Night Log, so that I know what happened while I wasn't watching.

## Acceptance criteria (EARS style)

### Mood engine

- WHEN metric samples arrive, THE SYSTEM SHALL derive mood, energy,
  and hunger per familiar from rolling windows of real samples
  (CPU → fever/energy, disk → hunger, incidents → anxiety/grumpiness).
- WHEN a single anomalous sample occurs, THE SYSTEM SHALL NOT flip a
  familiar's mood state (smoothing/hysteresis: state changes require
  sustained conditions, e.g. N consecutive samples or a time-weighted
  threshold).
- WHEN a familiar's mood changes, THE SYSTEM SHALL reflect it in the
  rendered expression and LED colour on the next render, in the lab
  and over SSE.
- THE SYSTEM SHALL be able to explain any derived value: given a
  familiar and a timestamp, it SHALL identify the samples/windows that
  produced its mood, energy, and hunger.

### Growth stages

- WHEN a familiar's total lived time and milestones cross a stage
  threshold, THE SYSTEM SHALL advance its stage (I → II → III) and
  record the growth event in the Night Log.
- WHEN a familiar is rendered, THE SYSTEM SHALL apply stage-appropriate
  proportions/features per creature-design.md.

### Incidents, scars, stitches

- WHEN the detector observes OOM, disk-full, crash/unexpected reboot,
  or extended downtime from the sample stream, THE SYSTEM SHALL record
  an `incident` row with kind, severity, and evidence.
- WHEN an incident of scar-eligible severity is recorded, THE SYSTEM
  SHALL award the corresponding scar (kind → canonical mark) exactly
  once per incident, and the scar SHALL render permanently thereafter.
- WHEN a familiar resumes reporting after a downtime gap ≥ the
  configured threshold, THE SYSTEM SHALL record a resurrection and add
  **exactly +1 stitch** to that familiar.
- IF a familiar merely has a brief network blip below the downtime
  threshold, THEN THE SYSTEM SHALL NOT award a stitch.

### Diary

- WHEN the diary generator runs for a familiar (on schedule and after
  major events), THE SYSTEM SHALL produce a 2–5 sentence entry grounded
  only in that familiar's real recent metrics and events.
- IF an LLM endpoint is configured, THEN THE SYSTEM SHALL generate the
  entry via the LLM with the grounded context, cache the result, and
  count it against the daily token budget.
- IF no LLM is configured (or the budget is exhausted), THEN THE SYSTEM
  SHALL generate the entry from templates using the same grounded data,
  and mark the entry's source as `template`.
- THE SYSTEM SHALL NEVER write a diary claim (fever, death, visitor,
  number) that is absent from the grounded context.

### Night Log and UI

- WHEN any notable event occurs (incident, scar, stitch, stage change,
  mood swing, resurrection, diary written), THE SYSTEM SHALL append it
  to the Night Log feed in chronological order with a link to the
  familiar.
- WHEN the operator opens the lab, THE SYSTEM SHALL offer status
  filters (All / Awake / Sleeping / Healing / Dead) that filter the
  creature grid by familiar status.
- WHEN the operator opens a familiar's detail page, THE SYSTEM SHALL
  show its character sheet: species, stage, mood/energy/hunger, scar
  list with dates and causes, stitch count, and its diary.
