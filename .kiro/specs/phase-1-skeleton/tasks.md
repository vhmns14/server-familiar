# Phase 1 — Skeleton: Tasks

Work in order. After each numbered group: `gofmt -l .`,
`go vet ./...`, `go build ./...`, `go test ./...` must pass.

## 1. Scaffold

- [x] 1.1 Initialise Go module `server-familiar`; create `cmd/agent`,
      `cmd/lab`, and the `internal/` packages from design.md (empty
      but compiling).
- [x] 1.2 Write `.gitignore` (binaries, `*.db*`, `.env`, editor/OS junk).
- [x] 1.3 Write the README skeleton (pitch, placeholder quick start,
      link to PROJECT_PLAN.md).
- [x] 1.4 Verify: clean clone → `go build ./...` succeeds.

Added while scaffolding, beyond the task text: `.env.example` documenting
every environment variable with no secrets in it, and a `Makefile` whose
`check` target is the mandated `gofmt` / `vet` / `build` / `test` gate.

## 2. Store

- [x] 2.1 Add migration runner + schema v1: `familiar`, `life`,
      `metric_sample` (fields per PROJECT_PLAN.md §8, Phase 1 subset),
      SQLite WAL mode.
- [x] 2.2 Implement store functions: create familiar + first life,
      insert sample, list familiars with latest sample, touch last-seen.
- [x] 2.3 Table-driven tests for the store against a temp DB file.

Decisions taken here, recorded because later phases inherit them:

- The dependency conversation `design.md` anticipated is settled:
  **`modernc.org/sqlite` v1.62.0**, approved by the keeper. Pure Go, so
  `CGO_ENABLED=0` still yields a static binary. Only `cmd/lab` imports it;
  `go list -deps ./cmd/agent` shows no third-party package, and the agent
  binary is unchanged at 1.7 MB against the lab's 10.3 MB.
- **A host profile belongs to `life`, not to `familiar`.** A familiar is the
  soul, a life is a body, and a body is a machine. Reincarnation in Phase 4
  therefore closes one life and opens another without touching the familiar
  row. A partial unique index enforces at most one open life per familiar.
- **Timestamps are INTEGER Unix milliseconds**, not text: text timestamps
  mis-sort as soon as fractional-second width varies, and "the latest
  sample" has to be exact.
- Tables are `STRICT`; tokens are stored only as a SHA-256 hash; foreign
  keys cascade, so deleting a familiar takes its lives and samples with it.
- Migrations run when `store.Open` is called, which `cmd/lab` does at
  startup — the hard rule in `tech-rules.md`.

## 3. Agent

- [ ] 3.1 Implement `/proc` samplers: CPU %, memory, disk (statfs),
      uptime. Unit-test parsers against fixture `/proc` files.
- [ ] 3.2 Implement the 10 s loop + ingest client with backoff and a
      bounded retry buffer.
- [ ] 3.3 Verify: static `CGO_ENABLED=0` build runs and posts from
      this machine; measure RSS on the smallest available host or
      container (target: tens of MB).

## 4. Ingest and enrolment

- [ ] 4.1 `POST /api/enrol` (fleet token) → creates familiar + life,
      returns per-familiar token. Species derived from host profile.
- [ ] 4.2 `POST /api/ingest` (per-familiar token) → stores sample,
      updates last-seen, publishes to SSE hub. 401 on bad token.
- [ ] 4.3 Write `install/install.sh`: download/copy binary, enrol,
      write env file, print the familiar's name and species.
- [ ] 4.4 Tests: enrol → ingest → row visible; bad token rejected.

## 5. Renderer v1

- [ ] 5.1 Implement `render.Creature(input) string (SVG)` with the full
      input tuple from creature-design.md; v1 bodies for the six demo
      species.
- [ ] 5.2 Size-from-RAM scaling; LED eyes; deterministic output.
- [ ] 5.3 Golden tests: same input → byte-identical SVG; each species
      renders valid, non-empty SVG.

## 6. Web UI + SSE

- [ ] 6.1 Lab home page: server-rendered creature cards (inline SVG +
      vitals), hand-written CSS, laboratory theme (dark storm lab —
      see mockup direction in PROJECT_PLAN.md).
- [ ] 6.2 SSE hub + `GET /events`; vanilla JS patches vitals in place;
      unreachable state after 60 s without a sample.
- [ ] 6.3 Familiar detail stub page: big render + raw vitals table
      (character sheet arrives in Phase 2).

## 7. Demo mode

- [ ] 7.1 `internal/demo`: seed the six showcase familiars with the
      exact facts from product.md (Anls pinned at 1 Mbps egress, etc.).
- [ ] 7.2 Simulator: per-familiar telemetry loops with plausible
      diurnal wobble, written through the real store path.
- [ ] 7.3 `LAB_DEMO=1` wiring + a visible "DEMO MODE" banner.
- [ ] 7.4 Verify: fresh DB + demo mode → six live cards updating,
      no agent running.

## Phase exit check

- [ ] E.1 `go build ./...` and `go test ./...` green from a clean clone.
- [ ] E.2 A real host enrolled via install.sh and reporting.
- [ ] E.3 Summarise for the user: what works, what was tested, what is
      stubbed. **Wait for approval before Phase 2.**
