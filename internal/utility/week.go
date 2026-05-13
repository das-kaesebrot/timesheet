package utility

import (
	"errors"
	"math"
	"time"
)

// generates a range of week start dates lying between start (inclusive) and end (exclusive) aligned to the given weekStartDay
// if a week doesnt start in the range, returns an empty list
func GetWeekRange(startInclusive time.Time, endExclusive time.Time, weekStartDay time.Weekday) ([]time.Time, error) {
	outputTimezone := startInclusive.Location()
	startInclusive = AsUTC(ZeroTimeComponents(startInclusive))
	endExclusive = AsUTC(ZeroTimeComponents(endExclusive))

	if endExclusive.Equal(startInclusive) || endExclusive.Before(startInclusive) {
		return nil, errors.New("end date can't be equal or before start date!")
	}

	containsWeekDay, err := ContainsWeekDay(startInclusive, endExclusive, weekStartDay)
	if err != nil {
		return nil, err
	}
	if !containsWeekDay {
		return []time.Time{}, nil
	}

	diffDays := endExclusive.Sub(startInclusive).Hours() / (24 * 7)
	weekStartsToGenerate := int(math.Ceil(diffDays))

	weekRange := make([]time.Time, weekStartsToGenerate)

	for i := range weekStartsToGenerate {
		weekRange[i] = AsTimezone(GetNextWeekStartDate(startInclusive.AddDate(0, 0, i*7), weekStartDay), outputTimezone)
	}

	return weekRange, nil
}

// Returns true if given range of dates contains the given weekday at least once.
func ContainsWeekDay(startInclusive time.Time, endExclusive time.Time, weekDay time.Weekday) (bool, error) {
	// to work around daylight savings time bugs (where a day might only be 23 hours long),
	// we always "cast" time objects to UTC when doing any calculations on them
	startInclusive = AsUTC(ZeroTimeComponents(startInclusive))
	endExclusive = AsUTC(ZeroTimeComponents(endExclusive))

	if endExclusive.Equal(startInclusive) || endExclusive.Before(startInclusive) {
		return false, errors.New("end date can't be equal or before start date!")
	}

	diff := endExclusive.Sub(startInclusive)

	// if at least a week is between the given dates, the weekday has been seen
	if diff.Hours() >= (7 * 24) {
		return true, nil
	}

	startDay := startInclusive.Weekday()
	endDay := endExclusive.Weekday()
	if endDay < startDay {
		endDay += 7
	}
	if weekDay < startDay {
		weekDay += 7
	}

	return (startDay <= weekDay && endDay > weekDay), nil
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

// Check whether a given time object lies within a given week, starting at start and ending at start + 7 days.
// start is always zeroed using ZeroTimeComponents.
// Might not cover edge cases such as daylight savings time changes.
func IsInWeekFromStartDate(start time.Time, timeToCheck time.Time) bool {
	start = ZeroTimeComponents(start)
	end := start.AddDate(0, 0, 7)

	return (timeToCheck.Equal(start) || timeToCheck.After(start)) && timeToCheck.Before(end)
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
