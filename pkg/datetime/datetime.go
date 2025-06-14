// Package datatime provides utility functions for working with date and time data.
package datetime

import "time"

// DatetimeGetter defines the interface for date/time utility functions.
type Getter interface {
	GetMonthStartAndEnd(monthOffset int) (start, end time.Time)
	GetLastDayOfMonth(t time.Time) int
}

// Datetime is the implementation of DatetimeGetter.
type Datetime struct {
	funcTime func() time.Time
}

func NewDatetime(funcTime func() time.Time) *Datetime {
	return &Datetime{funcTime: funcTime}
}

// GetMonthStartAndEnd calculates the start and end dates of a given month offset.
func (dt *Datetime) GetMonthStartAndEnd(monthOffset int) (start, end time.Time) {
	currentTime := dt.funcTime()

	timestamp := currentTime.AddDate(0, -monthOffset, 0)
	start = time.Date(timestamp.Year(), timestamp.Month(), 1, 0, 0, 0, 0, time.Local)
	end = time.Date(timestamp.Year(), timestamp.Month(), dt.GetLastDayOfMonth(timestamp), 23, 59, 59, 0, time.Local)

	return start, end
}

// GetLastDayOfMonth returns the last day of the month for a given time.
func (dt *Datetime) GetLastDayOfMonth(t time.Time) int {
	nextMonth := t.AddDate(0, 1, 0)
	firstDayOfNextMonth := time.Date(nextMonth.Year(), nextMonth.Month(), 1, 0, 0, 0, 0, t.Location())
	lastDay := firstDayOfNextMonth.AddDate(0, 0, -1).Day()

	return lastDay
}
