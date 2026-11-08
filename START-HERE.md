# Server Familiar — Start Here

This kit is everything you need to hand the project to Kiro and start building.

**Event:** Kiroween 2026 — Frankenstein category
**Submission deadline:** 5 December 2026, 05:00 WIB
**Repo:** public GitHub repository named `server-familiar`, built with Kiro

---

## Setup (5 minutes)

1. **Create the repository**
   Create a new **public** GitHub repository named `server-familiar`.
   Do not initialise it with a README, `.gitignore`, or licence — this kit
   and Kiro will create those, and an empty repo avoids merge conflicts.

2. **Unzip this kit into the repo root**
   Extract the ZIP so the repo root looks exactly like this:

   ```text
   server-familiar/
   ├── START-HERE.md
   ├── KICKOFF-PROMPT.md
   ├── PROJECT_PLAN.md
   └── .kiro/
       ├── steering/
       └── specs/
   ```

   The `/.kiro` folder **must sit at the repository root**. That is a
   Kiroween rule — the judges check for it there.

3. **Open the repo in Kiro**
   Clone the repo locally, commit and push the kit files as the first
   commit, then open the folder in Kiro.

4. **Paste the kickoff prompt**
   Open `KICKOFF-PROMPT.md`, copy its entire contents, and paste it as
   your **first message** to Kiro in that workspace. It tells Kiro what
   to read, how to scaffold, and what to build first.

5. **Work phase by phase**
   The specs in `/.kiro/specs/` are ordered:

   1. `phase-1-skeleton` — agent, lab, storage, creature cards, demo mode
   2. `phase-2-alive` — mood, growth, incidents, scars, diary, Night Log
   3. `phase-3-acts-and-talks` — actions, Telegram, guard senses, chat, MCP
   4. `phase-4-living-features` — reincarnation, dreams, graveyard, bonds,
      runs, growth records, public terrarium, Life Book
   5. `phase-5-submission` — polish, README, Devpost text, demo video, checks

   **Kiro should complete Phase 1 fully — building and passing tests —
   before starting Phase 2**, and so on. Each phase spec ends with tasks;
   ask Kiro to confirm a phase is done before you approve the next one.

## Where the reasoning lives

`PROJECT_PLAN.md` is the full plan: why this project exists, how it
differs from ordinary monitoring, the creature system, the architecture,
the tech stack, the data model, the build schedule, the Kiro credit
budget, and the 3-minute demo script for the judges. Read it once
yourself; Kiro will read it too.

## Files in this kit

| File | For | Purpose |
| --- | --- | --- |
| `START-HERE.md` | You | This setup guide |
| `KICKOFF-PROMPT.md` | Kiro | The first message to paste into Kiro |
| `PROJECT_PLAN.md` | You + Kiro | The complete project plan |
| `.kiro/steering/product.md` | Kiro | Product brief, always loaded |
| `.kiro/steering/tech-rules.md` | Kiro | Hard technical rules, always loaded |
| `.kiro/steering/creature-design.md` | Kiro | Creature design rules, always loaded |
| `.kiro/specs/phase-*/` | Kiro | Requirements, design, and tasks per phase |

## Ground rules to hold Kiro to

- Go, standard library first. No frontend framework, no build step, no npm.
- The `/.kiro` folder stays intact at the repo root.
- `go build ./...` and `go test ./...` pass before any phase is called done.
- Secrets come from environment variables only — never commit one.
- Everything in English: code, comments, commits, and docs.
