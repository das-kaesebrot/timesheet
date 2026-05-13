package utility

import "time"

func GenerateWeekRanges(start time.Time, end time.Time, weekStartDay time.Weekday) ([]time.Time, error) {
	start = time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, start.Location())

	return []time.Time{}, nil
}

// going from the given start time (or rather, the corresponding date), find and return the start of the next week, aligned to the given week start day.
// if weekStartDay is the same day, the same day is returned
func GetNextWeekStartDate(start time.Time, weekStartDay time.Weekday) time.Time {
	start = ZeroTimeComponents(start)

	delta := int((weekStartDay - start.Weekday() + 7) % 7)
	start = start.AddDate(0, 0, delta)

	return start
}

func GetPreviousWeekStartDate(start time.Time, weekStartDay time.Weekday) time.Time {
	start = ZeroTimeComponents(start)

	delta := int((start.Weekday() - weekStartDay + 7) % 7)
	start = start.AddDate(0, 0, -delta)

	return start
}

// returns a new date from the given date that has the hours, minutes, seconds and nanoseconds components zeroed.
func ZeroTimeComponents(toZero time.Time) time.Time {
	return time.Date(toZero.Year(), toZero.Month(), toZero.Day(), 0, 0, 0, 0, toZero.Location())
}

// returns a copy of the time with the exact same values, but with the timezone set to UTC
func AsUTC(toAdjust time.Time) time.Time {
	return AsTimezone(toAdjust, time.UTC)
}

// returns a copy of the time with the exact same values, but with the timezone set to loc
func AsTimezone(toAdjust time.Time, loc *time.Location) time.Time {
	return time.Date(toAdjust.Year(), toAdjust.Month(), toAdjust.Day(), toAdjust.Hour(), toAdjust.Minute(), toAdjust.Second(), toAdjust.Nanosecond(), loc)
}
