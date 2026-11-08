# Phase 4 — Living Features: Design

## New components

```text
internal/migrate/    migration workflow + ceremony state machine
internal/dreams/     dream scheduler + generator (reuses diary machinery)
internal/graveyard/  burial flow, epitaph generator, graveyard page
internal/bonds/      bond derivation (co-tenant, correspondent, shiftmate)
                     + diary cross-reference injector
internal/runs/       benchmark orchestration via agent actions channel,
                     run_record storage, PR detection
internal/growth/     host-profile diffing on ingest → growth records
internal/terrarium/  tokenised read-only public page
internal/keepsake/   adoption certificate + Life Book (printable HTML
                     views designed for Save-as-PDF; no PDF library)
store migrations v4: bond, run_record tables (if not yet present),
                     growth records, terrarium token on familiar/lab
                     settings, life.cause_of_death values completed
```

## Key decisions

- **Identity vs. body is the core modelling rule of this phase:**
  `familiar` = the soul (name, seed, diary, bonds, records);
  `life` = one body (host profile, species, scars earned in that body,
  born/died). Everything in this phase is an application of that split.
  Scars carry over as *remembered* marks: rendered in later lives as
  faint "memory scars" plus kept in the record; stitches are on the
  soul and always fully counted.
- **The ceremony is a UI state machine, not theatre code:** steps
  (farewell → orb → assembly → carried-over summary) are driven by the
  real migration API responses, so the animation cannot claim a step
  the backend hasn't done. Reduced-motion: respect
  `prefers-reduced-motion` with a fade-only variant.
- **Runs reuse the actions channel** from Phase 3: a `benchmark`
  internal action (not operator-configurable, not part of the
  permission matrix) executed by the agent; demo mode short-circuits
  with the familiar's documented cap (Anls 1 Mbps) and jitter.
- **Growth detection** compares each enrolment/ingest host profile
  against `life.host_profile`; RAM/disk changes beyond a small epsilon
  create a growth record and update the life profile. No agent change
  needed beyond Phase 1's profile reporting (add profile to ingest if
  missing — small Phase 4 agent touch, allowed).
- **Terrarium is a separate handler + template set** sharing only the
  renderer and read queries. It must be impossible to reach admin
  routes from the terrarium token path (dedicated middleware; test it).
- **Keepsakes are printable HTML**, styled for A4, using the same SVG
  renderer. "PDF export" = a dedicated print view + a download button
  that triggers the browser's Save-as-PDF. This keeps the no-dependency
  rule and produces genuinely good PDFs. If a server-side PDF library
  is proposed instead, it needs the standard dependency approval.
- **Epitaphs and dreams go through the diary grounding pipeline**
  (context snapshot stored, template fallback, budget counted). No new
  text-generation path is introduced in this phase.
- **Bonds are derived, not declared:** a nightly (and on-event) job in
  the lab recomputes bonds from co-tenancy (same host id in life
  profiles), latency measurements (from run records / agent-reported
  peer pings where available; demo values in demo mode), and activity
  overlap. Diary cross-references are injected into the grounded
  context — a familiar may only mention a bond that exists in the table.

## Demo mode

Every feature gets a demo path: migration ceremony on Sail, a dream
for Alpi on demand, three pre-buried graveyard residents with epitaphs
(e.g. a Minibox lost to `trial expired`), visible bonds (No. 1021 ↔
his six neighbours; Anls ↔ Hermuse at 190 ms letters; Vee ↔ Hermuse
night shift), a runnable Anls (1 Mbps PR jokes write themselves — keep
them grounded in the recorded numbers), and a working terrarium link.
