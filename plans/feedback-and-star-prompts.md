# Plan: Feedback & GitHub Star Prompts

> Source PRD: `diny/prd/feedback-and-star-prompts.md`

## Architectural decisions

Durable decisions that apply across all phases:

- **Routes**: new `POST /api/feedback` on `diny-next`. Existing `/api/requests` is not modified.
- **Schema**: new `feedback` table (Drizzle migration). Columns: `id` (uuid pk), `user_id` (fk → `users.id`), `type` (`'rating' | 'star'`), `value` (text; `'1'..'5'` for rating, `'starred' | 'already_given' | 'dismissed'` for star), `version`, `system`, `repo_name`, `created_at` (timestamptz default now). Stored as text so star outcomes are self-describing.
- **User identity**: reuse the email-based upsert pattern from `/api/requests` — no new auth, no API key.
- **State file**: `~/.config/diny/state.yaml`, separate from `config.yaml`. Schema: `prompts.commit_count` (int) + `prompts.rating.{status,value,answered_at}` + `prompts.star.{status,answered_at}`. `status` enums: rating = `pending | rated | dismissed`; star = `pending | starred | already_given | dismissed`. Missing file = all-defaults. Never overlaid by project-local configs.
- **CLI packages**:
  - `prompts/` — eligibility gate, prompt UI (simple blocking keypress via `ui.Box`), state I/O. Exposes `MaybeShow(cfg, state)` called once at the end of a successful commit.
  - `feedback/` — HTTP client mirroring `groq/sendRequest`. POSTs to `/api/feedback` with 5s timeout; failures are swallowed silently.
- **Config**: new `prompts.enabled` boolean, default `true`, in `defaults.yaml`, `Config`, and `LocalConfig` (as `*bool` for override semantics).
- **Tunable constants** (compile-time, single file): `MinCommitsBeforePrompt = 5`, `TriggerProbability = 0.15`, `GitHubRepoURL = "https://github.com/dinoDanic/diny"`.
- **Invocation surface**: `diny commit` only. No prompts on `timeline`, `yolo`, `--print`, non-TTY, or CI environments.
- **Telemetry policy**: rating `0` (dismiss) is **not** POSTed — local only. All four star outcomes are POSTed. Dev env (`NODE_ENV === "development"`) skips DB writes on the backend, matching `/api/requests`.

---

## Phase 1: Rating prompt end-to-end

**User stories**: 1, 2, 3, 4, 5, 18, 20, 23, 24, 26, 27, 29, 31, 32, 33, 34

### What to build

A full vertical slice that delivers the rating prompt with no randomness yet — after a successful `diny commit`, once the user has reached the 5-commit gate, the rating prompt appears every time until answered. This slice proves out every layer: state file, config flag, eligibility gate, prompt UI, HTTP client, backend endpoint, DB schema, user upsert.

- CLI: new `prompts/` package with state loader/saver for `~/.config/diny/state.yaml`, eligibility gate (commit count + `prompts.enabled` + pending-status check), rating UI rendered via `ui.Box` with blocking keypress read.
- CLI: new `feedback/` package that POSTs `{type, value, email, name, version, system, repoName}` to `/api/feedback` with a short timeout, silent on failure.
- CLI: new `prompts.enabled` config entry (default `true`) in `defaults.yaml`, `Config`, and `LocalConfig`.
- CLI: hook `prompts.MaybeShow` into the successful-commit path (post-`TryCommit`), bumping `commit_count` and gating on ≥ 5.
- Backend: new `app/api/feedback/route.ts` with Zod validation (`type: "rating"`, `value: 1..5`, required identity fields), email upsert into `users`, insert into new `feedback` table. Rejects malformed payloads with 400. Skips DB write in dev.
- Backend: Drizzle migration creating the `feedback` table per the architectural-decisions shape.

For this phase only, `TriggerProbability` behaves as `1.0` (always show when eligible and pending) so the slice is easy to demo. No star prompt, no random selection, no `--print` / CI skip yet.

### Acceptance criteria

- [ ] `~/.config/diny/state.yaml` is created on first successful commit with `commit_count: 1`.
- [ ] On the 5th successful commit, the rating prompt appears inside a themed `ui.Box` with a visible 1–5 + 0 key hint line.
- [ ] Pressing `1`–`5` persists `rating.status = rated` with the chosen value + timestamp in `state.yaml`, prints a short ack, and POSTs to `/api/feedback`.
- [ ] Pressing `0` persists `rating.status = dismissed` + timestamp and does **not** POST.
- [ ] Once `rating.status != pending`, subsequent successful commits never re-show the rating prompt.
- [ ] `prompts.enabled: false` in `config.yaml` suppresses the prompt entirely (state file still tracks `commit_count`).
- [ ] A Drizzle migration for the `feedback` table is committed and applies cleanly; a successful rating produces one row joinable to the correct `users` row by email.
- [ ] Hitting `/api/feedback` with `{type: "rating", value: 7}` or missing fields returns a 400 and no DB write.
- [ ] A failed HTTP call from the CLI (server unreachable, non-2xx, timeout) does not print an error to the user and does not block the commit flow.

---

## Phase 2: Star prompt end-to-end

**User stories**: 6, 7, 8, 9, 10, 14, 19, 28, 30

### What to build

Add the GitHub-star prompt as a second pending-capable prompt that lives next to the rating prompt. Still no randomness — for this phase, if rating is already answered and star is pending, the star prompt shows on every eligible commit. Cross-platform browser opening and the fallback path land here, as does independent per-prompt state so answering one never locks the other.

- CLI: star UI in `prompts/` with key map `[1] star · [2] don't ask again · [3] already starred · [0] close`.
- CLI: cross-platform "open URL" helper — `xdg-open` on Linux, `open` on macOS, `rundll32 url.dll,FileProtocolHandler` (or `start`) on Windows. If the command fails or is absent, print the URL and still mark `star.status = starred`.
- CLI: extend `feedback` client to send `type: "star"` with `value` as one of the four outcome strings (`starred`, `already_given`, `dismissed`, `dismissed` — `0` and `2` collapse to `dismissed`).
- CLI: prompt-selection logic updated so that when only one of rating/star is pending, that one is shown; when both are pending, still deterministic in this phase (prefer whichever order is simpler — say rating first — to keep Phase 3's addition of randomness minimal).
- Backend: extend Zod validator on `/api/feedback` to accept `type: "star"` with `value` in the outcome enum; insert into `feedback` table using the same columns.

### Acceptance criteria

- [ ] After rating is answered (any value 1–5 or 0), the next eligible commit shows the star prompt.
- [ ] Pressing `1` opens the GitHub repo URL in the default browser on Linux, macOS, and Windows; `star.status = starred` persists; a `star / starred` row is POSTed.
- [ ] Pressing `3` persists `star.status = already_given`, POSTs `star / already_given`, and does not open a browser.
- [ ] Pressing `2` or `0` persists `star.status = dismissed` and POSTs `star / dismissed`. Both keys produce an identical DB row and persisted state.
- [ ] When the browser-open command fails (simulated, e.g. no `xdg-open`), the CLI prints the GitHub URL to stdout and still marks `star.status = starred` + POSTs.
- [ ] Once `star.status != pending`, subsequent eligible commits never re-show the star prompt.
- [ ] Backend rejects `{type: "star", value: "nope"}` with 400; rejects `{type: "star", value: 5}` with 400.
- [ ] Rating and star state are independent: answering one does not mutate the other's fields in `state.yaml`.

---

## Phase 3: Random trigger, exclusivity, non-interactive skips

**User stories**: 11, 12, 13, 15, 16, 21, 22, 25

### What to build

Apply the real eligibility rules so the feature ships its intended behavior. Probability is the headline change, but this phase also enforces success-only triggering, one-prompt-per-session exclusivity, and non-interactive skips (`--print`, non-TTY stdout, common CI env vars).

- CLI: replace the "always show when eligible" stub with a 15% (`TriggerProbability`) random roll on each eligible commit. Seed from `time.Now().UnixNano()` per invocation — no persisted RNG state.
- CLI: when both rating and star are pending and the roll succeeds, pick 50/50.
- CLI: guarantee at most one prompt per `diny commit` invocation.
- CLI: skip prompts (and skip incrementing `commit_count`? — **no, keep incrementing**; only skip the prompt itself) when:
  - `--print` flag is set.
  - stdout is not a TTY (use `golang.org/x/term` or `mattn/go-isatty`).
  - Any of `CI`, `GITHUB_ACTIONS`, `NONINTERACTIVE` env vars is truthy.
- CLI: verify that save-draft, cancel, regenerate, and commit-failure paths never reach `prompts.MaybeShow`. Only the `TryCommit`-success path does.

### Acceptance criteria

- [ ] With rating and star both pending and `commit_count >= 5`, running 100 scripted successful commits results in a prompt appearing in roughly 10–20 of them (empirical, not asserted in a test — a counter log is enough).
- [ ] Across many eligible commits where both prompts are pending, each prompt is chosen roughly half the time when a roll succeeds.
- [ ] A single `diny commit` invocation never shows both prompts back-to-back.
- [ ] `diny commit --print` never shows a prompt, even if eligible and the roll would succeed.
- [ ] Running under `CI=true` never shows a prompt.
- [ ] Piping stdout to a file (`diny commit > out.txt`) never shows a prompt.
- [ ] Save-draft, cancel, regenerate, and failed commits do **not** increment `commit_count` and do **not** trigger `MaybeShow`.
- [ ] A successful commit **does** increment `commit_count` even when the random roll fails or prompts are suppressed by TTY/CI/`--print`, so eligibility continues to accrue.

---

## Phase 4: Resilience & polish

**User stories**: 17, 28, 32, 33, 34

### What to build

Harden the edges: recover from a corrupt state file, verify the global opt-out behaves correctly through the config override chain, confirm the client never blocks the user on network issues, and validate backend rejection paths.

- CLI: `state.yaml` load/recover mirroring `config.LoadOrRecover` — on invalid YAML or schema violation, rename to `state.backupN.yaml` and recreate defaults. Non-fatal; log once.
- CLI: confirm `feedback` client uses a short timeout (~5s), swallows all errors, and never writes to stderr on failure.
- CLI: verify `prompts.enabled` respects the full override chain — global → versioned `.diny.yaml` → local `<gitdir>/diny/config.yaml` — using the existing `mergeConfigWithLocal` machinery.
- Backend: tighten Zod schema; confirm 400 responses for all malformed payload shapes, including missing identity fields and out-of-range values. Ensure the endpoint is idempotent-friendly — repeated identical payloads are harmless inserts.
- Documentation: a short section in the CLI README describing the prompt behavior and the `prompts.enabled` flag.

### Acceptance criteria

- [ ] Corrupting `state.yaml` by hand (invalid YAML or an unknown `status`) causes diny to back it up as `state.backup1.yaml`, recreate a default state file, continue the commit flow, and not crash.
- [ ] Killing network access before a rating submission still lets the user complete the commit; no error surfaces to the user; state is persisted as if the POST succeeded.
- [ ] Setting `prompts.enabled: false` in a versioned or local project config correctly suppresses prompts for that repo while leaving global behavior unchanged elsewhere.
- [ ] `/api/feedback` returns 400 for: missing `email`, missing `type`, `rating` with `value: 0`, `star` with `value: 6`, unknown `type`.
- [ ] Sending the same valid rating payload twice produces two `feedback` rows (harmless duplicate); no unique-constraint violation.
- [ ] CLI README documents the prompt behavior, the 5-commit gate, and the `prompts.enabled` opt-out.
