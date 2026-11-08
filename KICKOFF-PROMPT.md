# Kickoff Prompt — paste this entire file as your first message to Kiro

---

You are building **Server Familiar** with me for **Kiroween 2026**,
**Frankenstein category**. The submission deadline is **5 December 2026,
05:00 WIB**. The repository will be public on GitHub, and the `/.kiro`
folder at the repo root is part of the submission — keep it intact.

## Read first, before writing any code

1. `PROJECT_PLAN.md` at the repo root — the full plan: why this exists,
   the three pillars, the creature system, the nine living features,
   architecture, tech stack, data model, build phases, and demo script.
2. Every file in `/.kiro/steering/`:
   - `product.md` — what the product is, the showcase familiars, tone rules
   - `tech-rules.md` — hard technical rules you must not break
   - `creature-design.md` — how creatures are designed and rendered
3. The Phase 1 spec in `/.kiro/specs/phase-1-skeleton/`:
   `requirements.md`, `design.md`, and `tasks.md`.

Follow the specs in `/.kiro/specs/` in order, starting with
`phase-1-skeleton`. Do not start a later phase until I approve it.

## What Server Familiar is, in one paragraph

A fleet of servers, reimagined as living creatures in a Frankenstein
laboratory. A tiny Go agent on each server reports real metrics; a Go
"lab" server turns each machine into a procedurally drawn monster whose
body *is* its hardware — mood from load, hunger from disk, scars from
incidents, a new stitch for every resurrection. Familiars act (self-care
actions under my permission), talk (a diary and chat grounded in their
real metrics), die, are buried, and are reincarnated when a server
migrates to a new machine. It is a monitoring tool whose output is a
life, not a chart.

## Your first job: scaffold the repo

- Initialise a Go module (`module server-familiar`, current stable Go).
- Layout:
  - `cmd/agent/` — the agent binary that runs on each monitored server
  - `cmd/lab/` — the lab (web server + API) binary
  - `internal/` packages as the Phase 1 design specifies (metrics,
    store, render, web, demo, and similar — keep packages small and
    single-purpose)
- A `README.md` skeleton: title, one-paragraph pitch, screenshots
  placeholder, quick start placeholder, link to `PROJECT_PLAN.md`.
- A `.gitignore` covering: built binaries (`/agent`, `/lab`,
  `*.exe`), the SQLite database (`*.db`, `*.db-wal`, `*.db-shm`),
  `.env`, editor files, and OS junk.
- Keep `START-HERE.md`, `KICKOFF-PROMPT.md`, `PROJECT_PLAN.md`, and the
  whole `/.kiro` folder exactly where they are.

## Then implement Phase 1 fully

Per `/.kiro/specs/phase-1-skeleton/`:

- **Agent skeleton:** a static Go binary that samples CPU, RAM, disk,
  and uptime from `/proc` every 10 seconds and POSTs them as JSON to the
  lab with a fleet token. It must run comfortably on a 512 MB Alpine
  (musl) host — CGO off, standard library only.
- **Lab:** a Go monolith (standard library `net/http`, no framework)
  with SQLite storage (single file, WAL mode), a familiar enrolment
  endpoint, and an install script that enrols a new host.
- **Creature renderer v1:** server-side procedural SVG. Species comes
  from the host profile, body size scales with RAM. No image assets.
- **Web UI:** server-rendered HTML + vanilla JS, creature cards for the
  fleet, and **live vitals over SSE** (Server-Sent Events). No frontend
  framework, no build step, no npm.
- **Demo mode:** a flag / environment variable that runs the lab with a
  simulated 6-familiar fleet, so judges need no servers of their own.

## Working rules

- Run `go build ./...` and `go test ./...` after each task group in the
  spec's `tasks.md`, and fix failures before moving on. A phase is not
  done until the build and tests pass.
- Go standard library first. Propose any third-party dependency in chat
  and wait for my approval before adding it.
- Never use a frontend framework or a build step. Server-rendered HTML,
  vanilla JS, SSE.
- Secrets only via environment variables. Never write a real token,
  password, or key into code, fixtures, or commits.
- All code, comments, commit messages, and docs in English.
- Work through `tasks.md` in order; tick tasks off as you complete them.
- **Ask me before starting each new phase.** When Phase 1 is complete,
  summarise what works, what you tested, and what is stubbed, then wait.

## Begin

Begin by summarising your understanding of the product in **5 bullets**,
then scaffold the repo and implement Phase 1.
