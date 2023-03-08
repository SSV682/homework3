package stock

import (
	"strconv"
)

func queryParamsToUInt64(value string, baseValue uint64) (uint64, error) {
	if value == "" {
		return baseValue, nil
	}

	limit, err := strconv.Atoi(value)
	if err != nil {
		return 0, err
	}

	return uint64(limit), nil
}
