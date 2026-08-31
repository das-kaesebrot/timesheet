package utility

import "time"

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
