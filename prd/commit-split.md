# Commit Split

## Problem Statement

When a developer runs `diny commit`, diny produces a single commit message for everything currently staged. In practice, real work rarely arrives in neat, single-purpose chunks — a developer might have finished a feature, fixed an unrelated typo, and renamed a few files, all in the same working session. The staged diff then contains multiple distinct concerns: a feature, a fix, a refactor, a doc update.

Today diny has no way to acknowledge that. It does one of two things:

1. Generates a single commit message that tries to summarize everything, producing vague commits like "update auth and misc cleanup".
2. Forces the developer to abandon diny, manually unstage, re-stage subsets with `git add -p` or `git add <file>`, and run `diny commit` (or `git commit`) multiple times in a row — slow, error-prone, and it defeats the reason they reached for diny in the first place.

The developer knows their staged diff covers several concerns, but diny treats it as one monolith. The result is either lumpy commits that nobody wants to review, or a manual unstage/stage/commit loop that the tool was supposed to eliminate.

## Solution

Add a "Split into multiple commits" option to the `diny commit` ready screen. When the user picks it, diny sends the staged diff to the API with a new `split` request type. The LLM returns an ordered plan: a list of groups, where each group has a proposed commit message and the set of files it covers, ordered so that dependencies land before dependents.

diny renders the plan in the TUI: one row per group showing its message and its file list. The user can:

- Edit any group's commit message in `$EDITOR`.
- Reassign a file from one group to another.
- Regenerate the whole plan with free-form feedback ("keep tests with their feature", "merge the two refactors").
- Confirm the plan, at which point diny commits each group in order by unstaging everything and re-staging group N's files before each `git commit`.

If a commit fails mid-sequence (e.g., a pre-commit hook rejects it), diny stops. Commits that already landed stay committed; files belonging to pending groups remain staged so the user can fix and resume.

The split flow honors the same flags as `diny commit`: `--no-verify` applies to every commit in the sequence, and `--push` pushes once after the final commit succeeds. Only currently staged files are in scope — split never reaches into unstaged or untracked files.

## User Stories

1. As a developer who just finished a chunk of work that touched a feature and an unrelated bug, I want diny to offer a "Split into multiple commits" option, so that I don't have to manually unstage and re-stage files to get clean history.
2. As a developer running `diny commit`, I want the split option to always appear on the ready screen, so that I don't have to remember a separate command or flag to reach for it.
3. As a developer, I want the split option to sit alongside the existing commit/edit/regenerate options, so that discovery costs me nothing — I see it every time.
4. As a developer picking "Split", I want diny to analyze my staged diff and propose a grouping by concern, so that I get an intelligent starting point instead of grouping files manually.
5. As a developer, I want each proposed group to have both a commit type (feat/fix/refactor/etc.) and a human-readable message, so that the plan reads like a real commit history.
6. As a developer, I want the proposed groups returned in the order they should be committed (dependencies first), so that intermediate commits in history are buildable and don't break bisect.
7. As a developer viewing the plan, I want to see each group's message and the files it covers, so that I can verify the split matches my intent before anything lands.
8. As a developer, I want to move a file from one group to another via the TUI, so that I can correct the AI when it puts a file in the wrong bucket.
9. As a developer, I want to edit any group's commit message in `$EDITOR`, so that I can polish the wording the same way I do for single commits.
10. As a developer, I want to provide free-form feedback like "keep tests with their feature" and get a regenerated plan, so that I can steer grouping without manual file shuffling.
11. As a developer unhappy with the whole proposed plan, I want a way to regenerate it from scratch, so that I get a different split without typing feedback.
12. As a developer, I want moving the last file out of a group to auto-delete that group, so that the plan stays consistent without manual cleanup.
13. As a developer, I want the TUI to show which group is focused and let me navigate with arrow keys, so that the interaction matches the rest of diny.
14. As a developer, I want a single "Confirm all" action to commit the whole plan, so that I approve once and diny handles the sequence.
15. As a developer, I want diny to run `--no-verify` on every commit in the sequence if I passed `--no-verify`, so that the flag works the same as it does for single commits.
16. As a developer, I want `--push` to push once after the last commit in the sequence succeeds, so that I don't end up with N pushes for N commits.
17. As a developer whose pre-commit hook fails partway through the sequence, I want diny to stop, leave the already-committed groups in place, and leave the pending groups' files staged, so that I can fix the failure and continue without losing work.
18. As a developer, I want a clear error message that names which group failed and why, so that I can diagnose the failure without re-running anything.
19. As a developer, I want only my currently staged files to be in scope for the split, so that the flow matches my staging mental model and doesn't surprise me by touching unstaged changes.
20. As a developer, I want the split option to work with my existing config (conventional commits, emoji, tone, length, custom instructions), so that every commit in the split follows the same style rules as a single-commit diny run.
21. As a developer, I want the split option to be available regardless of diff size, so that I can use it even for a two-file commit if I want separate messages.
22. As a developer, I want the TUI to let me cancel the split flow and fall back to the single-commit ready screen, so that I can back out without quitting diny entirely.
23. As a developer, I want the plan view to indicate commit order (e.g., numbered 1/2/3), so that I know what will land first, second, and third.
24. As a developer, I want each group's file list to show the git status (M/A/D/R) per file, so that I can tell modifications from additions at a glance.
25. As a developer, I want the "edit message" action to behave exactly like the existing edit action for single commits, so that I don't have to learn a new interaction.
26. As a developer with only one concern in my staged diff, I want the LLM to return a single-group plan when splitting doesn't make sense, so that the tool doesn't invent fake splits.
27. As a developer, I want free-form feedback on the plan to preserve the previous plan as history context, so that the LLM doesn't propose the same flawed split again after I reject it.
28. As a developer, I want diny to show a loading state while the split plan is being generated, so that I know the tool is working and roughly how long to wait.
29. As a developer, I want the successful completion screen to list each commit that was made (hash + first line of message), so that I can confirm the history I just produced.
30. As a developer who copied a file assignment into the wrong group, I want to undo the assignment, so that I don't have to re-shuffle everything back manually.

## Implementation Decisions

### Entry point

- The `diny commit` ready screen gets a new action: **Split into multiple commits**. It appears for every run, regardless of diff size or file count. No new CLI flag, no config toggle.
- The option behaves like any other ready-screen action (edit, regenerate, custom feedback, save draft): it transitions into a new TUI sub-state dedicated to the split flow.
- Cancelling out of the split flow returns the user to the existing single-commit ready screen with the original single commit message still generated and intact.

### Request type and API shape

- A new request `type: "split"` is added to the existing `/api/requests` endpoint. The payload keeps the same envelope (`userPrompt`, `config`, `version`, `name`, `email`, `repoName`, `system`) with the staged diff as `userPrompt`.
- The response extends the current shape. On success the API returns a structured `data.groups` array instead of (or in addition to) a `data.message` string. Each group:
  - `type`: conventional-commit type (`feat`, `fix`, `refactor`, etc.).
  - `message`: the full commit message for that group (including title and optional body, respecting the user's conventional/emoji/tone/length config).
  - `files`: an array of file paths belonging to this group. File paths match exactly what `git diff --cached --name-only` would output.
  - `order`: an integer starting at 1 indicating commit order within the plan. Groups are returned already ordered, but the field is included for explicitness.
- Errors are returned in the existing `error` string field.
- Telemetry: the existing `requests` table captures split requests under `type = "split"`. No schema change needed beyond allowing the new type value.
- Migration ownership: the user runs `pnpm db:generate` and `pnpm db:migrate` in `diny-next/` manually. The implementer updates the Drizzle schema; Claude does not run the migration commands.

### System prompt for split

- `buildSystemPrompt` in `diny-next` gets a new branch for `type: "split"`. The prompt instructs the LLM to:
  - Read the full git diff.
  - Group files by concern, where each concern is a distinct feature/fix/refactor/etc.
  - Produce between 1 and N groups (N to be tuned; start at 8 as a soft upper bound).
  - Return groups ordered by logical dependency — schemas and shared types before code that depends on them, feature code before its tests.
  - Use the user's config (conventional, emoji, tone, length, custom instructions) to produce each group's message, exactly like the single-commit prompt does.
  - Return a single-group plan when the diff really is one concern.

### Grouping granularity

- File-level only. Every staged file belongs to exactly one group. Hunk-level splitting (splitting within a single file) is out of scope for v1.
- If the LLM returns a file in zero groups, diny adds it to the last group automatically and warns the user.
- If the LLM returns a file in multiple groups, diny rejects the plan and re-requests with a hint (or shows an error to the user with a regenerate option).

### Threshold and auto-offer

- No threshold. The "Split into multiple commits" option is always present on the ready screen. The decision to split is entirely the user's.
- No config keys are added. Nothing to disable, nothing to tune. If the menu becomes noisy in practice, a future iteration can add `commit.split.enabled`.

### TUI: split plan view

- New TUI state (e.g., `stateSplitPlan`) in the commit app model. Entered from the ready screen via the new action.
- Layout: each group renders as a card with:
  - Header: `[1/3] feat — "add auth middleware"` (order, type, first line of message).
  - Body: list of files with git status prefix (`M foo.go`, `A new.go`, `D gone.go`).
  - Selected/focused card is highlighted with the active theme's primary color.
- Navigation:
  - `↑` / `↓` or `k` / `j`: move focus between groups.
  - `enter` or `space`: expand/collapse a group to see the full message body.
  - `e`: edit focused group's message in `$EDITOR` (reuses existing single-commit edit flow).
  - `m`: enter "move mode" — focus shifts to the group's file list; arrow keys pick a file; pressing a digit (`1`–`9`) reassigns the file to the corresponding group number. If the source group becomes empty, it's deleted automatically.
  - `f`: open a text input for plan-level feedback. Submitting regenerates the entire plan with the previous plan included in the context so the LLM doesn't repeat the same grouping.
  - `r`: regenerate the plan from scratch (no feedback).
  - `enter` on the global footer action **Confirm all**: proceed to commit execution.
  - `esc` / `q`: cancel and return to the single-commit ready screen.
- Footer shows the active key hints contextually (same convention as the timeline picker).

### Commit execution

When the user confirms the plan, diny executes groups in order. For each group:

1. `git reset` (mixed, default) to unstage everything. This preserves working-tree state and only clears the index.
2. `git add <files>` for every file in the current group. For deleted files, use `git add -u <path>` semantics or `git rm --cached` as appropriate (the existing single-commit path handles this via the normal index state; the split path needs to replicate it).
3. Run `git commit -m <message>` (with `--no-verify` if the flag is set, same as `TryCommit` in `commit/helpers.go`).
4. On success, record the short hash and move to the next group.
5. On failure, stop. The files already committed stay committed. The files for the current group remain staged (because we just `git add`'d them). The files for remaining groups return to their pre-split staged state (they're still in the working tree; diny re-stages them as a best-effort recovery step).

After all groups succeed:

- If `--push` was passed, run a single `git push`.
- If `hash_after_commit` is enabled in config, copy the hash of the final commit to the clipboard (matching the existing single-commit behavior).
- The success screen lists every commit that was made: `<hash> <first line>`.

### Flag parity

- `--no-verify`: applied to every `git commit` in the sequence.
- `--push`: applied once after the last successful commit. If the sequence stops early, no push happens.
- `--print`: not supported in split mode. Selecting "Split" while `--print` is active shows an error and returns to the ready screen. (`--print` is meant for piping a single message to `git commit -F`; it has no meaningful semantics for a multi-commit plan.)

### Staging scope

- Only files returned by `git diff --cached --name-status` are in scope. The existing `GetStagedFiles()` helper is the source of truth.
- Unstaged modifications and untracked files are ignored entirely. No file picker, no opt-in to include them. The user's staging decision defines the split's scope.

### Failure mode and recovery

- Stop-on-failure: diny halts on the first `git commit` failure. Successful prior commits are not rolled back.
- After a failure, diny prints:
  - The hash(es) of commits that already landed.
  - Which group failed and the stderr from git.
  - The files that remain staged (the failed group's files).
  - The files that belong to groups that never ran (left in the working tree, un-staged at that point).
- The user fixes the root cause (hook error, lint failure, etc.) and can re-run `diny commit` — they'll land back on the ready screen with the still-staged files.

### Regeneration with history

- "Regenerate with feedback" (`f`) sends the previous plan as part of the context in the next `split` request. Payload shape to be finalized server-side, but conceptually: include a `previousPlans` array in the request so the LLM sees what was rejected.
- "Regenerate from scratch" (`r`) does a fresh `split` request with no history.
- Both actions transition the TUI back into a loading state with the existing spinner; on completion the plan view re-renders.

### Telemetry

- Every split request logs to the `requests` table with `type = "split"`. `suggestion` stores the serialized plan (JSON) so the backend has the same audit trail as today's single-commit rows.
- No new tables. No new environment variables.

## Out of Scope

- Hunk-level splitting (splitting a single file's changes across multiple commits). File-level only for v1.
- Auto-triggering the split flow based on file count, folder count, or line count. The option is always manually picked.
- A `diny split` standalone command. Everything lives inside `diny commit`.
- A `--split` CLI flag on `diny commit` to skip the single-message generation. The user always sees the single message first and picks "Split" from the ready screen if they want it.
- Merging two proposed groups into a single commit. If the user wants this, they can regenerate with feedback ("merge the refactor and the fix"). A dedicated merge action can be added later.
- Per-group regeneration of a single commit message (separate from plan-level regenerate). Scope-creep for v1; the message edit action (`e`) covers ad-hoc wording changes.
- Reordering groups manually in the TUI. Order is whatever the LLM returned; users who want a different order regenerate with feedback.
- Undoing a file reassignment. Users who make a wrong move reassign the file back manually.
- Rolling back successful commits when a later commit in the sequence fails. Stop-on-failure preserves completed work.
- Atomic "all-or-nothing" mode. Same reason as above.
- Config keys for split (enabled flag, group caps, thresholds). Everything hard-coded for v1.
- Including unstaged or untracked files in the split. Staged-only.
- Working with `--print`. Incompatible with multi-commit output.
- Tests. Existing coverage in this repo is minimal, and adding tests here is not a blocker.
- Localization, accessibility features beyond what the existing TUI already supports.
- Any change to the `timeline` request path.

## Further Notes

- The LLM's grouping quality is the single biggest risk. The system prompt needs to be specific about what "a concern" means, with examples (a feat + its tests count as one concern; a feat + an unrelated typo fix are two). A round of prompt iteration on representative real-world diffs is necessary before shipping.
- The move-mode interaction (press `m`, then a digit to reassign) is the most novel TUI interaction in the flow. If group count > 9, the digit shortcut breaks; fall back to arrow-key selection of the destination group in that case.
- Because split runs `git reset` between commits, any pre-commit hook side effects (like auto-formatting that mutates files) will interact with the sequence. If a hook mutates a file, the next `git add` for that file in a later group may pick up unintended changes. Initial version: document this as a known limitation; later versions can snapshot file contents between groups.
- Splitting against a diff that already contains merge conflicts or unmerged paths should be rejected up front with a clear error.
- The existing `LAZYGIT_PENDING_COMMIT` / `COMMIT_EDITMSG` draft-save flow doesn't apply to split plans — there's no single draft file that represents a multi-commit plan. Consider a future JSON draft file if users ask for persistence across sessions.
- The API response shape change (`data.message` → `data.groups`) for `type: "split"` is additive. Existing `commit`/`timeline` callers keep their current shape. Versioning the endpoint is not required.
- Because diny supports multiple themes, the focused-group highlight, the move-mode cursor, and the footer key hints should all route through the theme's primary/border/text colors. No hard-coded colors in the split view.
- For repositories with many small staged files (100+), the TUI plan view should virtualize or scroll rather than render every file inline. First version can assume < 50 files and defer virtualization.
- Follow-up ideas worth capturing but not doing now: hunk-level split, `diny split` as a standalone command, per-group regen, manual reordering, atomic rollback, group merging, draft persistence, keyboard shortcut to expand/collapse all groups.
