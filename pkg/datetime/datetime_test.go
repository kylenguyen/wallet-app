package datetime_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"bitbucket.org/ntuclink/ff-order-history-go/pkg/datetime"
)

func TestDatetime_GetMonthStartAndEnd(t *testing.T) {
	testCases := []struct {
		name          string
		monthOffset   int
		currentTime   time.Time
		expectedStart time.Time
		expectedEnd   time.Time
	}{
		{
			name:          "Current Month",
			monthOffset:   0,
			currentTime:   time.Date(2023, 11, 15, 10, 0, 0, 0, time.Local),
			expectedStart: time.Date(2023, 11, 1, 0, 0, 0, 0, time.Local),
			expectedEnd:   time.Date(2023, 11, 30, 23, 59, 59, 0, time.Local),
		},
		{
			name:          "Last Month",
			monthOffset:   1,
			currentTime:   time.Date(2023, 11, 15, 10, 0, 0, 0, time.Local),
			expectedStart: time.Date(2023, 10, 1, 0, 0, 0, 0, time.Local),
			expectedEnd:   time.Date(2023, 10, 31, 23, 59, 59, 0, time.Local),
		},
		{
			name:          "Two Months Ago",
			monthOffset:   2,
			currentTime:   time.Date(2023, 11, 15, 10, 0, 0, 0, time.Local),
			expectedStart: time.Date(2023, 9, 1, 0, 0, 0, 0, time.Local),
			expectedEnd:   time.Date(2023, 9, 30, 23, 59, 59, 0, time.Local),
		},
		{
			name:          "Start Of Year",
			monthOffset:   1,
			currentTime:   time.Date(2024, 1, 15, 10, 0, 0, 0, time.Local),
			expectedStart: time.Date(2023, 12, 1, 0, 0, 0, 0, time.Local),
			expectedEnd:   time.Date(2023, 12, 31, 23, 59, 59, 0, time.Local),
		},
		{
			name:          "End Of Year",
			monthOffset:   1,
			currentTime:   time.Date(2023, 12, 15, 10, 0, 0, 0, time.Local),
			expectedStart: time.Date(2023, 11, 1, 0, 0, 0, 0, time.Local),
			expectedEnd:   time.Date(2023, 11, 30, 23, 59, 59, 0, time.Local),
		},
		{
			name:          "Leap Year",
			monthOffset:   1,
			currentTime:   time.Date(2024, 3, 15, 10, 0, 0, 0, time.Local),
			expectedStart: time.Date(2024, 2, 1, 0, 0, 0, 0, time.Local),
			expectedEnd:   time.Date(2024, 2, 29, 23, 59, 59, 0, time.Local),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			dt := datetime.NewDatetime(func() time.Time {
				return tc.currentTime
			})
			start, end := dt.GetMonthStartAndEnd(tc.monthOffset)

			assert.Equal(t, tc.expectedStart, start)
			assert.Equal(t, tc.expectedEnd, end)
		})
	}
}

func TestDatetime_GetLastDayOfMonth(t *testing.T) {
	testCases := []struct {
		name            string
		inputTime       time.Time
		expectedLastDay int
	}{
		{
			name:            "Normal Month",
			inputTime:       time.Date(2023, 11, 15, 10, 0, 0, 0, time.Local),
			expectedLastDay: 30,
		},
		{
			name:            "February Non-leap Year",
			inputTime:       time.Date(2023, 2, 15, 10, 0, 0, 0, time.Local),
			expectedLastDay: 28,
		},
		{
			name:            "February Leap Year",
			inputTime:       time.Date(2024, 2, 15, 10, 0, 0, 0, time.Local),
			expectedLastDay: 29,
		},
		{
			name:            "January",
			inputTime:       time.Date(2023, 1, 15, 10, 0, 0, 0, time.Local),
			expectedLastDay: 31,
		},
		{
			name:            "December",
			inputTime:       time.Date(2023, 12, 15, 10, 0, 0, 0, time.Local),
			expectedLastDay: 31,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			dt := datetime.NewDatetime(time.Now)
			lastDay := dt.GetLastDayOfMonth(tc.inputTime)
			assert.Equal(t, tc.expectedLastDay, lastDay)
		})
	}
}
