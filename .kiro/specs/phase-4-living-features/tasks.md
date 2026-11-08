# Phase 4 — Living Features: Tasks

Work in order. Build + tests green after each group.

## 1. Reincarnation

- [ ] 1.1 `internal/migrate`: close life / open life workflow, species
      from new profile, +1 stitch, carried-over summary data.
- [ ] 1.2 Ceremony UI (farewell → soul-orb → assembly → summary),
      driven by real API steps; reduced-motion variant.
- [ ] 1.3 Past-lives timeline on the detail page.
- [ ] 1.4 Demo: Sail's Frankfurt → Jakarta ceremony on a button.
- [ ] 1.5 Test: diary/scars/records/bonds identical before and after
      migration; only the life changed.

## 2. Dreams

- [ ] 2.1 Dream scheduler for asleep/down familiars; generator via the
      diary pipeline (kind=`dream`, snapshot stored); skip when no
      material.
- [ ] 2.2 Dreams shown on the detail page + Night Log with moon marker.
- [ ] 2.3 Demo: generate Alpi's dream on demand.

## 3. Graveyard

- [ ] 3.1 Burial flow (retire / die without migration) + epitaph
      generator (grounded).
- [ ] 3.2 Graveyard page: headstones with full life summary; buried
      familiars remain readable, not deletable.
- [ ] 3.3 Demo: three pre-buried residents with epitaphs.

## 4. Bonds

- [ ] 4.1 `internal/bonds` derivation job (co_tenant, correspondent,
      shiftmate) + `bond` store functions.
- [ ] 4.2 Bonds section on detail pages; diary context cross-references
      (only for bonds that exist — test this).
- [ ] 4.3 Demo bonds: No. 1021's neighbours, Anls↔Hermuse 190 ms
      letters, Vee↔Hermuse shiftmates.

## 5. Take for a run

- [ ] 5.1 `benchmark` internal action on the agent + lab orchestration;
      run animation in the UI; `run_record` storage + PR detection.
- [ ] 5.2 Record book on the detail page.
- [ ] 5.3 Demo: Anls runs at his 1 Mbps cap; PR logic exercised.

## 6. Growth records

- [ ] 6.1 Host-profile diffing on ingest/enrolment → growth records;
      renderer reflects new size/bays.
- [ ] 6.2 Test: simulated RAM upgrade grows the render; downgrade
      shrinks it; both recorded.

## 7. Public terrarium

- [ ] 7.1 Token generation/revocation; read-only terrarium handler +
      templates; privacy review (no controls, no addresses, no private
      diary).
- [ ] 7.2 Test: old token 404s after regeneration; terrarium cannot
      reach admin routes.

## 8. Keepsakes

- [ ] 8.1 Adoption certificate print view (issued at enrolment,
      re-viewable any time) with lab seal.
- [ ] 8.2 Life Book print view: renders across lives, diary, dreams,
      scars, stitches, records, bonds, epitaph if buried.
- [ ] 8.3 Verify print CSS on A4; Save-as-PDF output checked visually.

## Phase exit check

- [ ] E.1 Build + tests green.
- [ ] E.2 Full demo-mode walkthrough of all nine living features with
      no real servers involved.
- [ ] E.3 Summarise for the user and **wait for approval before Phase 5.**
