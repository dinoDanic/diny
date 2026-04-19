package timeline

import (
	"fmt"
	"strings"
	"time"
)

type datePreset struct {
	name  string
	start time.Time
	end   time.Time
}

const presetCount = 6
const dateMenuCount = presetCount + 2 // presets + Specific date + Date range

func resolvePresets(now time.Time) []datePreset {
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	yesterday := today.AddDate(0, 0, -1)

	// This week: most recent Monday (today if today is Monday)
	daysFromMonday := (int(today.Weekday()) - int(time.Monday) + 7) % 7
	thisWeekStart := today.AddDate(0, 0, -daysFromMonday)

	// Last week: previous ISO week's Monday through Sunday
	lastWeekEnd := thisWeekStart.AddDate(0, 0, -1)
	lastWeekStart := lastWeekEnd.AddDate(0, 0, -6)

	// This month: 1st of current month through today
	thisMonthStart := time.Date(today.Year(), today.Month(), 1, 0, 0, 0, 0, today.Location())

	// Last month: 1st of previous month through last day of previous month
	lastMonthEnd := thisMonthStart.AddDate(0, 0, -1)
	lastMonthStart := time.Date(lastMonthEnd.Year(), lastMonthEnd.Month(), 1, 0, 0, 0, 0, today.Location())

	return []datePreset{
		{name: "Today", start: today, end: today},
		{name: "Yesterday", start: yesterday, end: yesterday},
		{name: "This week", start: thisWeekStart, end: today},
		{name: "Last week", start: lastWeekStart, end: lastWeekEnd},
		{name: "This month", start: thisMonthStart, end: today},
		{name: "Last month", start: lastMonthStart, end: lastMonthEnd},
	}
}

func formatDateShort(t time.Time) string {
	return t.Format("Jan 2")
}

func formatPresetLabel(p datePreset) string {
	if p.start.Equal(p.end) {
		return fmt.Sprintf("%s (%s)", p.name, formatDateShort(p.start))
	}
	return fmt.Sprintf("%s (%s \u2013 %s)", p.name, formatDateShort(p.start), formatDateShort(p.end))
}

func presetDateRange(p datePreset) string {
	if p.name == "Today" {
		return "today"
	}
	if p.start.Equal(p.end) {
		return strings.ToLower(p.name)
	}
	return fmt.Sprintf("%s (%s \u2013 %s)", strings.ToLower(p.name), formatDateShort(p.start), formatDateShort(p.end))
}

func dateMenuLabels() []string {
	presets := resolvePresets(time.Now())
	labels := make([]string, 0, dateMenuCount)
	for _, p := range presets {
		labels = append(labels, formatPresetLabel(p))
	}
	labels = append(labels, "Specific date", "Date range")
	return labels
}
