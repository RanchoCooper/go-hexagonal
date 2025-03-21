package custom

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

// Min password length constant
const (
	// MinPasswordLength defines the minimum length for password validation
	MinPasswordLength = 8
)

// RegisterValidators registers custom validators
func RegisterValidators(v *validator.Validate) {
	// Register phone number validator
	_ = v.RegisterValidation("phone", validatePhone)

	// Register username validator
	_ = v.RegisterValidation("username", validateUsername)

	// Register password validator
	_ = v.RegisterValidation("password", validatePassword)
}

// validatePhone validates phone numbers
func validatePhone(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	// Phone number validation
	// Matches: +1234567890, 1234567890, +86-1234567890, etc.
	match, _ := regexp.MatchString(`^(\+\d{1,3}[-]?)?\d{10}$`, phone)
	return match
}

// validateUsername validates usernames
func validateUsername(fl validator.FieldLevel) bool {
	username := fl.Field().String()
	// Username must be 3-20 characters long and contain only letters, numbers, and underscores
	match, _ := regexp.MatchString(`^[a-zA-Z0-9_]{3,20}$`, username)
	return match
}

// validatePassword validates passwords
func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	// Password must be at least 8 characters long and contain:
	// - At least one uppercase letter
	// - At least one lowercase letter
	// - At least one number
	// - At least one special character
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	hasSpecial := regexp.MustCompile(`[!@#$%^&*]`).MatchString(password)
	hasLength := len(password) >= MinPasswordLength
	return hasUpper && hasLower && hasNumber && hasSpecial && hasLength
}
