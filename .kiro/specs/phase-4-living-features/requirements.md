# Phase 4 — Living Features: Requirements

**Goal:** the nine living features that separate Server Familiar from
every monitoring tool — the features the demo is remembered for. Core
and Acts/Talks are done; this phase makes the fleet's *lives*
complete: death and reincarnation, dreams, the graveyard, bonds,
runs, growth, the public terrarium, and the keepsakes (adoption
certificate, Life Book).

Depends on: Phase 3 complete and approved.

## Acceptance criteria (EARS style)

### 1. Reincarnation on migration

- WHEN the operator starts a migration for a familiar and enrols its
  new host, THE SYSTEM SHALL close the current `life` (died_at, cause:
  `migrated`), open a new `life` whose species derives from the new
  host profile, and keep the same `familiar` identity, diary, scars
  (as remembered marks), character/palette seed, records, and bonds.
- WHEN a migration completes, THE SYSTEM SHALL award +1 stitch and
  present the ceremony in the UI: old body fading → stitched soul-orb
  → new body assembling, ending with a "what was carried over" summary.
- WHEN the operator opens a familiar, THE SYSTEM SHALL show its
  **past-lives timeline**: every life with host, species, dates, and
  cause of death.
- IN DEMO MODE, THE SYSTEM SHALL run the full ceremony on the showcase
  example (Sail: Minibox, Frankfurt — died of `trial expired` → Tower
  Brute, Jakarta) at the press of a button.

### 2. Dreams

- WHEN a familiar is asleep or down for a sustained period, THE SYSTEM
  SHALL generate a short dream entry (kind=`dream`) composed from its
  own recent logs, diary, and events, via the LLM/template machinery.
- IF a familiar has no recent material, THEN THE SYSTEM SHALL skip the
  dream rather than invent one.

### 3. The Graveyard

- WHEN a familiar is retired or dies without migration, THE SYSTEM
  SHALL bury it: it leaves the living grid and appears in the
  Graveyard with a headstone showing name, species(es), dates, age,
  lives lived, stitch count, cause of death, and an epitaph.
- WHEN a burial occurs, THE SYSTEM SHALL generate the epitaph from the
  familiar's actual life record (grounded; templates acceptable).
- BURIED FAMILIARS SHALL remain fully readable (diary, lives, records)
  and SHALL NOT be deletable from the UI.

### 4. Fleet bonds

- WHEN two familiars share a host (co-tenants), THE SYSTEM SHALL create
  a `bond` of kind `co_tenant` and they SHALL occasionally reference
  each other in diary entries (grounded in the bond + shared events).
- WHEN two familiars' hosts have measured inter-server latency, THE
  SYSTEM SHALL maintain a `correspondent` bond expressing latency as
  letter-delivery time in diary/UI text.
- WHEN familiars are active in the same night hours repeatedly, THE
  SYSTEM SHALL form `shiftmate` bonds.

### 5. Take for a run

- WHEN the operator presses "Take for a run" on a familiar, THE SYSTEM
  SHALL run a benchmark via the agent (throughput/latency probe;
  in demo mode, simulated with the familiar's known caps — Anls: 1 Mbps),
  animate the creature running, and record a `run_record`.
- WHEN a run result exceeds the familiar's previous best, THE SYSTEM
  SHALL mark it a personal record and note it in the Night Log/diary.

### 6. Evolution from upgrades

- WHEN a familiar's host profile changes materially (RAM or disk size
  differs from the current life's recorded profile), THE SYSTEM SHALL
  record a growth record, update the render (size/bays per
  creature-design.md), and note the change in the diary context.
- IF the change is a downgrade, THEN THE SYSTEM SHALL render the
  familiar visibly smaller and record it just as faithfully.

### 7. Public terrarium

- WHEN the operator enables the terrarium, THE SYSTEM SHALL serve a
  read-only public page at an unguessable token URL showing the living
  fleet as creatures with vitals and moods — and NO controls, NO
  diary text marked private, NO host addresses.
- WHEN the operator revokes/regenerates the token, THE SYSTEM SHALL
  stop serving the old URL immediately.

### 8–9. Adoption certificate + Life Book

- WHEN a familiar is enrolled, THE SYSTEM SHALL generate an adoption
  certificate (printable HTML/PDF view): name, species, adoption date,
  host description, palette, and the lab seal.
- WHEN the operator requests a Life Book export, THE SYSTEM SHALL
  produce a printable PDF (via the browser print path / server-side
  HTML→PDF-free approach chosen in design) containing the familiar's
  renders across lives, diary, dreams, scars, stitch count, run
  records, bonds, and graveyard entry if buried.
