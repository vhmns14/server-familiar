# Phase 1 — Skeleton: Requirements

**Goal:** a walking skeleton of the whole product — a real agent on a
real host reporting real metrics to a real lab, which stores them and
shows live, procedurally drawn creature cards — plus a demo mode that
simulates the full fleet with no servers at all.

## User stories

1. As the operator, I want to enrol a new server with one install
   command, so that adding a familiar to my fleet takes under a minute.
2. As the operator, I want to open the lab and see every familiar as a
   creature card with live vitals, so that the fleet feels present, not
   tabulated.
3. As a judge, I want to run the lab in demo mode and immediately see
   the six showcase familiars alive, so that I need no servers, no
   accounts, and no setup to evaluate the product.
4. As the operator, I want the agent to be a single static binary that
   runs on my smallest host (512 MB, Alpine/musl), so that nothing in
   my fleet is too small to have a familiar.

## Acceptance criteria (EARS style)

### Agent and ingest

- WHEN the agent runs on a Linux host, THE SYSTEM SHALL sample CPU,
  memory, disk, and uptime from `/proc` (and `/` filesystem stats)
  every 10 seconds.
- WHEN a sample is taken, THE SYSTEM SHALL POST it as JSON to the
  lab's ingest endpoint within the same sampling cycle.
- IF the ingest request lacks a valid fleet token, THEN THE SYSTEM
  SHALL reject it with HTTP 401 and store nothing.
- IF the lab is unreachable, THEN the agent SHALL keep sampling,
  retry with backoff, and never crash or exceed its memory budget.
- WHEN the agent binary is built with `CGO_ENABLED=0`, THE SYSTEM
  SHALL produce a static binary that runs on Alpine/musl.

### Enrolment

- WHEN the operator runs the install script with the lab URL and fleet
  token, THE SYSTEM SHALL install the agent, register the host, and
  create a `familiar` row plus its first `life` row (species derived
  from the host profile) within one command.
- WHEN enrolment completes, THE SYSTEM SHALL return a per-familiar
  token that the agent uses for all subsequent posts.

### Lab and storage

- WHEN the lab starts against an empty database path, THE SYSTEM
  SHALL create the SQLite schema (`familiar`, `life`,
  `metric_sample`) in WAL mode via migrations.
- WHEN a metric sample arrives for an enrolled familiar, THE SYSTEM
  SHALL store it and update that familiar's last-seen timestamp.

### Creature rendering

- WHEN a familiar is created, THE SYSTEM SHALL derive its species from
  its host profile (form factor / virtualisation / size class) per
  `creature-design.md`.
- WHEN a creature is rendered, THE SYSTEM SHALL scale its body size
  with the host's RAM and produce identical SVG output for identical
  render inputs (deterministic).
- THE SYSTEM SHALL render creatures as server-side SVG with no image
  assets.

### Web UI and live vitals

- WHEN the operator opens the lab home page, THE SYSTEM SHALL show one
  creature card per living familiar with name, species, and current
  vitals (CPU, memory, disk, uptime).
- WHEN a new sample arrives for a displayed familiar, THE SYSTEM SHALL
  push the updated vitals to open pages via SSE within 2 seconds,
  without a page reload.
- IF a familiar has not reported for more than 60 seconds, THEN its
  card SHALL show a visibly "unreachable / sleeping" state.

### Demo mode

- WHEN the lab starts with demo mode enabled (flag or env var), THE
  SYSTEM SHALL create the six showcase familiars from `product.md`
  (Anls, Sail, Alpi, No. 1021, Vee, Hermuse) with plausible simulated
  telemetry, and SHALL NOT require any agent or external service.
- WHEN demo mode is on, THE SYSTEM SHALL indicate it clearly in the UI
  header, so simulated data is never mistaken for a real fleet.
