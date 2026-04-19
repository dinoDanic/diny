# Timeline Date Picker

## Problem Statement

The `diny timeline` command forces users to type dates character-by-character in `YYYY-MM-DD` format whenever they want to analyze anything other than today. Picking a specific date or a date range means typing one or two dates by hand, with no defaults, no shortcuts for common periods like "yesterday" or "this month", and no way to tweak a date with arrow keys. This is tedious for the most common use cases (checking recent work), error-prone (a single typo or wrong format yields no commits or a confusing error), and feels out of step with how the rest of the TUI works, which is almost entirely keyboard-driven with arrow navigation.

## Solution

Replace the text-entry date flow with two things:

1. A richer preset menu that covers the common timeframes (today, yesterday, this week, last week, this month, last month) so most users never need to touch a calendar at all. Each preset label shows the dates it resolves to, so users can verify at a glance.
2. An arrow-key date picker for the "Specific date" and "Date range" options, with three adjustable fields (Year / Month / Day). Users move between fields with left/right, adjust values with up/down (wrapping at bounds), and jump by 10 with PageUp/PageDown or Shift+Up/Down. The picker pre-fills with today's date and focuses the Day field, so the common case of "a few days ago" is just a couple of keystrokes.

Date ranges use the same picker twice (start, then end), with the end field defaulting to today. Invalid states like Feb 31 auto-clamp as the user moves through fields, and a start-after-end range is rejected on confirm with a clear error.

## User Stories

1. As a developer running `diny timeline`, I want to pick "Yesterday" from a menu, so that I can review what I committed yesterday without typing a date.
2. As a developer, I want a "This week" preset that covers Monday through today, so that I can summarize my week in one keystroke.
3. As a developer, I want a "Last week" preset, so that I can generate a summary for weekly standups or status updates.
4. As a developer, I want a "This month" preset, so that I can review my month's work without thinking about exact dates.
5. As a developer, I want a "Last month" preset, so that I can generate a monthly report for the prior month.
6. As a developer, I want each preset label to also show the resolved date range (e.g., "This week (Apr 13 – Apr 19)"), so that I know exactly what the preset means before committing to it.
7. As a developer, I want the preset menu listed first and custom options ("Specific date", "Date range") listed at the bottom, so that the fastest choices are closest to my fingers.
8. As a developer picking a specific custom date, I want Year / Month / Day columns pre-filled with today's date, so that I only change what I need to.
9. As a developer picking a specific date, I want the Day field focused first, so that I can nudge the date a few days back without tabbing across fields.
10. As a developer, I want to move between Year, Month, and Day fields with left/right arrow keys, so that navigation matches my mental model of the picker as three columns.
11. As a developer, I want to adjust the focused field's value with up/down arrow keys, so that I can mutate one piece of the date at a time.
12. As a developer, I want up/down to wrap around at bounds (month 12 + 1 → 1, day 31 + 1 → 1), so that I can cycle through values quickly without hitting a wall.
13. As a developer, I want PageUp/PageDown (or Shift+Up/Shift+Down) to jump the focused field by 10, so that I can reach distant dates in a few keystrokes.
14. As a developer, I want invalid date combinations (e.g., switching to Feb while Day=31) to auto-clamp the Day to the last valid day of that month, so that the picker never shows an invalid date.
15. As a developer, I want the same auto-clamp behavior when I change the Year across a Feb 29 leap boundary, so that I never end up with Feb 29 in a non-leap year.
16. As a developer confirming a date, I want to press Enter to accept the picked date, so that confirmation is consistent with the rest of the TUI.
17. As a developer, I want to press Esc to go back to the preset menu from the picker, so that I can change my mind without quitting.
18. As a developer picking a custom date range, I want to pick start and end dates using the same arrow picker I used for single dates, so that I only have to learn one interaction.
19. As a developer, I want the end date to pre-fill as today when I enter the end picker, so that the common "from date X to now" range is a single keystroke.
20. As a developer, I want a start-after-end range to be rejected with a clear error message, so that I don't accidentally run an analysis on an inverted range.
21. As a developer, I want a rejected range to keep me on the end-date picker with my entry intact, so that I can correct the error without starting over.
22. As a developer, I want the footer help line in the picker to show the active key hints (←/→ move, ↑/↓ adjust, PgUp/Dn jump, Enter confirm, Esc back), so that I can discover the picker's capabilities without reading docs.
23. As a developer on a small terminal, I want the picker to remain legible when the window is narrow, so that I can use it in split panes.
24. As a developer using any of the existing themes (Catppuccin, Tokyo Night, Nord, etc.), I want the picker to honor my theme colors, so that the UI stays visually consistent.
25. As a developer, I want to see what date range the analysis ran against in the results screen, so that I can confirm the summary is for the period I meant.
26. As a developer, I want to be able to pick "Analyze different period" from the results screen and land back on the same preset menu, so that I can iterate on timeframes quickly.
27. As a developer, I want the old typed-date input to be fully removed, so that I'm not presented with a stale or inconsistent alternative flow.

## Implementation Decisions

### Menu structure

The top-level timeline menu replaces today's three options (`Today`, `Specific date`, `Date range`) with an expanded list in this order:

1. Today
2. Yesterday
3. This week
4. Last week
5. This month
6. Last month
7. Specific date
8. Date range

Presets 1–6 resolve to a concrete `[start, end]` range at render time and display their resolved dates in the label (e.g., `This week (Apr 13 – Apr 19)`).

### Preset resolution rules

- **Today**: start = today, end = today.
- **Yesterday**: start = today − 1, end = today − 1.
- **This week**: start = most recent Monday (today if today is Monday), end = today.
- **Last week**: start = previous ISO week's Monday, end = previous ISO week's Sunday.
- **This month**: start = 1st of current month, end = today.
- **Last month**: start = 1st of previous month, end = last day of previous month.

Week boundaries use Monday as the first day of the week (ISO 8601).

### Arrow-key date picker

New picker component with three fields: Year, Month, Day. State includes the current `time.Time`, the focused field index, and whether the picker is in "start" or "end" mode for ranges.

- **Layout**: horizontal three-column layout with the focused column visually highlighted.
- **Defaults**: pre-fill with today's date; focus is on the Day field.
- **Keys**:
  - `←` / `→` or `h` / `l`: move focus between fields (wraps at ends).
  - `↑` / `↓` or `k` / `j`: increment / decrement the focused field by 1, wrapping at field bounds (year has a sane lower bound, e.g., 2000, and upper bound of current year).
  - `PgUp` / `PgDn` and `Shift+↑` / `Shift+↓`: increment / decrement by 10.
  - `Enter`: confirm the picked date and advance the flow.
  - `Esc`: return to the preset menu (or to the start-date picker when currently on end).
  - `q` / `Ctrl+C`: quit.
- **Auto-clamp**: after any field change, the Day is clamped to the last valid day of the current Month/Year combination. Feb 29 in a non-leap year clamps to Feb 28.
- **Bounds**: Year wraps between 2000 and the current year (configurable constant in the package). Month wraps 1–12. Day wraps 1–`daysInMonth(year, month)`.

### Date-range flow

Range mode reuses the picker in two consecutive steps:

1. Start-date picker (pre-filled today, Day focused).
2. End-date picker (pre-filled today, Day focused).

On confirming the end date, validate `start <= end`. If the range is inverted, keep the user on the end picker, surface an inline error message, and do not advance.

### State machine changes (tui/timeline)

- Existing states `stateEnterDate`, `stateEnterStartDate`, `stateEnterEndDate` are replaced with `statePickDate` (single custom date) and `statePickStartDate` / `statePickEndDate` (range).
- `dateMenuItems` grows from three to eight entries; `dateCursor` logic stays the same but ranges over the larger list.
- `confirmDateChoice` branches into six preset handlers (compute range then go straight to `stateFetching`) and two custom-picker handlers.
- A new picker sub-model encapsulates the Year/Month/Day state plus its own `Update` and `View`. The timeline model embeds it for `statePickDate`, `statePickStartDate`, and `statePickEndDate`.
- The existing text-input (`textinput.Model`) usage for date entry is removed. It remains for the feedback-refinement flow.

### Rendering

- Picker renders three columns with consistent spacing, the focused column bordered or highlighted using the active theme's primary color.
- Footer shows the picker-specific key hints, replacing the current "enter confirm / esc back" footer when the picker is active.
- Preset labels are generated each time the timeline command starts (so "this week" reflects the current date).

### Backend / git layer

No changes to `git.GetCommitsByDate`, `git.GetCommitsByDateRange`, or the Groq request shape. Preset handlers resolve to `start`/`end` strings in the existing format and call the same functions. The `dateRange` string shown in the analysis header and filename continues to use the existing human-readable form (e.g., `today`, `2026-04-13 to 2026-04-19`, `specific-date-2026-04-15`), with preset names preserved so saved files remain self-describing.

### Removal

The old text-entry approach is removed entirely. No toggle, no fallback key. The `dateInputPrompt`, `generateDayOptions`, `generateMonthOptions`, `generateYearOptions` helpers (and any related dead code) are deleted.

## Out of Scope

- Adding "This year", "Last year", "Last 7 days", "Last 30 days", or any other presets not listed above. They can be added later if the menu proves too narrow in practice.
- A full mini-calendar grid view (Mo–Su layout with highlighted day). The arrow-field picker is intentionally simpler.
- Locale-aware week start (Sunday vs Monday). Monday is hard-coded.
- Locale-aware date formatting in the picker or in labels. Dates remain in unambiguous `YYYY-MM-DD` or `Mon DD` form in English.
- Keyboard shortcuts on the preset menu (e.g., pressing `1`–`6` to jump straight to a preset). Arrow/j/k + Enter stays the only selection method.
- Changes to the Groq prompt, the analysis output, the results screen's action keys (c/s/r/f/n/q), or the feedback-refinement flow.
- Any change to the server API (`/api/requests`) or telemetry schema.
- Tests. Existing coverage is minimal, and adding tests here is not a blocker.

## Further Notes

- `dateMenuItems` growing from 3 to 8 means the menu Height used by huh-style rendering needs re-checking; the current hand-rolled menu renders all items directly and just needs the cursor clamp to respect the new length.
- The picker can be generalized later (e.g., for a "since last tag" comparison feature or for scheduling future reminders), but the first version should live under `tui/timeline` alongside the rest of the timeline TUI. Promoting it to a shared package is deferred until a second call site exists.
- The year lower bound (2000) is somewhat arbitrary but avoids letting the user scroll to year 0. If a user genuinely needs an older date, the bound can be lowered without affecting the rest of the UI.
- Because the picker auto-clamps Day when Month or Year changes, users who start on Day and nudge up into the next month get predictable behavior: Day wraps within the current month's day count, not into the next month.
- Preset label formatting uses the user's local time zone, matching the behavior of the existing `git.GetCommitsToday`. No explicit TZ handling is added.
