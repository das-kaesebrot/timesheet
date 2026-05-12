package utility

import "time"

func GenerateWeekRanges(start time.Time, end time.Time, weekStartDay time.Weekday) ([]time.Time, error) {
	start = time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, start.Location())

	return []time.Time{}, nil
}

// going from the given start time (or rather, the corresponding date), find and return the start of the next week, aligned to the given week start day.
// if weekStartDay is the same day, the same day is returned
func GetNextWeekStartDate(start time.Time, weekStartDay time.Weekday) time.Time {
	start = time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, start.Location())

	delta := int((weekStartDay - start.Weekday() + 7) % 7)
	start = start.AddDate(0, 0, delta)

	return start
}

func GetPreviousWeekStartDate(start time.Time, weekStartDay time.Weekday) time.Time {
	start = time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, start.Location())

	delta := int((start.Weekday() - weekStartDay + 7) % 7)
	start = start.AddDate(0, 0, -delta)

	return start
}
