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
}

func TestPreviousWeekStartGetter(t *testing.T) {
	for _, tt := range previousWeekStartGetterTests {
		result := GetPreviousWeekStartDate(tt.in, tt.weekStartDay)
		if result != tt.out {
			t.Errorf("result is not expected value! result=%v, expected=%v", result, tt.out)
		}
	}
}
