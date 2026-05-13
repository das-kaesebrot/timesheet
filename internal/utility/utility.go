package utility

func assert[T any](i T, err error) T {
	if err != nil {
		panic(err)
	}
	return i
}
