# Phase 2 — Alive: Design

## New components

```text
internal/life/       mood engine, stage calculator, incident detector,
                     scar/stitch awarding, resurrection detection
internal/diary/      diary generator: grounded context builder,
                     LLM client (lab-side), template engine, cache
internal/web/        Night Log feed, status filters, character sheet page
store migrations v2: incident, diary_entry tables; derived-state cache
                     columns on familiar/life (mood, energy, hunger,
                     stage, scar list JSON, stitch count)
```

## Mood engine

- Inputs: rolling windows over `metric_sample` (5 min and 1 h), recent
  `incident` rows, time since last sample.
- Outputs: `mood` (enum: content, cheerful, anxious, feverish, sleepy,
  grumpy, unwell), `energy` 0–100, `hunger` 0–100.
- Implementation: pure function `Derive(samples, incidents, now) State`
  in `internal/life`, called on each ingest and on a 30 s ticker for
  decay (sleepiness grows as samples go stale). Pure = table-driven
  tests with synthetic sample streams as fixtures.
- Hysteresis: a state transition requires the target condition to hold
  across the short window, not one sample; fever additionally requires
  sustained CPU ≥ 90% for ≥ 3 minutes.
- Explainability: `Derive` also returns the evidence summary (window
  averages, triggering incidents) stored alongside diary context and
  shown on the character sheet ("Feverish — CPU 94% avg, last 5 min").

## Incident detection

Detectors run on the ingest path over consecutive samples:

| Signal | Detection | Scar kind |
| --- | --- | --- |
| OOM / memory pressure | mem available < 3% sustained, or agent-reported oom flag | crack |
| Disk full | disk ≥ 98% | patch plate |
| Crash / unexpected reboot | uptime decreases between samples without a graceful-stop marker | dent |
| Downtime | sample gap ≥ 5 min (configurable) | — (resurrection instead) |
| Thermal/CPU saturation event | fever ≥ 10 min | scorch |

Scar awarding is idempotent per incident id. Scars are stored on the
`life` row (a life keeps its own scars; a reincarnated familiar's
*character* remembers them — full past-lives handling is Phase 4, but
the diary context already includes scar history).

## Resurrection

On ingest, if `now - last_seen ≥ downtime threshold` and the familiar
was marked down by the liveness ticker, emit `resurrection`: stitch
count +1, Night Log entry, diary prompt context flag `returned_from_dead
with duration`. The liveness ticker (30 s) owns alive → sleeping →
dead transitions so the UI state never depends on a page load.

## Diary generation

```text
context = { character sheet, mood evidence, last 24h events,
            recent samples summary, scar/stitch deltas }
        → LLM (if configured & budget remains) → cache by context hash
        → else template engine (per-species sentence banks,
          slots filled only from context values)
```

- Scheduler: one entry per familiar per day-part (morning/evening) +
  event-triggered entries (resurrection, scar, stage change), rate
  limited so demo mode produces a lively log without spamming.
- Every stored entry keeps its `grounded snapshot` (the context JSON)
  so any sentence can be audited back to telemetry — this is also the
  demo talking point.
- Budget: `LLM_DAILY_TOKEN_BUDGET` env; a counter row per day;
  template fallback engages automatically at the cap.

## Night Log

A single reverse-chronological feed over a `nightlog` view: union of
incidents, diary_entries (kind=nightlog summaries), stage changes,
resurrections. Rendered server-side with the pinned-journal styling
from the mockup direction; SSE appends new entries live.

## Demo mode extensions

The simulator gains scripts: a daily fever episode for Vee, a disk
creep for Hermuse, a scheduled "death" (gap) + resurrection for one
rotating familiar, so judges see scars and stitches appear without
waiting days. Demo timelines are compressed, and the UI labels demo
time compression honestly in demo mode only.
