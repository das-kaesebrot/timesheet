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

// Source - https://stackoverflow.com/a/43945812
// Posted by Kaedys
// Retrieved 2026-07-07, License - CC BY-SA 3.0

func Divmod(numerator, denominator int64) (quotient, remainder int64) {
	quotient = numerator / denominator // integer division, decimals are truncated
	remainder = numerator % denominator
	// fix remainder so that it always returns a positive value
	// https://www.reddit.com/r/golang/comments/bnvik4/modulo_in_golang/
	if remainder < 0 {
		remainder += denominator
	}
	return
}
