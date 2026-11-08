# Phase 2 — Alive: Tasks

Work in order. Build + tests green after each group.

## 1. Schema v2

- [ ] 1.1 Migrations: `incident`, `diary_entry` tables; derived-state
      cache + scars/stitches columns on `life`/`familiar`.
- [ ] 1.2 Store functions for incidents, diary entries, scar/stitch
      updates, Night Log queries.

## 2. Mood engine

- [ ] 2.1 `internal/life.Derive` as a pure function with the mood enum,
      energy, hunger, and evidence summary.
- [ ] 2.2 Hysteresis/smoothing rules per design.md.
- [ ] 2.3 Table-driven tests: fever episode, disk filling, idle →
      sleepy, single-spike does NOT flip mood.
- [ ] 2.4 Wire into ingest + 30 s liveness/decay ticker; publish mood
      changes over SSE; renderer consumes mood (expression + LED colour).

## 3. Stages

- [ ] 3.1 Stage calculator (lived time + milestones) + growth events to
      the Night Log.
- [ ] 3.2 Renderer stage variants for the six demo species; golden tests
      updated.

## 4. Incidents, scars, stitches

- [ ] 4.1 Detectors: OOM, disk-full, reboot, downtime (per design table).
- [ ] 4.2 Scar awarding (idempotent, canonical marks in the renderer,
      permanent rendering) + scar rendering for all scar kinds.
- [ ] 4.3 Resurrection detection: +1 stitch exactly once per qualifying
      downtime; stitch seam placement in the renderer that stays tidy at
      high counts (test with 23 stitches — Hermuse).
- [ ] 4.4 Tests: simulated streams for each detector; blip < threshold
      earns nothing.

## 5. Diary

- [ ] 5.1 Grounded context builder + stored snapshot per entry.
- [ ] 5.2 Template engine with per-species sentence banks; slots filled
      only from context (test: no context value → no claim in output).
- [ ] 5.3 LLM client (OpenAI-compatible, env-configured), response
      cache, daily token budget counter, automatic template fallback.
- [ ] 5.4 Scheduler (day-part + event-triggered, rate-limited).

## 6. UI

- [ ] 6.1 Night Log page/section (reverse-chronological, live via SSE).
- [ ] 6.2 Status filters on the grid (All / Awake / Sleeping / Healing /
      Dead).
- [ ] 6.3 Character sheet on the detail page: mood evidence, scars with
      dates/causes, stitch count, diary list.

## 7. Demo mode

- [ ] 7.1 Simulator scripts: fever episode, disk creep, rotating
      death + resurrection (compressed timeline, labelled in demo UI).
- [ ] 7.2 Verify end-to-end in demo mode: within a few minutes a judge
      sees a mood change, a Night Log entry, and a diary entry appear.

## Phase exit check

- [ ] E.1 Build + tests green; renderer golden files reviewed.
- [ ] E.2 Live demo: kill a real agent for > threshold, restart it,
      familiar returns with +1 stitch and a diary entry about it.
- [ ] E.3 Summarise for the user and **wait for approval before Phase 3.**
