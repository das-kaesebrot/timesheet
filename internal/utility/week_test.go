package utility

import (
	"reflect"
	"testing"
	"time"
)

var loc *time.Location = assert(time.LoadLocation("Europe/Berlin"))
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

var containsWeekDayTests = []struct {
	inStart      time.Time
	inEnd        time.Time
	weekStartDay time.Weekday
	out          bool
}{
	{time.Date(2026, 05, 11, 22, 54, 8, 0, loc), time.Date(2026, 05, 12, 22, 54, 8, 0, loc), time.Monday, true},
	{time.Date(2026, 05, 12, 22, 54, 8, 0, loc), time.Date(2026, 05, 16, 22, 54, 8, 0, loc), time.Monday, false},
	{time.Date(2026, 05, 12, 22, 54, 8, 0, loc), time.Date(2026, 05, 16, 22, 54, 8, 0, loc), time.Wednesday, true},
	{time.Date(2026, 05, 12, 22, 54, 8, 0, loc), time.Date(2026, 05, 19, 22, 54, 8, 0, loc), time.Monday, true},
	{time.Date(2026, 05, 12, 22, 54, 8, 0, loc), time.Date(2026, 05, 30, 22, 54, 8, 0, loc), time.Monday, true},
	{time.Date(2026, 05, 12, 22, 54, 8, 0, loc), time.Date(2026, 05, 13, 22, 54, 8, 0, loc), time.Friday, false},
	{time.Date(2026, 03, 28, 0, 0, 0, 0, loc), time.Date(2026, 03, 30, 0, 0, 0, 0, loc), time.Sunday, true}, // DST start
	{time.Date(2026, 10, 24, 0, 0, 0, 0, loc), time.Date(2026, 10, 26, 0, 0, 0, 0, loc), time.Sunday, true}, // DST end
	{time.Date(2026, 10, 20, 0, 0, 0, 0, loc), time.Date(2026, 10, 26, 0, 0, 0, 0, loc), time.Sunday, true}, // DST end
}

func TestContainsWeekDay(t *testing.T) {
	for _, tt := range containsWeekDayTests {
		result, err := ContainsWeekDay(tt.inStart, tt.inEnd, tt.weekStartDay)
		if result != tt.out || err != nil {
			t.Errorf("result is not expected value! start=%v, end=%v, weekDay=%v, result=%v, expected=%v, err=%v", tt.inStart, tt.inEnd, tt.weekStartDay, result, tt.out, err)
		}
	}
}

var asUTCtests = []struct {
	in  time.Time
	out time.Time
}{
	{time.Date(2026, 05, 12, 22, 54, 8, 0, loc), time.Date(2026, 05, 12, 22, 54, 8, 0, time.UTC)},
	{time.Date(2026, 03, 29, 03, 0, 0, 0, loc), time.Date(2026, 03, 29, 03, 0, 0, 0, time.UTC)},
}

func TestAsUTC(t *testing.T) {
	for _, tt := range asUTCtests {
		result := AsUTC(tt.in)
		if result != tt.out {
			t.Errorf("result is not expected value! result=%v, expected=%v", result, tt.out)
		}
	}
}

var weekRangeTests = []struct {
	inStart  time.Time
	inEnd    time.Time
	weekDay  time.Weekday
	expected []time.Time
}{
	{time.Date(2026, 05, 11, 22, 54, 8, 0, loc), time.Date(2026, 05, 12, 22, 54, 8, 0, loc), time.Monday, []time.Time{time.Date(2026, 05, 11, 0, 0, 0, 0, loc)}},
	{time.Date(2026, 05, 12, 22, 54, 8, 0, loc), time.Date(2026, 05, 19, 22, 54, 8, 0, loc), time.Monday, []time.Time{time.Date(2026, 05, 18, 0, 0, 0, 0, loc)}},
	{time.Date(2026, 05, 11, 22, 54, 8, 0, loc), time.Date(2026, 06, 30, 22, 54, 8, 0, loc), time.Monday,
		[]time.Time{
			time.Date(2026, 05, 11, 0, 0, 0, 0, loc),
			time.Date(2026, 05, 18, 0, 0, 0, 0, loc),
			time.Date(2026, 05, 25, 0, 0, 0, 0, loc),
			time.Date(2026, 06, 01, 0, 0, 0, 0, loc),
			time.Date(2026, 06, 8, 0, 0, 0, 0, loc),
			time.Date(2026, 06, 15, 0, 0, 0, 0, loc),
			time.Date(2026, 06, 22, 0, 0, 0, 0, loc),
			time.Date(2026, 06, 29, 0, 0, 0, 0, loc),
		},
	},
	{time.Date(2026, 05, 11, 22, 54, 8, 0, loc), time.Date(2026, 05, 16, 22, 54, 8, 0, loc), time.Sunday, []time.Time{}},
}

func TestGetWeekRange(t *testing.T) {
	for _, tt := range weekRangeTests {
		result, err := GetWeekRange(tt.inStart, tt.inEnd, tt.weekDay)
		if !reflect.DeepEqual(result, tt.expected) || err != nil {
			t.Errorf("result is not expected value!\nstart=%v\nend=%v\nexpected=%v\nresult=%v\nerr=%v", tt.inStart, tt.inEnd, tt.expected, result, err)
		}
	}
}
