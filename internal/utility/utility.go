package utility

import (
	"strconv"
)

func ParseUint8(s string) (uint8, error) {
	parsedUint, err := strconv.ParseUint(s, 10, 8)

	if err != nil {
		return 0, err
	}

	return uint8(parsedUint), nil
}
