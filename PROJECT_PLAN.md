# Server Familiar — Project Plan

**Event:** Kiroween 2026 · **Category:** Frankenstein
**Submission deadline:** 5 December 2026, 05:00 WIB
**Repository:** public GitHub repo, built with Kiro, `/.kiro` at the root
**Working title for the lab:** the Laboratory

---

## 1. Why this exists

### Monitoring is a commodity

Checking whether a server is up is a solved problem. Uptime checkers,
SSH dashboards, Prometheus, Grafana, Netdata, Beszel, Uptime Kuma —
dozens of competent tools will tell you a host is down, draw its CPU
graph, and send you an alert. Another dashboard with nicer charts is
not a reason to build anything, and it is not a reason for a judge to
care.

What none of those tools model is that, for the people who run small
fleets, a server is not a row in a table. It has a history. It survived
things. Losing it means something.

### The fleet that inspired it

The builder runs a fleet of roughly **28 servers on a Tailscale
tailnet**. Almost none of them are respectable, permanent machines:

- Trial VPS instances that expire in one to three months, forcing a
  migration to whatever trial, credit, or free tier is available next.
- Unused cloud credits scattered across providers — AWS, boxd.sh,
  free-tier slots half-forgotten.
- A 512 MB Alpine box with half a CPU core.
- A shared LXC container (6 cores / 16 GB / 150 GB, no Docker) shared
  with strangers.
- An AWS Lightsail instance in Jakarta, a Tencent CVM whose egress is
  capped at roughly **1 Mbps**, a big shared Ryzen host, worker VMs
  running AI agents around the clock.

Servers in this fleet are perpetually dying of expired trials, being
reborn on different hardware in a different country, changing size,
changing neighbours, losing and regaining services. The fleet is not
*maintained* so much as *kept alive*.

### KKCS: the proof, and the gap

The builder already wrote the monitoring engine. **KKCS** is his Go
fleet monitor: agents on each host report metrics to a central panel.
It works. It is also, deliberately, a *passive dashboard* — names,
graphs, status lights. It proves the hard technical half of this
project (a small Go agent that runs anywhere, a collector, a panel)
already exists and runs in production on exactly this ragged fleet.

KKCS answers "is it up?" It cannot answer "how is it?", "what has it
been through?", or "what did you lose when it died?" Nobody would feel
anything deleting a dead row from KKCS. That indifference is the gap.

### The category is the architecture

Frankenstein's creature is not a metaphor pasted onto this project —
it is a literal description of the fleet: **a body stitched together
from borrowed, mismatched, previously-owned parts, brought to life,
which then has to live with what it is made of.**

- A familiar's host changes because a trial expired: that is a limb
  replaced.
- Three unrelated services share one small VM: that is a patchwork body.
- A server comes back after an outage: that is a resurrection, and it
  should leave a mark.
- The fleet itself is the creature — assembled from dead providers'
  free tiers and animated by one operator who refuses to let it die.

### The missing layer: a server having a life

Server Familiar adds the layer no monitoring tool has: **continuity of
identity through change, expressed as a life.** Metrics become a body
and a mood. Incidents become scars. Migrations become deaths and
reincarnations. Uptime becomes age. The operator stops administering
inventory and starts keeping creatures alive — which is, honestly,
what running a fleet of expiring trials already feels like at 2 a.m.

> Monitoring tells you a server is down. Server Familiar makes you
> care that it almost died — and helps it live.

---

## 2. What it is

**Server Familiar** turns each server in a fleet into a living
creature — a *familiar* — rendered as a procedurally drawn monster
whose body is its hardware. A tiny Go agent on each host reports real
metrics to a central **lab**. The lab derives each familiar's body,
mood, and behaviour from those metrics, lets familiars perform
self-care under the operator's permission, and writes their lives —
diary entries, dreams, scars, deaths, reincarnations — in their own
voices, always grounded in real telemetry.

The operator's home is the lab UI: a stormy Frankenstein laboratory
where the fleet lives in specimen jars, a Night Log records what the
creatures did and felt, and a graveyard remembers the ones that did
not make it.

### How it differs

| | Typical monitor | KKCS (builder's existing tool) | Server Familiar |
| --- | --- | --- | --- |
| Core question | "Is it up? What are the graphs?" | "Is the fleet up?" | "How is it living? What has it survived?" |
| Identity | Hostname + IP, deleted when decommissioned | Tailscale name, passive row | A named creature with a persistent soul across hardware |
| Migration | Re-register a new host; history resets | New row | **Death + reincarnation**: new body, same diary, scars, and character |
| Incidents | Log lines and alerts | Status change | Permanent **scars and stitches** on the body |
| Alerts | Threshold emails/webhooks | Panel status | Guard behaviour: growls, scratches, barks — with the telemetry attached |
| Benchmark | A number in a table | — | "Take for a run": play that produces a personal record |
| Downtime | Red graph segment | Red row | Sleep, dreams written from its own logs |
| Decommission | Row deleted | Row deleted | **Burial** in the graveyard, with an epitaph |
| Output for others | Status page | Private panel | A public **terrarium**: the fleet as creatures |
| Interface | Dashboard + API | Dashboard | Dashboard + Telegram + **MCP** (ask your fleet how it feels) |

---

## 3. The three pillars

### ALIVE — a body and a mood, derived from real metrics

Every familiar has continuously derived state. Nothing is decorative
randomness; everything maps to telemetry:

| Life signal | Derived from | Meaning |
| --- | --- | --- |
| **Mood** | CPU load, error/incident rate, recent events | Cheerful, anxious, feverish, sleepy, grumpy… |
| **Energy** | Load average vs. capacity, time since boot | Bustling, steady, exhausted, asleep |
| **Hunger** | Disk usage | Full disk = starving/full-bellied discomfort |
| **Fever** | CPU saturation (and temperature where available) | Sustained high CPU = fever, flushed LEDs |
| **Body size** | RAM | 512 MB is a sprite; 16 GB is a brute |
| **Appetite / stride** | Bandwidth, esp. egress cap | A 1 Mbps cap is a permanently slow gait — "Slow Blood" |
| **Age** | Uptime, and total lived time across lives | Young, adult, elder |
| **Stage** | Total lived time + milestones | Stage I (juvenile) → II (adult) → III (elder) |
| **Scars** | Incidents (OOM, disk-full, crash, attack) | Permanent marks, rendered on the body |
| **Stitches** | Resurrections | **+1 stitch every time it comes back from the dead** |

Mood uses smoothing/hysteresis so a creature does not flicker between
states on every sample, and every derived value can be traced back to
the samples that produced it.

### ACTS — self-care, under permission

Familiars can *do* things, through a strict allow-list executed by the
agent: restart a known service, clear caches, run a health ritual
(log rotation, package cache cleanup, a connectivity check). The
operator chooses per familiar / per action:

- **Auto** — the familiar heals itself and reports what it did.
- **Ask-first** — it requests permission (in the lab and via Telegram)
  and only acts on approval.

Every action is logged (`action_log`), reversible where the action
type allows, and never arbitrary shell: the agent exposes named
actions only, with fixed parameters.

### TALKS — a voice grounded in its own telemetry

Each familiar keeps a **diary** and can be **chatted with**. All
generated text obeys one hard rule: **it may only say things its real
metrics support.** If the diary says "I was feverish all afternoon,"
there are CPU samples proving it. The lab (never the agent) calls a
configurable OpenAI-compatible LLM endpoint with a grounded prompt:
the familiar's character sheet, its recent metrics summary, and its
recent events. Responses are cached, a daily token budget caps cost,
and if no LLM is configured the system falls back to hand-written
templates that are still fully grounded — the diary never goes silent
for lack of an API key.

The fleet is also exposed over **MCP** (Model Context Protocol), so an
AI assistant can be asked "how is my fleet feeling?" and answer from
live familiar state.

---

## 4. The creature system

### Species (14)

A familiar's species follows its hardware class:

| # | Species | Hardware |
| --- | --- | --- |
| 001 | Rackling | Rack server |
| 002 | Tower Brute | Tower PC |
| 003 | Board Sprite | Single-board computer / tiny VPS |
| 004 | Cabinet Colony | Shared rack / LXC container with neighbours |
| 005 | Screenrunner | Monitor-headed worker VM |
| 006 | Patchwork | Multi-service stitched host |
| 007 | Blade Imp | Blade server |
| 008 | Vaultback | NAS / storage box (belly of drive bays) |
| 009 | Routling | Router / switch |
| 010 | Crunchbeast | GPU rig (fan-covered, steaming) |
| 011 | Oldframe | Mainframe |
| 012 | Lapdog | Laptop (four-legged, loyal, sleepy) |
| 013 | Cloudpuff | Pure cloud instance (no chassis; it floats) |
| 014 | Minibox | Mini-PC / NUC |

### Rendering

Creatures are **procedural SVG, rendered server-side** — no image
assets, no sprite sheets. The design language is *"the body IS the
machine"*: chassis bodies, drive-bay teeth, LED eyes, cooling-fan
mouths, network cables as tails and tentacles, antennas as ears,
Frankenstein stitches across panel seams, neck bolts, mismatched
patchwork panels.

Rendering is **deterministic** from a tuple:

```text
render(species, stage I/II/III, palette seed, scar list, mood, size-from-RAM)
```

Same inputs, same creature — so a familiar looks like itself on every
page, in every export, and after every restart. Variety is therefore
combinatorial: **14 species × 3 stages × LED palettes × scar,
stitch, and accessory layers earned from real incidents** means the
50th server enrolled in a fleet still gets a face no other familiar
has. Mood changes expression and LED colour; scars and stitches, once
earned, render permanently for the rest of that life.

---

## 5. Feature set

### Core

- Agent enrolment with a one-line install script and a fleet token.
- Live fleet view: creature cards with SSE-streamed vitals.
- Mood / energy / hunger engine with incident detection.
- Diary, Night Log event feed, familiar detail view (character sheet,
  vitals, scars, record book).
- **Demo mode:** a simulated 6-familiar fleet, so judges (and anyone
  else) can run the whole product with zero servers.

### The nine living features

1. **Reincarnation on migration.** When a server is replaced (trial
   expired, new provider), the operator runs the migration ceremony:
   the old body dies, a **stitched soul-orb** carries the familiar's
   diary, scars, character, and records to a new body whose species
   matches the *new* hardware. A past-lives timeline records every
   body it has worn.
2. **Guard senses.** Security telemetry as animal instinct: an SSH
   login from an unknown IP makes the familiar **growl** (Telegram
   alert with IP/geo attached); a burst of failed logins leaves
   **scratches on its door** (counted, shown, scar-eligible); a port
   scan makes it **bark**. The alert is the instinct; the evidence is
   one click behind it.
3. **Take for a run.** A speedtest / benchmark framed as play: press
   the button and the creature runs (animated) while the agent
   measures. Results become **personal records** in its record book —
   including an honourably consistent 1 Mbps.
4. **Dreams.** While a familiar sleeps (idle or down), the lab composes
   a short dream **from its own recent logs and diary**. Downtime
   becomes the most poetic content in the product instead of a red gap.
5. **The Graveyard.** Retired and truly-dead familiars are buried, not
   deleted: a mossy headstone page with dates, age, lives lived,
   stitch count, cause of death ("trial expired"), and an epitaph.
6. **Fleet bonds.** Familiars relate: co-tenants on one host gossip in
   each other's diaries; inter-server latency is framed as letter
   delivery time ("his letters take 190 ms to arrive"); shared night
   shifts form working friendships.
7. **Evolution from real upgrades.** Add RAM or disk and the body
   **grows** — shoulders widen, a new drive bay appears. Downgrade and
   it visibly shrinks. Infrastructure changes are recorded as growth
   records in its life story.
8. **Public terrarium.** A tokenised, read-only public page showing
   the fleet as creatures — a status page people actually enjoy
   opening, suitable for a studio to show clients.
9. **Adoption certificate + Life Book.** Enrolment issues a framed
   **adoption certificate** with the lab's seal. At any time, a
   familiar's whole life — diary, dreams, scars, records, lives — can
   be exported as a printable **Life Book (PDF)**.

---

## 6. Architecture

```text
 ┌─ Server A ──────────┐   ┌─ Server B ──────────┐        ┌─ Demo fleet ─┐
 │ agent (Go, static)  │   │ agent (Go, static)  │        │ simulator    │
 │ /proc sampler, 10 s │   │ action executor     │        │ (in-lab)     │
 └──────────┬──────────┘   └──────────┬──────────┘        └──────┬───────┘
            │  HTTPS POST metrics      │  HTTPS POST metrics        │
            │  (fleet token)           │  (polls for approved       │
            │                          │   actions)                 │
            ▼                          ▼                            ▼
 ┌───────────────────────────────────────────────────────────────────┐
 │ THE LAB (Go monolith, one binary)                                 │
 │  ingest API · mood engine · incident/scar engine · renderer (SVG) │
 │  web UI (server-rendered HTML + vanilla JS + SSE)                 │
 │  diary/dream/chat generator ──► LLM (OpenAI-compatible, optional) │
 │  Telegram bot (approvals, guard alerts, death notices)            │
 │  MCP server (list_familiars, how_is, fleet_mood)                  │
 │  storage: SQLite, single file, WAL mode                           │
 └───────────────────────────────────────────────────────────────────┘
```

Key properties:

- **The agent is a fork of the KKCS monitoring engine** — the proven
  sampler and reporting loop — reduced to the standard library, sampling
  `/proc` every 10 seconds. It is a single static binary (CGO off) that
  runs on a 512 MB Alpine/musl box as happily as on a large host.
- **Agents are outbound-only.** They initiate every connection; the lab
  never dials into a server. Approved actions are picked up by the
  agent on its next poll. There is no inbound port and no remote shell.
- **Actions are an allow-list, never commands.** The agent exposes
  named actions with fixed behaviour (restart a configured service,
  clear caches, health ritual). No arbitrary command string ever
  crosses the wire.
- **The lab is a monolith on purpose.** One Go binary + one SQLite
  file is the whole deployment. The intended users run fleets of tiny
  borrowed machines; the lab must not demand a database server, Redis,
  or a container platform (one showcase host has no Docker at all).
- **The LLM lives only in the lab**, behind a small internal interface
  with: grounded prompt construction (character sheet + metric summary
  + events), response caching keyed by input state, a configurable
  daily token budget, and a template fallback that keeps diary and
  dreams working with no key configured.
- **Telegram** is the out-of-band nervous system: action approvals,
  guard-sense alerts, death and resurrection notices.
- **MCP** exposes the fleet to AI assistants as tools, reusing the same
  state the UI shows.

### Security model

- Fleet token on ingest; per-familiar token after enrolment.
- Terrarium links are unguessable tokens, read-only, revocable.
- Secrets only via environment variables; none in the repo, ever.
- Action allow-list with per-familiar permission mode (Auto / Ask).
- Every action and approval written to `action_log` — the audit trail
  is part of the creature's life story.

---

## 7. Tech stack

| Layer | Choice | Why |
| --- | --- | --- |
| Agent | Go, stdlib only, CGO off, static binary | Runs on 512 MB Alpine/musl and every other odd host in the fleet; fork of the proven KKCS sampler |
| Lab | Go monolith, `net/http` stdlib, no framework | One binary to deploy anywhere; no dependency roulette before a deadline |
| Storage | SQLite, single file, WAL mode | Zero-ops, backup = copy one file, plenty for a personal/small-studio fleet |
| Web UI | Server-rendered HTML, vanilla JS, SSE | Live vitals without a frontend framework, build step, or npm; pages work on any device |
| Creatures | Procedural server-side SVG | Deterministic from state; no assets to draw, license, or ship |
| LLM | Configurable OpenAI-compatible endpoint, lab-side only | Any provider or self-hosted model; caching + budget caps; template fallback |
| Notifications / approvals | Telegram Bot API | The operator already lives there; approvals need to reach a phone |
| AI-assistant interface | MCP server (Go) | Fleet state becomes tools: `list_familiars`, `how_is`, `fleet_mood` |
| Tests / build | `go test ./...`, `go vet`, `gofmt` | Boring, fast, built in |
| Process docs | `/.kiro` steering + specs | Kiro-native workflow; also a hackathon submission requirement |

Explicitly **not** used: a frontend framework, a JS build step, npm,
an external database, a message queue, containers as a requirement.
Every one of those would exclude some real host in the target fleet.

---

## 8. Data model

SQLite, eight core tables (migrations applied by the lab at startup):

| Table | Holds |
| --- | --- |
| `familiar` | Identity that survives hardware: id, name, character seed / palette seed, permission mode, status (alive / sleeping / dead / buried), terrarium token, created_at |
| `life` | One row per body a familiar has worn: familiar_id, life number, host profile (provider, region, CPU, RAM, disk), species, born_at, died_at, cause of death |
| `metric_sample` | Raw telemetry: familiar_id, life_id, ts, cpu_pct, mem_used/total, disk_used/total, load1, uptime_s, net rates; pruned on a retention schedule |
| `incident` | Detected events: kind (oom, disk_full, crash, unknown_login, failed_login_burst, port_scan, downtime…), severity, ts, detail JSON, scar/stitch awarded |
| `diary_entry` | Diary, dreams, and Night Log narration: familiar_id, kind (diary / dream / nightlog), body text, grounded-metric snapshot, source (llm / template), ts |
| `bond` | Relationships: familiar_a, familiar_b, kind (co_tenant, correspondent, shiftmate), strength, last_interaction_at |
| `run_record` | "Take for a run" results: familiar_id, ts, measured throughput/latency, is_personal_record |
| `action_log` | Every self-care action: familiar_id, action, mode (auto / approved / denied), requested_by, result, ts — the audit trail |

Derived state (current mood, energy, hunger, stage, scar list) is
computed from samples + incidents and cached on the `familiar`/`life`
rows; it can always be recomputed from the raw tables.

---

## 9. Build plan

Deadline-anchored: everything ships by **5 Dec 2026, 05:00 WIB**.
Phases map 1:1 to the specs in `/.kiro/specs/`.

| Phase | Scope | Exit condition |
| --- | --- | --- |
| **0 — Kit & scaffold** | This kit in the repo root, Go module, `cmd/agent`, `cmd/lab`, README skeleton, `.gitignore` | Repo clones, `go build ./...` passes on empty binaries |
| **1 — Skeleton** | Agent sampling `/proc` every 10 s → lab ingest; SQLite (`familiar`, `metric_sample`); enrolment + install script; renderer v1 (species from host profile, size from RAM); creature cards + SSE vitals; **demo mode** with the simulated 6-familiar fleet | A stranger can run demo mode and watch six live creatures; a real host can enrol |
| **2 — Alive** | Mood/energy/hunger engine, growth stages, incident → scar/stitch engine, resurrection detection (+1 stitch), diary (templates first, LLM optional), Night Log, status filters | Creatures visibly react to load; kill and revive one and it returns with a new stitch and a diary entry about dying |
| **3 — Acts & Talks** | Action allow-list + Auto/Ask permission modes; Telegram approvals and alerts; guard senses (growl / scratches / bark); grounded familiar chat; MCP server (`list_familiars`, `how_is`, `fleet_mood`) | Approve a restart from a phone; simulate a failed-login burst and get scratched-door evidence; ask an MCP client how the fleet feels |
| **4 — Living features** | Reincarnation ceremony + past-lives timeline; dreams; graveyard + epitaphs; fleet bonds; Take-for-a-run + personal records; growth records; public terrarium; adoption certificate + Life Book PDF | All nine living features demonstrable in demo mode |
| **5 — Submission** | Polish, README (incl. `/.kiro` and how Kiro was used), demo-mode hardening, Devpost text, demo video, final verification checklist | Submission checklist green before the deadline, with buffer |

Rules of the build: phases finish in order; build + tests pass before
a phase is called done; demo mode must demonstrate *every* feature,
because demo mode is what the judges will actually run.

### Kiro credit budget

**4,000 Kiro credits** are available (Kiroween grant, expiring
31 Dec 2026). Planned allocation:

| Share | For |
| --- | --- |
| 40% (~1,600) | Phases 1–2: scaffold, agent, lab, renderer, mood engine — the foundation gets the depth |
| 30% (~1,200) | Phases 3–4: actions, Telegram, MCP, and the nine living features |
| 20% (~800) | Polish, tests, edge cases, README and submission assets |
| 10% (~400) | **Reserve** — held back for deadline-week surprises; never planned away |

The `/.kiro` structure that steers the spend: `steering/product.md`,
`steering/tech-rules.md`, and `steering/creature-design.md` are loaded
on every interaction (so Kiro never re-derives the rules), and each
phase is a spec (`requirements.md` / `design.md` / `tasks.md`) so work
proceeds task-by-task instead of in open-ended chat.

---

## 10. Demo script (3 minutes, for the judges)

1. **0:00 — The lab.** Open demo mode. Six creatures in specimen jars,
   vitals streaming. "This is my real fleet. Every monitoring tool
   shows you graphs of these machines. This one shows you *them*."
2. **0:30 — Alive.** Push one familiar's CPU: its LEDs flush, mood
   turns feverish; fill its disk and hunger spikes. Introduce **Anls**,
   whose chest meter is pinned at **1 Mbps** — "Slow Blood. He has
   never run faster. He is proud of it."
3. **1:10 — Acts & guard senses.** Simulate a burst of failed SSH
   logins: the familiar growls, a Telegram alert lands with the
   attacker's IP, and scratch marks appear on its door. Then approve a
   self-care restart from the Telegram prompt.
4. **1:50 — Death and reincarnation.** Kill a familiar's host
   ("his trial expired — this happens to my fleet every month").
   Run the migration ceremony: the body fades to a stitched soul-orb,
   a new body assembles — different species, matching the new
   hardware — with the same diary, the same scars, **plus one stitch**.
   Show its past-lives timeline.
5. **2:30 — The soft stuff that lands.** A dream it wrote from its own
   logs while it was down. The graveyard, and one good epitaph. Take
   Anls for a run: 1 Mbps, a new personal record by 0 Mbps.
6. **2:55 — Close.** Show the MCP answer to "how is my fleet feeling?"
   and the public terrarium link. Final line below.

---

## 11. Risks & mitigations

| Risk | Mitigation |
| --- | --- |
| Scope: nine living features is a lot for one hackathon | Phases are ordered by demo value; Phase 1+2 alone is a submittable product, 3–4 are additive |
| LLM cost / unavailability during judging | Template fallback is a first-class path, not a stub; caching + daily budget cap; demo mode ships pre-grounded sample text paths |
| Agent must run on hostile little hosts (512 MB, musl, no Docker) | Go stdlib only, CGO off, static binary; tested against the smallest real host in the fleet (Alpi) from Phase 1 |
| Security: an agent that *acts* is an attack surface | Outbound-only polling, named allow-list actions, no shell, per-familiar tokens, Auto/Ask permission modes, full `action_log` audit |
| Judges have no servers to enrol | Demo mode simulates the full fleet and every feature, including migration and death |
| Deadline crunch / credit exhaustion | 10% credit reserve; Phase 5 starts from a "submission checklist", not from hope; README and Devpost text drafted inside Phase 5 tasks, not after |
| Renderer variety becomes a time sink | Deterministic parameterised SVG (species × stage × palette × scars); new species are parameter sets + a body path, not new art pipelines |

## 12. Post-hackathon roadmap

- **More species & stages** as the community's hardware demands;
  a species submission guide for contributors.
- **Smarter guard senses:** IP reputation lookups, impossible-travel
  detection, per-familiar "stranger" learning.
- **Fleet stories:** weekly generated "fleet chronicle" — the lab's
  newsletter of births, deaths, feuds, and records.
- **KKCS convergence:** offer Server Familiar as the living front-end
  for existing KKCS deployments; one agent, two ways of seeing a fleet.
- **Studio use:** public terrariums as client status pages for
  Kode Kecil (the builder's web studio) — clients watch their servers
  live as creatures instead of reading uptime charts.
- **Self-hosting polish:** one-line lab install, backup/restore of the
  single SQLite file, terrarium theming.

---

*Monitoring tells you a server is down. Server Familiar makes you care that it almost died — and helps it live.*
