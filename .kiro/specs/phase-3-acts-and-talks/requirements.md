# Phase 3 — Acts & Talks: Requirements

**Goal:** familiars start *doing* and *speaking* — permissioned
self-care actions, a Telegram nervous system, guard senses that turn
security telemetry into instinct, grounded chat, and an MCP interface.

Depends on: Phase 2 complete and approved.

## User stories

1. As the operator, I want my familiars to fix small problems
   themselves — but only with my permission where I've asked for it —
   so that the fleet heals without me babysitting it, and never acts
   behind my back.
2. As the operator, I want approvals and security alerts on Telegram,
   so that the fleet can reach me on my phone in seconds.
3. As the operator, I want break-in attempts presented as a creature
   reacting (growl, scratches, bark) with the evidence attached, so
   that security is glanceable and still auditable.
4. As the operator, I want to chat with a familiar and ask my AI
   assistant "how is my fleet?", so that the fleet is conversational,
   not just visible.

## Acceptance criteria (EARS style)

### Actions and permissions

- THE SYSTEM SHALL offer exactly the allow-listed actions: `restart_service`
  (a service named in the agent's local config), `clear_cache`
  (package/tmp caches), and `health_ritual` (log rotation + connectivity
  and disk checks). THE SYSTEM SHALL NOT accept any other action name
  or any free-form command.
- WHEN a familiar's permission mode for an action is **Auto** and its
  trigger condition holds (per action rules), THE SYSTEM SHALL dispatch
  the action and write an `action_log` row with mode=`auto`.
- WHEN the mode is **Ask** and a trigger holds, THE SYSTEM SHALL create
  a pending approval shown in the lab and sent to Telegram, and SHALL
  dispatch the action ONLY IF the operator approves before it expires
  (default 15 min).
- WHEN the agent polls and finds an approved action, THE SYSTEM SHALL
  execute it, and SHALL report the result back so the lab closes the
  `action_log` row with success/failure and detail.
- IF an action fails on the agent, THEN THE SYSTEM SHALL record the
  failure, notify the operator, and NOT retry automatically.
- Agents SHALL receive actions only via their outbound poll — THE
  SYSTEM SHALL NOT open any inbound connection to an agent host.

### Telegram

- WHEN an approval is pending, a guard event fires, or a familiar dies
  or is resurrected, THE SYSTEM SHALL send a Telegram message to the
  configured operator chat within seconds.
- WHEN the operator taps Approve/Deny on a Telegram approval, THE
  SYSTEM SHALL resolve the pending approval identically to a lab UI
  decision.
- IF Telegram is not configured, THEN approvals SHALL remain fully
  usable in the lab UI and the system SHALL NOT error.

### Guard senses

- WHEN an SSH login succeeds from an IP the familiar has not seen
  before, THE SYSTEM SHALL emit a **growl** event: Night Log entry,
  Telegram alert with IP and geo (where a local lookup is available),
  and the familiar's startled expression briefly.
- WHEN failed SSH logins reach the burst threshold (default ≥ 10 in
  10 min), THE SYSTEM SHALL add a **door-scratch** mark (visible
  counter + render mark on the detail page), record an incident, and
  alert via Telegram.
- WHEN a port scan pattern is detected (many distinct ports from one
  source in a short window), THE SYSTEM SHALL emit a **bark** event
  with the same logging/alerting path.
- THE SYSTEM SHALL obtain auth/scan signals from the agent's local log
  tailing (`auth.log` / journal) where available, and in demo mode from
  the simulator.

### Chat

- WHEN the operator sends a chat message to a familiar, THE SYSTEM
  SHALL reply in character, grounded in that familiar's current derived
  state, recent metrics, and events, using the same LLM/template
  machinery and budget as the diary.
- IF the operator asks about a metric the familiar does not have, THEN
  the familiar SHALL say so in character rather than invent a value.

### MCP server

- THE SYSTEM SHALL expose an MCP server with at least the tools
  `list_familiars` (name, species, status, mood), `how_is` (one
  familiar's grounded status summary), and `fleet_mood` (aggregate
  fleet state).
- WHEN an MCP tool is called, THE SYSTEM SHALL answer from the same
  stored state the UI shows — never from a separate cache that can drift.
