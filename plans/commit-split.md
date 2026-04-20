# Plan: Commit Split

> Source PRD: `diny/prd/commit-split.md`

## Architectural decisions

Durable decisions that apply across all phases:

- **Entry point**: new action on the existing `diny commit` ready screen labeled "Split into multiple commits". No new subcommand, no new CLI flag. Always visible regardless of diff size.
- **Scope**: only files returned by `git diff --cached --name-status` are in scope. Unstaged and untracked files are ignored. No file picker before splitting.
- **Granularity**: file-level. Every staged file belongs to exactly one group. No hunk-level splitting in v1.
- **API route**: reuse existing `POST /api/requests` with a new `type: "split"` value. No new endpoint.
- **Request payload**: same envelope as `commit`/`timeline` — `type`, `userPrompt` (the staged diff), `config`, `version`, `name`, `email`, `repoName`, `system`.
- **Response shape for `type: "split"`**: `data.groups: [{ order: number, type: string, message: string, files: string[] }]`. Groups arrive pre-ordered by logical dependency. Existing `commit`/`timeline` responses keep returning `data.message`.
- **DB migration**: additive — extend the `request_type` Postgres enum from `["commit", "timeline"]` to `["commit", "timeline", "split"]`. Existing rows untouched. The user runs `pnpm db:generate` and `pnpm db:migrate` in `diny-next/` manually after the Drizzle schema is updated.
- **Telemetry**: split requests log to the existing `requests` table with `type = "split"`. `suggestion` stores the serialized plan as JSON. No new tables, no new env vars.
- **System prompt**: new branch in `buildSystemPrompt` for `"split"` that honors the same user config (conventional, emoji, tone, length, custom instructions) as the single-commit branch.
- **Commit execution**: sequential `git reset` → `git add <group files>` → `git commit` per group. Stop-on-failure; successful commits are not rolled back.
- **Flag parity**: `--no-verify` applies to every commit in the sequence. `--push` pushes once after the final commit. `--print` is incompatible with Split and is rejected with an error.
- **Release ordering**: deploy `diny-next` before releasing the CLI version that ships the Split option so new CLIs always find a server that understands `type: "split"`.

---

## Phase 1: Tracer — end-to-end split with static plan

**User stories**: 1, 2, 3, 4, 5, 6, 7, 14, 19, 20, 21, 22, 23, 24, 26, 28

### What to build

The thinnest vertical slice that takes a user from the ready screen to multiple commits on disk.

- Extend the `request_type` Postgres enum in the Drizzle schema to include `"split"`. The user runs `pnpm db:generate` and `pnpm db:migrate` manually after the schema change lands.
- Accept `type: "split"` in the server's Zod schema and the `/api/requests` handler.
- Add a `"split"` branch to `buildSystemPrompt` that instructs the LLM to group staged files by concern, return 1–N groups ordered by logical dependency, and produce a commit message per group honoring the user's config.
- Return `data.groups[]` (order, type, message, files) for `type: "split"`. Log the serialized plan to the `requests` telemetry table under `type = "split"`.
- In the CLI, add a "Split into multiple commits" action to the ready screen. Picking it calls `/api/requests` with `type: "split"`, shows a spinner while waiting, and transitions the TUI into a new "plan view" state.
- Render the plan as a read-only list: one entry per group showing `[order/total] type — first line of message` and the group's files with git status prefix (`M`, `A`, `D`, `R`).
- Provide navigation (`↑`/`↓`, `k`/`j`), expand/collapse (`enter`/`space`), a "Confirm all" footer action, and cancel (`esc`/`q`) that returns to the single-commit ready screen with the original message intact.
- On confirm, execute each group in order: `git reset`, `git add <files>`, `git commit -m <message>`. No flag handling yet (no `--no-verify`, no `--push`). On any failure, surface a plain error and stop.
- Validate the plan on the CLI side before execution: every staged file appears in exactly one group. Reject with an error on duplicates or missing files.

### Acceptance criteria

- [ ] Drizzle migration adds `"split"` to the `request_type` enum; existing rows and inserts unaffected.
- [ ] `POST /api/requests` with `type: "split"` returns `data.groups[]` in dependency order.
- [ ] Request row is written with `type = "split"` and the serialized plan in `suggestion`.
- [ ] `diny commit` ready screen shows a "Split into multiple commits" action on every run.
- [ ] Picking Split renders a plan view that lists groups with order, type, message, and files-with-status.
- [ ] "Confirm all" produces N commits on disk, one per group, in the returned order.
- [ ] A single-group plan produces exactly one commit, equivalent in outcome to a normal `diny commit`.
- [ ] Cancelling the plan view returns the user to the ready screen with the originally generated single-commit message intact.
- [ ] Any mismatch between plan files and currently staged files surfaces an error instead of committing.

---

## Phase 2: Failure handling and flag parity

**User stories**: 15, 16, 17, 18, 29

### What to build

Make the split flow robust and consistent with the existing `diny commit` flags.

- Apply `--no-verify` to every `git commit` in the sequence when the flag is set on the parent `diny commit` invocation.
- After the last group commits successfully, run a single `git push` if `--push` was set. No pushes happen if the sequence stops early.
- Reject `--print` + Split up front with a clear error message and return to the ready screen.
- Replace the plain error from Phase 1 with a structured failure report: the hash(es) of commits that already landed, which group failed (order, type, first line), the stderr from `git commit`, the files left staged (the failed group's files), and the files belonging to groups that never ran.
- On sequence completion (success path), show a success screen that lists every commit made: `<short hash> <first line of message>`.
- Honor the existing `hash_after_commit` config: on success, copy the final commit's short hash to the clipboard, matching single-commit behavior.

### Acceptance criteria

- [ ] `diny commit --no-verify` → Split produces N commits, each made with `git commit --no-verify`.
- [ ] `diny commit --push` → Split runs exactly one `git push` after the final commit; no push fires if any commit fails.
- [ ] `diny commit --print` → Split is refused with a clear error and no commits are made.
- [ ] A pre-commit hook failure on group K stops the sequence, leaves groups 1..K-1 committed, leaves group K's files staged, and surfaces a report naming the failing group and showing git's stderr.
- [ ] On full success, a completion screen lists each commit's short hash and first message line.
- [ ] When `hash_after_commit` is true, the final commit's short hash is on the clipboard after a successful split.

---

## Phase 3: Plan editing (edit message + reassign files)

**User stories**: 8, 9, 12, 13, 25, 30

### What to build

Turn the read-only plan view into an editable one so users can correct the AI's output without regenerating.

- Add an `e` action that opens the focused group's message in `$EDITOR`, reusing the existing single-commit edit implementation. On editor exit, update the focused group's message in place.
- Add an `m` action that enters "move mode": focus shifts to the focused group's file list; arrow keys (`↑`/`↓`, `k`/`j`) pick a file; pressing a digit (`1`–`9`) reassigns the picked file to the group with that order number. `esc` exits move mode without a reassignment.
- Fall back gracefully when group count > 9: pressing `m` uses arrow-key destination selection instead of digit shortcuts.
- Auto-delete any group that becomes empty after a reassignment and renumber the remaining groups' `order` values so the sequence stays contiguous.
- Extend focus highlighting, move-mode cursor, and footer key hints through the active theme's primary/border/text colors. No hard-coded colors.

### Acceptance criteria

- [ ] Pressing `e` on a focused group opens `$EDITOR` with that group's current message; saving updates the message in the plan view and in the eventual commit.
- [ ] Move mode lets the user reassign any file from one group to another via digit shortcuts (≤9 groups) or arrow selection (>9 groups).
- [ ] Reassigning the last file out of a group deletes that group and renumbers the remaining groups so `order` stays contiguous starting at 1.
- [ ] Group focus, file focus in move mode, and footer hints all render through the current theme's colors.
- [ ] Edits and reassignments persist through confirm: the committed sequence reflects the final edited plan, not the plan originally returned by the API.

---

## Phase 4: Regeneration (from scratch + with feedback)

**User stories**: 10, 11, 27

### What to build

Let users steer the AI when the whole plan is off, not just individual pieces.

- Add an `r` action in the plan view that regenerates the plan from scratch: shows the loading spinner, calls `/api/requests` with `type: "split"` again, replaces the current plan on success.
- Add an `f` action that opens a text input for free-form plan-level feedback. Submitting regenerates the plan with the rejected plan(s) included in the request context so the LLM doesn't repeat the same grouping.
- Extend the request payload (and Zod schema) with an optional `previousPlans` array. Server-side, thread the previous plans into the system prompt for `type: "split"` so the LLM is explicitly told what was rejected.
- Maintain a client-side history of rejected plans across feedback rounds within a single session so successive regenerations all see prior rejections.
- During regeneration, keep the user on the plan view with a spinner overlay; on success, replace the plan and reset the focus cursor to the first group.

### Acceptance criteria

- [ ] `r` triggers a fresh plan request and replaces the current plan on success.
- [ ] `f` opens a text input, accepts free-form feedback, and triggers a regeneration that includes the current plan in `previousPlans`.
- [ ] Successive feedback rounds accumulate prior plans in `previousPlans` so the LLM sees the full rejection history.
- [ ] The server's Zod schema accepts an optional `previousPlans` field on `type: "split"` requests and the system prompt for split incorporates them.
- [ ] Regeneration shows a loading state and does not require the user to leave the plan view.
