package utility

import (
	"fmt"
	"math"
	"time"
)

func assert[T any](i T, err error) T {
	if err != nil {
		panic(err)
	}
	return i
}

func GetWeekdays() []time.Weekday {
	return []time.Weekday{
		time.Sunday,
		time.Monday,
		time.Tuesday,
		time.Wednesday,
		time.Thursday,
		time.Friday,
		time.Saturday,
	}
}

// Source - https://stackoverflow.com/a/74940188
// Posted by user6705546, modified by community. See post 'Timeline' for change history
// Retrieved 2026-08-31, License - CC BY-SA 4.0
func Divmod(numerator, denominator int64) (quotient, remainder int64) {
	quotient = numerator / denominator
	remainder = numerator % denominator
	if remainder != 0 && numerator*denominator < 0 {
		remainder += denominator
		quotient--
	}
	return
}

func GetFormattedDuration(d time.Duration, compressed bool) string {
	result := ""

	totalMinutes := math.Round(d.Minutes())
	prefix := ""
	absTotalMinutes := math.Abs(totalMinutes)
	if absTotalMinutes != totalMinutes {
		prefix = "-"
		totalMinutes = absTotalMinutes
	}

	result = prefix
	hours, minutes := Divmod(int64(totalMinutes), 60)
	if !compressed || hours != 0 {
		result += fmt.Sprintf("%02dh", hours)
	}

	if !compressed || minutes != 0 {
		result += fmt.Sprintf("%02dm", minutes)
	}

	return result
}
