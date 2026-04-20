# Feedback & GitHub Star Prompts

## Problem Statement

Diny has no channel for hearing from the people who actually use it. There is no rating signal, no feature-request capture, and no ask for the one thing that materially helps a free OSS CLI grow — a GitHub star. The tool is used after every `git commit`, which is the highest-attention moment the user has with it, but that attention is currently thrown away. At the same time, the user doesn't want to build anything that nags: a prompt every commit would train people to mash `q` on reflex and poison any signal the tool did manage to collect.

## Solution

After a successful commit, occasionally (~15% chance), surface one of two lightweight, single-keypress prompts:

1. **Rating prompt** — "How's diny? 1–5, or 0 to never ask again." One keystroke, done.
2. **GitHub star prompt** — "Diny is free; a star helps. 1=star, 2=don't ask again, 3=already starred, 0=close." One keystroke, done. Picking `1` opens the repo in the browser.

Only one of the two prompts can appear in a session. Each prompt has its own "done" state: once a user has rated (any value 1–5) or dismissed (0), the rating prompt never shows again; the star prompt has the same finality on any answer. A new gate prevents either prompt from appearing until the user has made at least 5 successful commits, so evaluators aren't ambushed. Power users can disable the whole thing with a single config flag.

Ratings are sent to a new `/api/feedback` endpoint on `diny-next`. The existing `/api/requests` endpoint stays focused on LLM calls.

## User Stories

1. As a diny user, I want to be asked how I feel about the tool, so that the author knows whether it's working for me.
2. As a diny user, I want to give a rating with a single keypress, so that responding costs me less than a second.
3. As a diny user, I want to press `0` to permanently dismiss the rating prompt, so that if I don't want to rate I only have to decline once.
4. As a diny user, I want the rating scale to be 1–5, so that it matches the mental model I already have for app ratings.
5. As a diny user, I don't want a text input after rating, so that I can answer and move on instantly.
6. As a diny user, I want to be asked occasionally to star the project on GitHub, so that I have a nudge to support a free tool I use.
7. As a diny user who wants to star, I want pressing `1` to open the GitHub repo in my browser, so that I don't have to copy/paste a URL.
8. As a diny user who already starred the repo, I want to press `3` to tell diny that, so that I'm not nagged and the author gets a correct signal.
9. As a diny user who doesn't want a star prompt, I want to press `2` to permanently dismiss it, so that it never returns.
10. As a diny user who wants to defer the decision, I want to press `0` to close the prompt, so that I can skip it without committing to never seeing it again — though under the current policy this is equivalent to "don't ask again".
11. As a new diny user, I don't want to see any prompt until I've committed a handful of times, so that my first impression is about commit quality, not marketing.
12. As a diny user who commits frequently, I don't want prompts on every commit, so that diny doesn't feel pushy.
13. As a diny user, I don't want to see two prompts back-to-back after the same commit, so that one `git commit` never triggers multiple interruptions.
14. As a diny user who has answered both prompts, I want diny to go silent forever on this topic, so that I'm rewarded for engaging.
15. As a diny user running diny in CI, automation, or scripts, I don't want prompts to ever appear, so that non-interactive flows stay clean.
16. As a diny user who pipes output (`diny commit --print`), I don't want prompts to appear, so that scripted output isn't corrupted.
17. As a privacy-minded diny user, I want a single config flag to turn prompts off entirely, so that I don't have to wait for them to appear to dismiss them.
18. As the diny author, I want rating data sent to a backend so that I can see aggregate sentiment, not just per-user vibes.
19. As the diny author, I want star-prompt outcomes (`starred`, `already_given`, `dismissed`) sent to the backend, so that I know how many of the ~15% rolls convert to actual stars.
20. As the diny author, I want feedback requests identified by the same user email that `/api/requests` already upserts, so that I can join rating data to usage data in the DB.
21. As a diny user whose commit failed (hook error, empty commit, etc.), I don't want a rating prompt, so that negative moments don't bias the rating signal downward.
22. As a diny user who only saved a draft or cancelled, I don't want a rating prompt, so that prompts only follow the explicit "I committed" moment.
23. As a diny user, I want the prompt to render with the same theme/box styling as the rest of the TUI, so that it feels like part of diny and not a bolted-on popup.
24. As a diny user, I want the prompt to show below (not replace) the commit confirmation, so that I still see which hash/message was committed.
25. As a diny user, I want the random-trigger probability to be tunable by the author, so that annoyance levels can be adjusted based on real usage.
26. As a diny user, I want prompt state (commit count, whether I've rated, whether I've handled the star prompt) persisted across runs, so that diny remembers my decisions and doesn't re-ask.
27. As a diny user who blew away their config, I want prompt state kept separate from config, so that editing `config.yaml` doesn't reset my "don't ask again" choice.
28. As a diny user whose browser can't be opened (SSH, headless), I want a fallback path when `open`/`xdg-open` fails, so that I still see the URL and the prompt doesn't hang.
29. As a diny user, I want star and rating prompts each to track their own state independently, so that answering one doesn't disable the other.
30. As a diny user running on Linux, macOS, and Windows, I want the "open browser" action to work on all three, so that the prompt behaves the same everywhere.
31. As a diny user, I want the prompt help line to show the key mapping (e.g., `[1-5] rate · [0] never`), so that I don't have to remember the scheme.
32. As the diny author, I want the feedback endpoint to reject malformed payloads (wrong type, out-of-range value) with a 400, so that bad clients don't pollute the DB.
33. As the diny author, I want the feedback endpoint to be idempotent-friendly — a second identical send is a harmless insert — so that client-side retries don't corrupt state.
34. As a diny user, I want the prompt interaction to be non-blocking for background work (update check, etc.), so that answering a prompt doesn't delay anything else.

## Implementation Decisions

### Trigger & eligibility

- Prompts are considered **only after `TryCommit` returns successfully**. Save-draft, cancel, regenerate, and error paths do not increment the eligibility counter or trigger prompts.
- `--print` mode and non-TTY stdout (detected via `isatty` on stdout) skip prompts entirely.
- A common `CI`/`GITHUB_ACTIONS`/`NONINTERACTIVE` env-var check also skips prompts.
- The global `prompts.enabled` config flag (default `true`) short-circuits everything.
- After a successful commit, increment a persistent `commit_count`. If `commit_count >= 5` AND at least one of `rating` or `star` is still `pending`, roll a 15% die.
- If the roll succeeds, pick which prompt to show:
  - If only one of the two is still `pending`, show that one.
  - If both are `pending`, pick randomly 50/50.
- Only one prompt per `diny commit` invocation.

### Rating prompt UX

- Rendered inside a `ui.Box` (existing component), shown after the post-commit confirmation line.
- Body text (indicative): `How's diny working for you?`
- Scale line: `[1] 😐  [2] 🙂  [3] 😀  [4] 🤩  [5] ❤️   ·   [0] Don't ask again` (emojis optional; subject to style review — can be plain `1–5`).
- Blocking keypress reader (no Bubble Tea state machine required — a single `bufio` read of one rune).
- Any other key is ignored (re-prompt or fall through after a short timeout — decision: ignore and re-read).
- On `1–5`: mark `rating.status = rated`, `rating.value = N`, POST to `/api/feedback`, print a short `Thanks!` line.
- On `0`: mark `rating.status = dismissed`, no POST (see Backend), print nothing or a minimal ack.

### Star prompt UX

- Same container/styling as rating prompt.
- Body text (indicative): `Diny is free & open source. A GitHub star helps a lot!`
- Scale line: `[1] Star it  [2] Don't ask again  [3] Already starred  [0] Close`
- On `1`: open `https://github.com/dinoDanic/diny` via:
  - `xdg-open` on Linux
  - `open` on macOS
  - `rundll32 url.dll,FileProtocolHandler` (or `start`) on Windows
  - Fallback: print the URL if the command fails or isn't found.
  - Mark `star.status = starred`.
- On `2`: mark `star.status = dismissed`.
- On `3`: mark `star.status = already_given`.
- On `0`: mark `star.status = dismissed` (under the chosen re-ask policy, `0` and `2` have the same persisted effect; the two keys coexist because users may reach for either word).
- All four outcomes POST to `/api/feedback` (see Backend).

### State file

New file: `~/.config/diny/state.yaml` (path derivation mirrors `config.GetConfigPath`). Created on first write. Schema:

```yaml
prompts:
  commit_count: 0
  rating:
    status: pending      # pending | rated | dismissed
    value:               # 1..5 when rated, else unset
    answered_at:         # RFC3339 timestamp when status left pending
  star:
    status: pending      # pending | starred | already_given | dismissed
    answered_at:         # RFC3339 timestamp when status left pending
```

- Read at startup (once), mutations written atomically (temp file + rename).
- Missing file treated as all-default state; missing fields default to `pending` / `0`.
- Schema/validation errors are non-fatal: back up (`state.backupN.yaml`) and recreate, mirroring the `LoadOrRecover` pattern in `config.go`.
- The file is **not** overlaid by project-local configs; state is per-user-per-machine, never team-shared.

### Config changes

Add to `config.Config`:

```go
type PromptsConfig struct {
    Enabled bool `yaml:"enabled" json:"Enabled"`
}
// on Config: Prompts PromptsConfig `yaml:"prompts" json:"Prompts"`
```

- Default in `defaults.yaml`: `prompts.enabled: true`.
- Add mirroring `*bool` field to `LocalConfig` / `LocalCommitConfig`-sibling for overrides.
- Validation: boolean; no other constraints.

### Constants (tunables)

Kept in a single `prompts/constants.go` or similar:

- `MinCommitsBeforePrompt = 5`
- `TriggerProbability = 0.15`
- `GitHubRepoURL = "https://github.com/dinoDanic/diny"`

These are compile-time constants in v1; no runtime override. Changing probability ships in a release.

### Backend: new `/api/feedback` endpoint

New route under `diny-next/app/api/feedback/` with its own `route.ts`:

- **Method**: `POST`.
- **Payload (Zod-validated)**:

  ```json
  {
    "type": "rating" | "star",
    "value": integer,
    "email": string,
    "name": string,
    "version": string,
    "system": string,
    "repoName": string
  }
  ```

  - `type = "rating"` requires `value` in `1..5`.
  - `type = "star"` requires `value` in a small enum representing outcome (see DB schema).
- **Behavior**:
  - Validate payload; reject malformed with `400`.
  - Upsert user by email (reuse logic from `/api/requests`).
  - Insert a row into a new `feedback` table.
  - Skip DB write in `NODE_ENV === "development"` (parity with `/api/requests` telemetry).
  - Return `{ ok: true }`.
- **Schema (Drizzle, new migration)**:

  ```
  feedback
    id             uuid pk
    user_id        fk -> users.id
    type           text  -- 'rating' | 'star'
    value          text  -- for star: 'starred' | 'already_given' | 'dismissed'; for rating: '1'..'5'
    version        text
    system         text
    repo_name      text
    created_at     timestamptz default now()
  ```

  `value` is stored as text (not int) so `star` outcomes are self-describing and don't collide with rating integers.

- **Rating `0` (dismiss) is not sent**. Dismissal is purely local — nothing useful to report to the server, and we explicitly don't want to log "users who said no".
- **Star `0`/`2`/`3`/`1`** all send, with `value` mapped to the four outcome strings.

### Client: request shape to `/api/feedback`

New function in a new `feedback/` package (mirroring `groq/`):

- Reuses the same `git.GetGitName`, `git.GetGitEmail`, `git.GetRepoName`, `version.Get`, `runtime.GOOS` calls as `groq.sendRequest`.
- Short HTTP timeout (~5s). Failures are swallowed silently — a flaky network must never disrupt the commit flow.
- Runs synchronously inline with the prompt handler (it's already post-commit, so a brief wait is acceptable), but wrapped in `ui.WithSpinner` is unnecessary — a silent inline call is fine.

### Where prompts plug in to the flow

- New package `prompts/` with entry point `prompts.MaybeShow(cfg, state)` called once at the end of `ExecuteCommit` (only on success).
- The helper returns nothing; all side effects (printing, state mutation, HTTP) happen internally.
- `app.Run` / existing TUI model does **not** need to know about prompts — they run after the Bubble Tea program has exited, as plain terminal I/O below the final TUI render.

### Eligibility flow (pseudocode summary)

1. Commit succeeds.
2. Load state → increment `commit_count` → save.
3. If `!cfg.Prompts.Enabled` → return.
4. If not a TTY or CI detected → return.
5. If `commit_count < 5` → return.
6. If both `rating.status != pending` and `star.status != pending` → return.
7. Roll random float; if `>= 0.15` → return.
8. Choose pending prompt (one, or random if both).
9. Show prompt, read keypress, mutate state, POST if applicable.

## Out of Scope

- Text feedback / free-form feature request capture. Rating is purely 1–5. If we later want qualitative feedback, a separate `diny feedback` command is preferred (explicit, not opportunistic).
- Re-asking users periodically (e.g., every 90 days, or after major version bump). v1 is "asked once, done forever".
- A web dashboard or admin UI for viewing ratings. Reading the DB directly is fine for v1.
- NPS-style scoring (0–10, detractor/promoter math). Keep it 1–5.
- Randomized A/B variants of the prompt copy. One copy, tuned later if needed.
- Sending rating dismissals (`0`) to the backend. Local only.
- Auth or API keys on `/api/feedback`. Trust the email-upsert pattern used by `/api/requests`. Spam mitigation is deferred.
- Showing an aggregate rating inside the CLI (e.g., `diny --about` showing average stars). Out of scope for this PRD.
- Changing the commit flow UI, regeneration, drafts, or any existing behavior.
- Showing prompts on `diny timeline`, `diny yolo`, or any non-`diny commit` command. Only the `commit` flow is instrumented in v1.
- Migration of any existing users — no prior state exists.
- Shipping a "promo / release note" prompt using the same framework. The `prompts/` package should be written with this in mind, but the release-note prompt is a later PRD.

## Further Notes

- **`0` and `2` on the star prompt are functionally identical** under the chosen re-ask policy. They are both preserved because "Close" and "Don't ask again" are two different user intents linguistically, and offering both keeps the prompt self-explanatory. If telemetry later shows the distinction doesn't matter, consolidate in a v2.
- **Rating dismissal is deliberately un-telemetered.** Logging "users who refused to rate" is lower-value than respecting the dismissal cleanly — and avoids the optics of "I said no and you logged that".
- **Probability tuning** should be revisited once there's real data. 15% on a power user who commits 20×/day is ~3 rolls/day on average, meaning the first prompt likely appears within the first working day after crossing the 5-commit gate. If that feels too aggressive in practice, drop to 10%.
- **Commit counter never resets.** Even if a user uninstalls and reinstalls, a fresh `state.yaml` means they re-enter the gate at 0. This is acceptable — fresh install = fresh start.
- **Failure of the `open`/`xdg-open` command** (e.g., headless server) should print the GitHub URL and still mark the state as `starred` — the user expressed intent, even if the browser didn't open. The alternative (leaving as `pending` and re-rolling later) conflicts with the re-ask policy.
- **TUI vs plain prompt.** Prompts are rendered with `ui.Box` for visual consistency but use a simple blocking keypress read, not a Bubble Tea program. This keeps the `app.Run` exit clean and avoids re-entering the event loop for a single keystroke.
- **Future reuse.** The `prompts/` package should expose `ShowRating(...)` / `ShowStar(...)` as independent functions so that other commands (timeline, yolo) can opt in later by calling them directly, sharing the state file but with command-specific eligibility gates.
- **Naming.** `feedback` is used for the backend endpoint and table to leave conceptual room for future qualitative feedback; the client package can be named `feedback` too. The in-CLI subsystem is `prompts/` because its job is broader than "feedback" (it also handles the star nudge, and likely future prompts).
