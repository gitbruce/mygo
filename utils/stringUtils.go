package utils

import "strconv"

func IsInteger(val float64) bool {
	return val == float64(int(val))
}

func PrefixString(val float64) string {
	return strconv.FormatFloat(val, 'f', -1, 32)
}
