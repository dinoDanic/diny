# Plan: Timeline Date Picker

> Source PRD: `../prd/timeline-date-picker.md`

## Architectural decisions

Durable decisions that apply across all phases:

- **Scope**: all changes live in the `diny` CLI, primarily under the `tui/timeline` package. No changes to `diny-next`, the `/api/requests` endpoint, the Groq request shape, or the `git` package's commit-fetch functions.
- **Menu order**: presets first (Today, Yesterday, This week, Last week, This month, Last month), custom options last (Specific date, Date range).
- **Week boundaries**: ISO 8601 — weeks start on Monday. Hard-coded, not locale-aware.
- **Preset labels**: each preset displays its resolved date range (e.g., `This week (Apr 13 – Apr 19)`), computed from the current local date at menu render time.
- **Picker model**: three adjustable fields — Year, Month, Day. Defaults to today's date with Day focused. Year bounded 2000 through current year.
- **Picker keys**: ←/→ (h/l) move between fields; ↑/↓ (j/k) adjust by 1 with wrap; PgUp/PgDn or Shift+↑/↓ adjust by 10; Enter confirm; Esc back; q / Ctrl+C quit.
- **Invalid dates**: Day auto-clamps to the last valid day of the current Month/Year after any field change (Feb 29 in non-leap year → Feb 28).
- **Range flow**: same picker used twice. End-date picker pre-fills today. Start-after-end is rejected inline on the end-picker; user stays on that step with an error message.
- **Legacy text input**: removed for date entry. The `textinput.Model` remains in use for the feedback-refinement flow.
- **State machine**: existing `stateEnterDate` / `stateEnterStartDate` / `stateEnterEndDate` are replaced by `statePickDate` / `statePickStartDate` / `statePickEndDate`. Preset handlers resolve a concrete range and jump straight to `stateFetching`.
- **Backend unchanged**: `git.GetCommitsByDate`, `git.GetCommitsByDateRange`, and `groq.CreateTimelineWithGroq` are called exactly as today. Only the `dateRange` human-readable string and the resolved `start`/`end` inputs change shape.

---

## Phase 1: Preset menu expansion

**User stories**: 1, 2, 3, 4, 5, 6, 7, 26

### What to build

Expand the timeline's top-level menu from three entries to eight. The five new presets — Yesterday, This week, Last week, This month, Last month — each resolve to a concrete `[start, end]` date range computed from today's local date, and each menu label shows those resolved dates alongside the preset name. Selecting any preset skips all date-entry UI and goes straight to the existing fetch-and-analyze flow.

Custom Specific date and Date range remain at the bottom of the menu, wired to the existing text input for now; they are rewritten in later phases. The "Analyze different period" action on the results screen returns users to this expanded menu.

### Acceptance criteria

- [ ] Menu lists eight entries in the documented order with arrow/j-k navigation.
- [ ] Each preset label shows its resolved dates in a human-readable form (e.g., `Yesterday (Apr 18)`, `This week (Apr 13 – Apr 19)`).
- [ ] Selecting Today, Yesterday, This week, Last week, This month, or Last month fetches commits for that period and produces an analysis without any further prompts.
- [ ] Week-based presets treat Monday as the first day of the week.
- [ ] Month-based presets use calendar-month boundaries (1st through last-of-month for Last month; 1st through today for This month).
- [ ] Preset labels reflect the current local date each time the command starts.
- [ ] "Analyze different period" from the results screen returns to this menu, not to any legacy flow.

---

## Phase 2: Arrow-key picker for single custom date

**User stories**: 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 22, 23, 24, 25

### What to build

Replace the text-input flow behind the "Specific date" menu entry with an arrow-key picker. The picker displays three columns — Year, Month, Day — pre-filled with today's date and with Day focused. Users move focus between columns with left/right, adjust the focused value with up/down (wrapping at bounds), and jump by 10 with PageUp/PageDown or Shift+Up/Down. After any field change, the Day is clamped to the last valid day of the current Month/Year so the picker never shows Feb 31 or Feb 29 outside a leap year.

The focused column is visually highlighted using the active theme's primary color; the rendering honors the existing theme system and remains legible on narrow terminals. A footer shows the active key hints. Enter confirms the picked date and flows into the existing fetch-and-analyze path; Esc returns to the preset menu; q / Ctrl+C quits. The results screen shows the picked date in its header, matching how presets already display their range.

### Acceptance criteria

- [ ] Selecting "Specific date" from the menu opens the arrow-key picker instead of a text input.
- [ ] Picker opens with today's Year/Month/Day pre-filled and Day focused.
- [ ] Left/right (and h/l) cycle focus between Year, Month, and Day.
- [ ] Up/down (and k/j) increment/decrement the focused field by 1 and wrap at bounds (year 2000 ↔ current year, month 1 ↔ 12, day 1 ↔ days-in-month).
- [ ] PageUp/PageDown and Shift+Up/Shift+Down adjust the focused field by 10.
- [ ] Changing Month or Year auto-clamps Day to the last valid day of that Month/Year (Feb 29 in a non-leap year → Feb 28).
- [ ] Enter confirms the picked date and advances to the existing fetch → analysis → results flow.
- [ ] Esc returns to the preset menu.
- [ ] Footer help line lists the active keys (←/→ move, ↑/↓ adjust, PgUp/Dn jump, Enter confirm, Esc back).
- [ ] Rendering respects the active theme's primary color for the focused column and stays legible on narrow terminals.
- [ ] Results screen shows the picked date in its header so the user can confirm the analyzed period.

---

## Phase 3: Range picker and cleanup

**User stories**: 18, 19, 20, 21, 27

### What to build

Wire the arrow-key picker from Phase 2 into the "Date range" menu entry. Users pick a start date, then an end date, using the same interaction. The end-date picker pre-fills with today. On confirming the end date, validate that start ≤ end; an inverted range shows an inline error on the end-picker, preserves the user's entry, and blocks advancement until corrected. Esc from the end picker returns to the start picker with its previous value restored.

With both custom flows migrated to the picker, delete the legacy date-entry code: the text-input based date prompts, the huh-select day/month/year helpers, and any associated state machine entries. The `textinput.Model` remains for the feedback-refinement flow only.

### Acceptance criteria

- [ ] Selecting "Date range" opens the arrow-key picker for the start date.
- [ ] After confirming the start date, the end-date picker opens pre-filled with today's date and Day focused.
- [ ] Esc on the end-date picker returns to the start-date picker with its previous value intact.
- [ ] Confirming an end date earlier than the start date shows an inline error on the end-picker and does not advance.
- [ ] A valid range flows into the existing fetch → analysis → results path, with the header showing the picked range.
- [ ] Old date-entry code paths (text-input prompts for specific date / start date / end date, and the day/month/year helper lists) are removed from the codebase.
- [ ] No lingering references to the removed states or helpers remain in the timeline package.
- [ ] `go build ./...` succeeds and running `diny timeline` exercises all eight menu paths without errors.
