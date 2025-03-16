package error_code

import (
	"strconv"
)

// StringToInt converts a string to an integer
func StringToInt(str string) (int, error) {
	return strconv.Atoi(str)
}

// StringToInt64 converts a string to an int64
func StringToInt64(str string) (int64, error) {
	return strconv.ParseInt(str, 10, 64)
}

// StringToUint converts a string to an unsigned integer
func StringToUint(str string) (uint, error) {
	val, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(val), nil
}

// StringToUint64 converts a string to an uint64
func StringToUint64(str string) (uint64, error) {
	return strconv.ParseUint(str, 10, 64)
}

// StringToFloat64 converts a string to a float64
func StringToFloat64(str string) (float64, error) {
	return strconv.ParseFloat(str, 64)
}

// StringToBool converts a string to a boolean
func StringToBool(str string) (bool, error) {
	return strconv.ParseBool(str)
}
