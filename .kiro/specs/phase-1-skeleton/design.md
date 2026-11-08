# Phase 1 — Skeleton: Design

## Components

```text
cmd/agent/main.go        agent entrypoint: config from env/flags, sampler loop
cmd/lab/main.go          lab entrypoint: config, migrations, HTTP server
internal/model/          types shared by both binaries (sample, host profile,
                         species, stage); imports no other internal package
internal/agent/          /proc sampling, ingest client with backoff
internal/store/          SQLite access + migrations (familiar, life, metric_sample)
internal/ingest/         ingest + enrolment HTTP handlers, token checks
internal/render/         deterministic SVG creature renderer v1
internal/web/            pages (server-rendered HTML), static CSS/JS, SSE hub
internal/demo/           simulated fleet: 6 familiars, plausible telemetry generator
install/install.sh       one-command agent install + enrolment
```

## Data flow

1. Agent samples `/proc/stat`, `/proc/meminfo`, `/proc/uptime`, and
   `statfs("/")` every 10 s → builds a sample JSON.
2. Agent POSTs to `POST /api/ingest` with its per-familiar token
   (fleet token only during enrolment). On failure: exponential
   backoff (cap 5 min), drop-oldest in-memory ring of unsent samples —
   memory stays bounded on tiny hosts.
3. Lab validates the token, writes `metric_sample`, updates
   `familiar.last_seen_at`, and publishes the sample to the SSE hub.
4. Browser pages subscribe to `GET /events` (SSE); a small vanilla-JS
   handler patches card vitals in place. Cards are server-rendered on
   first load; SVG comes from `internal/render` inline in the HTML.
5. Demo mode replaces steps 1–2: `internal/demo` runs the same
   pipeline from a simulator goroutine per familiar, writing through
   the same store functions — demo familiars are real rows, so every
   later phase works on them unchanged.

## Key decisions

- **One module, two binaries.** Shared types (sample, host profile)
  live in `internal/`; no agent code depends on lab-only packages, so
  the agent binary stays small.
- **SQLite via a pure-Go driver or `modernc.org/sqlite`-class
  dependency is NOT pre-approved**: CGO is off, so if the stdlib has
  no answer, propose the driver in chat per `tech-rules.md` before
  adding it. (This is the one expected dependency conversation.)
- **Species derivation v1** is deliberately simple: a function from
  host profile (virtualisation type, RAM band, storage type, service
  count) to species per the table in `creature-design.md`, unit-tested
  with table cases. Refinement is welcome; the mapping table is law.
- **Renderer v1 scope:** Rackling, Tower Brute, Board Sprite, Cabinet
  Colony, Screenrunner, Patchwork bodies (the six demo species), size
  from RAM, LED eyes, basic expression. Scars/stitches/palettes land
  in Phase 2 — but the render function's signature already takes the
  full input tuple so Phase 2 does not change call sites.
- **SSE over websockets:** one-way live vitals are exactly SSE's job,
  and stdlib `http.Flusher` covers it with no dependency.
- **Tokens:** fleet token from env (`LAB_FLEET_TOKEN`); per-familiar
  tokens are random 32-byte hex stored hashed in `familiar`.
- **Config via env + flags only** (`LAB_ADDR`, `LAB_DB`, `LAB_DEMO`,
  `AGENT_LAB_URL`, `AGENT_TOKEN`). No config files in Phase 1.

## Non-goals for Phase 1

Mood, scars, diary, actions, Telegram, MCP — all later phases. Phase 1
vitals are raw numbers on a living-looking card; that is enough.
