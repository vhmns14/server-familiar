# Server Familiar

**Every server in your fleet is a living creature, stitched from its own
metrics.**

Monitoring tells you a server is down. Server Familiar makes you care that it
almost died — and helps it live.

A tiny Go agent lives on each of your machines and reports real telemetry. A
central Go *lab* turns every machine into a procedurally drawn monster whose
body **is** its hardware: RAM is body size, CPU is fever, disk is a full
stomach, bandwidth is appetite, uptime is age. Incidents leave permanent
scars. Every resurrection adds a stitch. Familiars perform their own
self-care under permissions you set, keep a diary grounded in their real
metrics, dream while they sleep, are mourned in a graveyard when they die, and
are reincarnated — same soul, new body — when a server migrates.

Built for **Kiroween 2026**, Frankenstein category. The fleet it was built
for is a patchwork of borrowed, expiring, mismatched machines held together by
habit and urgency, which is to say: it was already Frankenstein's creature.
This gives it a heartbeat.

---

## Screenshots

> _Placeholder — the lab, a familiar's lab file, the graveyard, and the public
> terrarium go here once Phase 1 renders them._

## Quick start

> _Placeholder — filled in as Phase 1 lands. The shape it will take:_

```sh
# Try it with no servers at all: the simulated six-familiar fleet.
go run ./cmd/lab -demo
# then open http://localhost:8787
```

```sh
# Run it for real: start the lab, then enrol a host.
LAB_FLEET_TOKEN=$(openssl rand -hex 32) go run ./cmd/lab
curl -fsSL https://your-lab/install.sh | sh -s -- --lab https://your-lab --token <fleet-token>
```

Configuration is environment variables only — see [`.env.example`](.env.example).
Secrets never live in the repository.

## How it is built

| Layer | Choice |
| --- | --- |
| Language | Go, standard library first. One module, two static binaries: `cmd/agent` and `cmd/lab`. |
| Storage | SQLite, one file, WAL mode. Backup is "copy one file". |
| Frontend | Server-rendered HTML, hand-written CSS, vanilla JS, SSE for live vitals. No framework, no build step, no npm. |
| Creatures | Procedural SVG generated server-side, deterministic from the familiar's species, stage, palette seed, scars, stitches, mood, and RAM. No image assets. |
| Brain | Any OpenAI-compatible endpoint, used only for diary, dreams, and chat — always grounded in real telemetry, with caching, a daily budget, and a template fallback. |
| Integrations | Telegram for approvals and alerts; an MCP server so other AI agents can ask the fleet how it feels. |

The agent reads `/proc` directly, runs with CGO off, is outbound-only, and
fits on a 512 MB Alpine host, because nothing in a fleet should be too small
to have a familiar.

## Repository layout

```text
cmd/agent/        the agent binary, one per monitored server
cmd/lab/          the lab binary: UI, API, SSE, MCP
internal/model/   types shared by the agent and the lab
internal/agent/   /proc sampling and the ingest client
internal/store/   SQLite schema, migrations, and queries
internal/ingest/  enrolment and metric ingest handlers
internal/render/  the procedural SVG creature renderer
internal/web/     pages, static assets, and the SSE hub
internal/demo/    the simulated six-familiar fleet
install/          one-command agent install and enrolment
.kiro/            steering and specs — how this was built with Kiro
```

## The plan

[`PROJECT_PLAN.md`](PROJECT_PLAN.md) is the complete build plan: why this
exists, how it differs from ordinary monitoring, the fourteen species, the
nine living features, the architecture, the data model, the build phases, and
the demo script.

The `.kiro/` directory is the build itself made legible: always-loaded
steering for the product, the technical rules, and the creature design
language, plus one spec — requirements, design, tasks — per phase.

## Status

Phase 1 (skeleton) in progress. See
[`.kiro/specs/phase-1-skeleton/tasks.md`](.kiro/specs/phase-1-skeleton/tasks.md).
