# Phase 3 — Acts & Talks: Design

## New components

```text
internal/actions/    action definitions, trigger rules, approval queue,
                     dispatch + result handling, action_log writes
internal/agent/      + action poll loop, local executor (allow-list),
                     auth-log tailer for guard signals
internal/telegram/   bot client (Bot API over net/http): send alerts,
                     inline Approve/Deny buttons, callback handling
internal/guard/      guard-sense detectors: known-IP tracking,
                     failed-login bursts, port-scan pattern
internal/chat/       familiar chat endpoint (reuses diary context builder)
internal/mcp/        MCP server: list_familiars, how_is, fleet_mood
store migrations v3: action_log table; approvals on familiar config;
                     known_ips table; chat handled statelessly per request
```

## Action flow

```text
trigger (lab rules engine, e.g. service down in samples / disk ≥ 95%)
  → permission check (familiar × action mode)
      Auto → action queued → agent polls (piggyback on ingest response
             or a lightweight GET /api/actions?token=…) → executes →
             POSTs result → action_log closed, Night Log entry
      Ask  → approval row + lab banner + Telegram message with buttons
             → approve: same dispatch path · deny/expire: action_log
             closed as denied/expired, familiar told (diary context)
```

- The executor on the agent maps action names to fixed local functions.
  `restart_service` may only restart services listed in the agent's own
  local config file (created at install); the lab can never introduce a
  new service name at dispatch time.
- Every state transition writes to `action_log` — it doubles as the
  audit trail and as diary/Night Log source material.

## Guard senses

- Agent-side signal collection: tail the auth log (path configurable;
  journal fallback) and count connection attempts per source per port
  from the same tail where possible. Signals are sent as a compact
  `guard` block on the ingest payload — raw log lines never leave the
  host.
- Lab-side `internal/guard` keeps per-familiar known IPs, burst
  windows, and scan windows; emits incidents with kind
  `unknown_login` / `failed_login_burst` / `port_scan`, mapped to the
  growl / scratches / bark presentation. Door scratches are a counter
  + incident list in Phase 3; making them a Phase 2 scar kind is
  allowed if trivial (scar kind `scratch cluster` already reserved).
- Demo mode: the simulator exposes buttons/endpoints to fire each guard
  event on demand — the judge demo depends on this.

## Telegram

- Plain Bot API via `net/http` (sendMessage + inline keyboard +
  answerCallbackQuery + getUpdates long-poll in the lab). Token and
  chat id from env (`TELEGRAM_BOT_TOKEN`, `TELEGRAM_CHAT_ID`); the
  whole package is a no-op stub when unset.
- Messages are short, in-world, and evidence-carrying: e.g. "🐾 Anls
  GROWLED — new SSH login from 203.0.113.7 (geo: —). [Open in lab]".

## Chat

`POST /api/familiars/{id}/chat` → builds the diary grounded context +
the operator's message → LLM or template responder → returns the reply
and stores nothing but an optional diary note. Templates cover the
common questions (how are you, are you hungry, what happened today)
with slot-filled grounded answers; free-form questions without an LLM
get an honest in-character "my keeper hasn't given me a voice for that
yet" fallback — never an invented answer.

## MCP

- A small MCP server (stdio or HTTP transport — pick per current MCP
  Go options; propose the library if one is needed, per tech-rules)
  in `cmd/lab` as a subcommand (`lab mcp`) so deployment stays one binary.
- Tools read through the same store queries as the web UI.
- `fleet_mood` aggregates: counts by status/mood, the unhappiest and
  happiest familiar, any active guard events — a paragraph-ready summary.
