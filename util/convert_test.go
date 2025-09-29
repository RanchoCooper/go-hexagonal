package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringToInt(t *testing.T) {
	testCases := []struct {
		name        string
		input       string
		expected    int
		expectError bool
	}{
		{"Valid positive integer", "123", 123, false},
		{"Valid negative integer", "-456", -456, false},
		{"Zero", "0", 0, false},
		{"Invalid string", "abc", 0, true},
		{"Empty string", "", 0, true},
		{"Float string", "123.45", 0, true},
		{"Large number", "2147483647", 2147483647, false},
		{"Very large number", "99999999999999999999", 0, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := StringToInt(tc.input)

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, result)
			}
		})
	}
}

func TestStringToInt64(t *testing.T) {
	testCases := []struct {
		name        string
		input       string
		expected    int64
		expectError bool
	}{
		{"Valid positive integer", "123", 123, false},
		{"Valid negative integer", "-456", -456, false},
		{"Zero", "0", 0, false},
		{"Invalid string", "abc", 0, true},
		{"Empty string", "", 0, true},
		{"Float string", "123.45", 0, true},
		{"Large number", "9223372036854775807", 9223372036854775807, false},
		{"Very large number", "99999999999999999999", 0, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := StringToInt64(tc.input)

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, result)
			}
		})
	}
}

func TestStringToUint(t *testing.T) {
	testCases := []struct {
		name        string
		input       string
		expected    uint
		expectError bool
	}{
		{"Valid positive integer", "123", 123, false},
		{"Zero", "0", 0, false},
		{"Invalid string", "abc", 0, true},
		{"Empty string", "", 0, true},
		{"Negative number", "-456", 0, true},
		{"Float string", "123.45", 0, true},
		{"Large number", "4294967295", 4294967295, false},
		{"Very large number", "99999999999999999999", 0, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := StringToUint(tc.input)

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, result)
			}
		})
	}
}

func TestStringToUint64(t *testing.T) {
	testCases := []struct {
		name        string
		input       string
		expected    uint64
		expectError bool
	}{
		{"Valid positive integer", "123", 123, false},
		{"Zero", "0", 0, false},
		{"Invalid string", "abc", 0, true},
		{"Empty string", "", 0, true},
		{"Negative number", "-456", 0, true},
		{"Float string", "123.45", 0, true},
		{"Large number", "18446744073709551615", 18446744073709551615, false},
		{"Very large number", "99999999999999999999", 0, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := StringToUint64(tc.input)

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, result)
			}
		})
	}
}

func TestStringToFloat64(t *testing.T) {
	testCases := []struct {
		name        string
		input       string
		expected    float64
		expectError bool
	}{
		{"Valid positive float", "123.45", 123.45, false},
		{"Valid negative float", "-456.78", -456.78, false},
		{"Integer string", "123", 123.0, false},
		{"Zero", "0", 0.0, false},
		{"Scientific notation", "1.23e4", 12300.0, false},
		{"Invalid string", "abc", 0, true},
		{"Empty string", "", 0, true},
		{"Multiple decimal points", "123.45.67", 0, true},
		{"Very large number", "1.7976931348623157e+308", 1.7976931348623157e+308, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := StringToFloat64(tc.input)

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.InDelta(t, tc.expected, result, 0.0001)
			}
		})
	}
}

func TestStringToBool(t *testing.T) {
	testCases := []struct {
		name        string
		input       string
		expected    bool
		expectError bool
	}{
		{"True (lowercase)", "true", true, false},
		{"True (uppercase)", "TRUE", true, false},
		{"True (mixed case)", "True", true, false},
		{"False (lowercase)", "false", false, false},
		{"False (uppercase)", "FALSE", false, false},
		{"False (mixed case)", "False", false, false},
		{"1", "1", true, false},
		{"0", "0", false, false},
		{"Invalid string", "abc", false, true},
		{"Empty string", "", false, true},
		{"Yes", "yes", false, true},
		{"No", "no", false, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := StringToBool(tc.input)

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, result)
			}
		})
	}
}

// Benchmark tests for performance comparison
func BenchmarkStringToInt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = StringToInt("12345")
	}
}

func BenchmarkStringToInt64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = StringToInt64("12345")
	}
}

func BenchmarkStringToUint(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = StringToUint("12345")
	}
}

func BenchmarkStringToUint64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = StringToUint64("12345")
	}
}

func BenchmarkStringToFloat64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = StringToFloat64("123.45")
	}
}

func BenchmarkStringToBool(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = StringToBool("true")
	}
}

// Test edge cases and boundary conditions
func TestStringToInt_EdgeCases(t *testing.T) {
	// Test with leading/trailing whitespace (should fail)
	result, err := StringToInt(" 123 ")
	assert.Error(t, err)
	assert.Equal(t, 0, result)

	// Test with plus sign (Go's strconv.Atoi can handle plus signs)
	result, err = StringToInt("+123")
	assert.NoError(t, err)
	assert.Equal(t, 123, result)
}

func TestStringToFloat64_EdgeCases(t *testing.T) {
	// Test with leading/trailing whitespace (should fail)
	result, err := StringToFloat64(" 123.45 ")
	assert.Error(t, err)
	assert.Equal(t, 0.0, result)

	// Test with plus sign
	result, err = StringToFloat64("+123.45")
	assert.NoError(t, err)
	assert.Equal(t, 123.45, result)

	// Test infinity (Go's strconv.ParseFloat can parse these special values)
	result, err = StringToFloat64("Inf")
	assert.NoError(t, err)
	assert.True(t, result > 1e308) // Positive infinity

	result, err = StringToFloat64("-Inf")
	assert.NoError(t, err)
	assert.True(t, result < -1e308) // Negative infinity

	result, err = StringToFloat64("NaN")
	assert.NoError(t, err)
	assert.True(t, result != result) // NaN is not equal to itself
}

// Test consistency across different numeric types
func TestNumericConversionConsistency(t *testing.T) {
	// Test that "123" converts consistently across different types
	intVal, err := StringToInt("123")
	assert.NoError(t, err)
	assert.Equal(t, 123, intVal)

	int64Val, err := StringToInt64("123")
	assert.NoError(t, err)
	assert.Equal(t, int64(123), int64Val)

	uintVal, err := StringToUint("123")
	assert.NoError(t, err)
	assert.Equal(t, uint(123), uintVal)

	uint64Val, err := StringToUint64("123")
	assert.NoError(t, err)
	assert.Equal(t, uint64(123), uint64Val)

	floatVal, err := StringToFloat64("123")
	assert.NoError(t, err)
	assert.Equal(t, 123.0, floatVal)

	boolVal, err := StringToBool("1")
	assert.NoError(t, err)
	assert.True(t, boolVal)
}

// Test error messages for better debugging
func TestErrorMessages(t *testing.T) {
	// Test that error messages are informative
	_, err := StringToInt("abc")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid syntax")

	_, err = StringToFloat64("xyz")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid syntax")

	_, err = StringToBool("maybe")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid syntax")
}
