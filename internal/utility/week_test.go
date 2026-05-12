package utility

import (
	"testing"
	"time"
)

var loc *time.Location = Assert(time.LoadLocation("Europe/Berlin"))
var nextWeekStartGetterTests = []struct {
	in           time.Time
	weekStartDay time.Weekday
	out          time.Time
}{
	{time.Date(2026, 05, 11, 22, 54, 8, 0, loc), time.Monday, time.Date(2026, 05, 11, 0, 0, 0, 0, loc)},
	{time.Date(2026, 05, 12, 22, 54, 8, 0, loc), time.Monday, time.Date(2026, 05, 18, 0, 0, 0, 0, loc)},
	{time.Date(2026, 05, 13, 22, 54, 8, 0, loc), time.Monday, time.Date(2026, 05, 18, 0, 0, 0, 0, loc)},
	{time.Date(2026, 05, 14, 22, 54, 8, 0, loc), time.Monday, time.Date(2026, 05, 18, 0, 0, 0, 0, loc)},
	{time.Date(2026, 05, 15, 22, 54, 8, 0, loc), time.Monday, time.Date(2026, 05, 18, 0, 0, 0, 0, loc)},
	{time.Date(2026, 05, 16, 22, 54, 8, 0, loc), time.Monday, time.Date(2026, 05, 18, 0, 0, 0, 0, loc)},
	{time.Date(2026, 05, 17, 22, 54, 8, 0, loc), time.Monday, time.Date(2026, 05, 18, 0, 0, 0, 0, loc)},
	{time.Date(2026, 05, 18, 22, 54, 8, 0, loc), time.Monday, time.Date(2026, 05, 18, 0, 0, 0, 0, loc)},
	{time.Date(2026, 05, 11, 22, 54, 8, 0, loc), time.Monday, time.Date(2026, 05, 11, 0, 0, 0, 0, loc)},
	{time.Date(2026, 05, 11, 22, 54, 8, 0, loc), time.Tuesday, time.Date(2026, 05, 12, 0, 0, 0, 0, loc)},
	{time.Date(2026, 05, 11, 22, 54, 8, 0, loc), time.Wednesday, time.Date(2026, 05, 13, 0, 0, 0, 0, loc)},
	{time.Date(2026, 05, 11, 22, 54, 8, 0, loc), time.Thursday, time.Date(2026, 05, 14, 0, 0, 0, 0, loc)},
	{time.Date(2026, 05, 11, 22, 54, 8, 0, loc), time.Friday, time.Date(2026, 05, 15, 0, 0, 0, 0, loc)},
	{time.Date(2026, 05, 11, 22, 54, 8, 0, loc), time.Saturday, time.Date(2026, 05, 16, 0, 0, 0, 0, loc)},
	{time.Date(2026, 05, 11, 22, 54, 8, 0, loc), time.Sunday, time.Date(2026, 05, 17, 0, 0, 0, 0, loc)},
}

func TestNextWeekStartGetter(t *testing.T) {
	for _, tt := range nextWeekStartGetterTests {
		result := GetNextWeekStartDate(tt.in, tt.weekStartDay)
		if result != tt.out {
			t.Errorf("result is not expected value! result=%v, expected=%v", result, tt.out)
		}
	}
}

var previousWeekStartGetterTests = []struct {
	in           time.Time
	weekStartDay time.Weekday
	out          time.Time
}{
	{time.Date(2026, 05, 11, 22, 54, 8, 0, loc), time.Monday, time.Date(2026, 05, 11, 0, 0, 0, 0, loc)},
	{time.Date(2026, 05, 12, 22, 54, 8, 0, loc), time.Monday, time.Date(2026, 05, 11, 0, 0, 0, 0, loc)},
	{time.Date(2026, 05, 13, 22, 54, 8, 0, loc), time.Monday, time.Date(2026, 05, 11, 0, 0, 0, 0, loc)},
	{time.Date(2026, 05, 14, 22, 54, 8, 0, loc), time.Monday, time.Date(2026, 05, 11, 0, 0, 0, 0, loc)},
	{time.Date(2026, 05, 15, 22, 54, 8, 0, loc), time.Monday, time.Date(2026, 05, 11, 0, 0, 0, 0, loc)},
	{time.Date(2026, 05, 16, 22, 54, 8, 0, loc), time.Monday, time.Date(2026, 05, 11, 0, 0, 0, 0, loc)},
	{time.Date(2026, 05, 17, 22, 54, 8, 0, loc), time.Monday, time.Date(2026, 05, 11, 0, 0, 0, 0, loc)},
	{time.Date(2026, 05, 18, 22, 54, 8, 0, loc), time.Monday, time.Date(2026, 05, 18, 0, 0, 0, 0, loc)},
	{time.Date(2026, 05, 11, 22, 54, 8, 0, loc), time.Monday, time.Date(2026, 05, 11, 0, 0, 0, 0, loc)},
	{time.Date(2026, 05, 11, 22, 54, 8, 0, loc), time.Tuesday, time.Date(2026, 05, 05, 0, 0, 0, 0, loc)},
	{time.Date(2026, 05, 11, 22, 54, 8, 0, loc), time.Wednesday, time.Date(2026, 05, 06, 0, 0, 0, 0, loc)},
	{time.Date(2026, 05, 11, 22, 54, 8, 0, loc), time.Thursday, time.Date(2026, 05, 07, 0, 0, 0, 0, loc)},
	{time.Date(2026, 05, 11, 22, 54, 8, 0, loc), time.Friday, time.Date(2026, 05, 8, 0, 0, 0, 0, loc)},
	{time.Date(2026, 05, 11, 22, 54, 8, 0, loc), time.Saturday, time.Date(2026, 05, 9, 0, 0, 0, 0, loc)},
	{time.Date(2026, 05, 11, 22, 54, 8, 0, loc), time.Sunday, time.Date(2026, 05, 10, 0, 0, 0, 0, loc)},
}

func TestPreviousWeekStartGetter(t *testing.T) {
	for _, tt := range previousWeekStartGetterTests {
		result := GetPreviousWeekStartDate(tt.in, tt.weekStartDay)
		if result != tt.out {
			t.Errorf("result is not expected value! result=%v, expected=%v", result, tt.out)
		}
	}
}
