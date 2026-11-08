# Phase 5 — Submission: Tasks

Work in order. Feature freeze: fixes and documents only.

## 1. Polish pass

- [ ] 1.1 Full route walk in demo mode and real mode; fix list triaged
      (broken → inconsistent → wording); fixes committed small.
- [ ] 1.2 Fallback-state check: no LLM key, no Telegram — every feature
      degrades as designed, no error pages.
- [ ] 1.3 Screenshot set captured from demo mode (lab wide shot,
      character sheet, ceremony, graveyard) into `docs/img/`.

## 2. README

- [ ] 2.1 Rewrite README.md per design.md structure (pitch, demo quick
      start, enrolment, feature tour, architecture, tech stack).
- [ ] 2.2 "Built with Kiro" section: `/.kiro` steering + specs
      explained, credit allocation 40/30/20/10, honest build notes.
- [ ] 2.3 Env-var configuration table (names and purpose only).
- [ ] 2.4 Proofread: every command in the README executed verbatim from
      a clean clone by the letter.

## 3. Demo-mode hardening

- [ ] 3.1 `scripts/demo-smoke.sh` (or Go test): clean build, demo boot,
      all key routes 200, SSE emits.
- [ ] 3.2 Dry-run the full 3-minute demo script end-to-end, timed;
      fix any beat that needs manual setup.
- [ ] 3.3 Confirm demo mode needs no network access.

## 4. Submission assets

- [ ] 4.1 `docs/devpost-submission.md` drafted in full (per design.md),
      ready to paste with only links to fill.
- [ ] 4.2 `docs/demo-video-shotlist.md`: per-beat screen, actions,
      narration, fallback takes.
- [ ] 4.3 Record the demo video from the shot list; upload; link it in
      the README and the Devpost draft.

## 5. Final verification checklist

- [ ] 5.1 Repository is **public**; `/.kiro` present at the root with
      steering + all five specs, tasks ticked.
- [ ] 5.2 Clean clone (fresh directory, no caches): `go build ./...`
      and `go test ./...` pass; demo smoke passes.
- [ ] 5.3 No secrets anywhere in the repo or history (grep for tokens,
      keys, `.env`); `.gitignore` covers the database and binaries.
- [ ] 5.4 README renders correctly on GitHub; all links live.
- [ ] 5.5 Devpost submission filled (text, repo link, video link,
      category: Frankenstein) and **submitted** — target: one full day
      before 5 Dec 2026, 05:00 WIB; the final day is buffer only.
- [ ] 5.6 Tag the submission commit (e.g. `kiroween-submission`) and
      push.

## Phase exit check

- [ ] E.1 The user confirms the Devpost dashboard shows the submission
      as received. Nothing else counts as done.
