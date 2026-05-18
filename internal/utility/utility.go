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
