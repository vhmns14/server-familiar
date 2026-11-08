# Phase 5 — Submission: Design

## Approach

Phase 5 is a sequence of audits and documents, not features. Each
workstream has a concrete artefact the user can check.

### 1. Polish pass

- Walk every route in demo mode and in real mode; fix in priority
  order: (a) anything broken, (b) anything unstyled/inconsistent,
  (c) wording. Log fixes as small commits; no refactors, no new
  features — feature freeze is a design decision of this phase.
- Screenshot set for the README taken from demo mode at a fixed seed
  so the six showcase familiars look their best (wide shot of the lab,
  one character sheet, the ceremony mid-orb, the graveyard).

### 2. README

Structure (single file, top to bottom):

1. Title + one-paragraph pitch + the closing line as a blockquote.
2. Hero screenshot.
3. "Run the demo" — two commands maximum to a living lab.
4. "Enrol a real server" — install script one-liner.
5. Feature tour (core + nine living features, one line each, linked
   to screenshots).
6. Architecture diagram (the ASCII diagram from PROJECT_PLAN.md §6,
   kept as text so it renders everywhere) + tech stack table.
7. **"Built with Kiro"** — the `/.kiro` explanation: steering files
   and what each enforces; the five specs; task counts executed;
   credit allocation 40/30/20/10; one honest paragraph on what Kiro
   did well and what needed steering.
8. Repository layout, development commands (`go build`, `go test`),
   configuration (env vars table — names only, never values).
9. Roadmap (from PROJECT_PLAN.md §12) + licence.

### 3. Demo-mode hardening

- A scripted smoke run: fresh temp dir → build → `LAB_DEMO=1` →
  assert via HTTP that all key routes return 200 (home, detail,
  night log, graveyard, terrarium, certificate) and the SSE stream
  emits a sample. Add it as `scripts/demo-smoke.sh` and a Go test
  where practical.
- Freeze demo seed data; verify the 3-minute script timings are
  achievable (simulator event schedule checked against §10 beats).

### 4. Devpost text draft

`docs/devpost-submission.md` in the repo: title, tagline, problem,
solution, creature system explainer, nine features as bullets, tech
stack, "built with Kiro" section mirroring the README, links (repo,
terrarium demo if public), and the closing line. Written to be pasted
into Devpost with only link fields to fill.

### 5. Demo video shot list

`docs/demo-video-shotlist.md`: a table per beat of PROJECT_PLAN.md
§10 — timestamp, screen, exact clicks/actions, narration line,
fallback take (what to show if a beat misbehaves live; prefer a
recorded take for the ceremony and the guard event).

### 6. Final verification

The checklist in tasks.md §5 is the contract: repo public, `/.kiro`
at root, clean-clone build + tests + demo smoke all green, README
complete, Devpost fields filled, video uploaded/linked, submission
confirmed on Devpost **before** 5 Dec 2026 05:00 WIB — with the target
of finishing a full day early and treating the last day as buffer.
