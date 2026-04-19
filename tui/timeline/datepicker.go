package timeline

import (
	"fmt"
	"time"
)

// Field indices for the date picker.
const (
	fieldYear  = 0
	fieldMonth = 1
	fieldDay   = 2
)

const minYear = 2000

type datePicker struct {
	year  int
	month int // 1–12
	day   int // 1–daysInMonth
	focus int // fieldYear, fieldMonth, or fieldDay

	// label shown above the picker (e.g. "Pick a date", "Pick start date")
	label string
	// error message shown below the picker (e.g. range validation)
	errMsg string
}

func newDatePicker(label string) datePicker {
	now := time.Now()
	return datePicker{
		year:  now.Year(),
		month: int(now.Month()),
		day:   now.Day(),
		focus: fieldDay,
		label: label,
	}
}

func (dp *datePicker) maxYear() int {
	return time.Now().Year()
}

// clampDay ensures day is valid for the current year/month.
func (dp *datePicker) clampDay() {
	max := daysInMonth(dp.year, dp.month)
	if dp.day > max {
		dp.day = max
	}
	if dp.day < 1 {
		dp.day = 1
	}
}

// adjust changes the focused field by delta, wrapping at bounds.
func (dp *datePicker) adjust(delta int) {
	dp.errMsg = ""
	switch dp.focus {
	case fieldYear:
		dp.year += delta
		maxY := dp.maxYear()
		// Wrap around
		for dp.year > maxY {
			dp.year = minYear + (dp.year - maxY - 1)
		}
		for dp.year < minYear {
			dp.year = maxY - (minYear - dp.year - 1)
		}
		dp.clampDay()

	case fieldMonth:
		dp.month += delta
		for dp.month > 12 {
			dp.month -= 12
		}
		for dp.month < 1 {
			dp.month += 12
		}
		dp.clampDay()

	case fieldDay:
		max := daysInMonth(dp.year, dp.month)
		dp.day += delta
		for dp.day > max {
			dp.day -= max
		}
		for dp.day < 1 {
			dp.day += max
		}
	}
}

// moveFocus shifts focus left or right, wrapping.
func (dp *datePicker) moveFocus(delta int) {
	dp.focus += delta
	if dp.focus > fieldDay {
		dp.focus = fieldYear
	}
	if dp.focus < fieldYear {
		dp.focus = fieldDay
	}
}

// toTime returns the currently selected date as a time.Time.
func (dp *datePicker) toTime() time.Time {
	return time.Date(dp.year, time.Month(dp.month), dp.day, 0, 0, 0, 0, time.Now().Location())
}

// dateString returns the date in YYYY-MM-DD format.
func (dp *datePicker) dateString() string {
	return fmt.Sprintf("%04d-%02d-%02d", dp.year, dp.month, dp.day)
}

// daysInMonth returns the number of days in the given month/year.
func daysInMonth(year, month int) int {
	// Day 0 of the next month gives the last day of this month.
	return time.Date(year, time.Month(month+1), 0, 0, 0, 0, 0, time.UTC).Day()
}
