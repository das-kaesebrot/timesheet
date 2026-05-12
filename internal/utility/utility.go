package utility

func Assert[T any](i T, err error) T {
	if err != nil {
		panic(err)
	}
	return i
}
