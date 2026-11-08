# Tech Rules — Server Familiar

*Steering file: always loaded. These are hard rules. If a task seems
to require breaking one, stop and ask the user instead.*

## Language and layout

- **Go** for both binaries, one module: `cmd/agent` and `cmd/lab`,
  shared code under `internal/`. Keep packages small and single-purpose.
- **Standard library first** — `net/http` for the lab, no web
  framework. Propose any third-party dependency in chat and wait for
  approval before adding it.
- **CGO off** (`CGO_ENABLED=0`). Both binaries must be single, static,
  copy-it-and-run files.
- The **agent must run on a 512 MB Alpine (musl) host** with half a
  CPU core. Budget its memory in tens of MB. No runtime dependencies,
  no shelling out to tools that may not exist — read `/proc` directly.

## Storage

- **SQLite, single file, WAL mode.** Migrations run at lab startup.
  Backup must be "copy one file". No external database, no Redis, no
  queue — ever, in this build.

## Web UI

- **Server-rendered HTML + vanilla JS + SSE** (Server-Sent Events) for
  live vitals.
- **NO frontend framework. NO build step. NO npm / node_modules.**
  If a page needs JavaScript, it is a hand-written `<script>` or a
  static file served by the lab.
- CSS is hand-written in static files; no preprocessor, no CDN
  dependency required for the app to work offline.

## Creatures

- **Procedural SVG only**, rendered server-side. No image assets,
  sprite sheets, or icon fonts for creatures. Rendering must be
  deterministic from (species, stage, palette seed, scar list, mood,
  size) — see `creature-design.md`.

## LLM usage

- **LLM access only from the lab** — never from the agent.
- Via a **configurable OpenAI-compatible endpoint**: base URL, model,
  and API key all from environment variables.
- **Response caching** keyed by the grounded input state, a
  **configurable daily token budget**, and a **template fallback**:
  when no LLM is configured (or the budget is spent), diary, dreams,
  and chat fall back to grounded hand-written templates. Generated
  features must never hard-fail for lack of an API key.

## Security and secrets

- **Secrets via environment variables only.** Never commit a token,
  key, or password; never put one in code, fixtures, tests, or docs.
- Agents are **outbound-only**: they POST metrics and poll for work;
  the lab never connects into a server.
- Agent actions are a **named allow-list with fixed behaviour** —
  never accept or execute arbitrary command strings.
- Fleet token required on ingest; per-familiar token after enrolment.
- The public terrarium is read-only behind an unguessable,
  revocable token.

## Process

- After each task group: `gofmt`, `go vet ./...`,
  `go build ./...`, `go test ./...` — all must pass before a phase is
  called done.
- Tests use the standard `go test` tooling; prefer table-driven tests
  for the mood engine, incident detection, and renderer determinism.
- **English** for all code, comments, commit messages, and docs.
- Keep the `/.kiro` folder intact at the repo root — it is part of the
  hackathon submission.
- Commit in small, working increments; never commit a broken build.
