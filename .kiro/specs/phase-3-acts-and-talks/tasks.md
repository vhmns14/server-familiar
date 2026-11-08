# Phase 3 — Acts & Talks: Tasks

Work in order. Build + tests green after each group.

## 1. Schema + action engine

- [ ] 1.1 Migration v3: `action_log`, per-familiar/action permission
      config, `known_ips`.
- [ ] 1.2 `internal/actions`: the three allow-listed actions, trigger
      rules, approval queue with expiry, dispatch + result handling.
- [ ] 1.3 Tests: Auto dispatches without approval; Ask never dispatches
      unapproved; expired approval closes as expired; unknown action
      name rejected.

## 2. Agent executor + guard signals

- [ ] 2.1 Agent poll for approved actions (on the ingest cycle);
      executor with fixed functions; `restart_service` limited to the
      agent's local service list.
- [ ] 2.2 Result reporting → lab closes the action_log row.
- [ ] 2.3 Auth-log tailer producing the compact `guard` ingest block
      (successful logins + source IPs, failed counts, distinct ports).
- [ ] 2.4 Verify on a real host: trigger a cache clear in Auto; trigger
      a restart in Ask and approve it from the lab.

## 3. Telegram

- [ ] 3.1 `internal/telegram` no-op-safe client; alerts for approvals,
      guard events, death, resurrection.
- [ ] 3.2 Inline Approve/Deny with callback handling resolving the same
      approval rows as the lab UI.
- [ ] 3.3 Verify: full Ask flow completed from a phone.

## 4. Guard senses

- [ ] 4.1 `internal/guard`: known-IP growl, failed-login burst →
      door scratches (counter + incident + detail-page marks),
      port-scan bark.
- [ ] 4.2 Night Log + Telegram wiring; startled expression pulse.
- [ ] 4.3 Demo-mode triggers for all three events (buttons in demo UI).
- [ ] 4.4 Tests with synthetic guard blocks for each detector + a
      below-threshold case that fires nothing.

## 5. Chat

- [ ] 5.1 Chat endpoint reusing the grounded context builder; LLM path
      with cache/budget; template responder for common questions;
      honest fallback for unanswerable questions.
- [ ] 5.2 Chat panel on the familiar detail page (vanilla JS).
- [ ] 5.3 Test: ask about an absent metric → in-character refusal, no
      invented value (assert against template path; spot-check LLM path).

## 6. MCP server

- [ ] 6.1 `lab mcp` subcommand; tools `list_familiars`, `how_is`,
      `fleet_mood` reading live store state.
- [ ] 6.2 Verify with an MCP client: "how is my fleet feeling?" returns
      the aggregate summary matching the lab UI.

## Phase exit check

- [ ] E.1 Build + tests green.
- [ ] E.2 Scripted run-through of the Phase 3 demo beats (approval from
      Telegram, guard burst in demo mode, MCP query) without errors.
- [ ] E.3 Summarise for the user and **wait for approval before Phase 4.**
