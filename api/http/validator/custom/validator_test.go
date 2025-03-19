package custom

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func TestValidatePhone(t *testing.T) {
	validate := validator.New()
	err := validate.RegisterValidation("phone", validatePhone)
	assert.NoError(t, err)

	type TestStruct struct {
		Phone string `validate:"phone"`
	}

	tests := []struct {
		name    string
		phone   string
		wantErr bool
	}{
		{
			name:    "Valid phone without country code",
			phone:   "1234567890",
			wantErr: false,
		},
		{
			name:    "Valid phone with country code",
			phone:   "+861234567890",
			wantErr: false,
		},
		{
			name:    "Valid phone with country code and hyphen",
			phone:   "+86-1234567890",
			wantErr: false,
		},
		{
			name:    "Invalid phone - too short",
			phone:   "123456",
			wantErr: true,
		},
		{
			name:    "Invalid phone - contains letters",
			phone:   "123abc4567",
			wantErr: true,
		},
		{
			name:    "Invalid phone - wrong format",
			phone:   "++1234567890",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			test := TestStruct{
				Phone: tt.phone,
			}
			err := validate.Struct(test)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateUsername(t *testing.T) {
	validate := validator.New()
	_ = validate.RegisterValidation("username", validateUsername)

	type TestStruct struct {
		Username string `validate:"username"`
	}

	tests := []struct {
		name     string
		username string
		wantErr  bool
	}{
		{"Valid username", "john_doe123", false},
		{"Valid username - minimum length", "abc", false},
		{"Invalid username - too short", "ab", true},
		{"Invalid username - too long", "abcdefghijklmnopqrstu", true},
		{"Invalid username - special chars", "john@doe", true},
		{"Invalid username - empty", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			test := TestStruct{Username: tt.username}
			err := validate.Struct(test)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidatePassword(t *testing.T) {
	validate := validator.New()
	_ = validate.RegisterValidation("password", validatePassword)

	type TestStruct struct {
		Password string `validate:"password"`
	}

	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{"Valid password", "Test123!", false},
		{"Valid password - complex", "P@ssw0rd!", false},
		{"Invalid password - no uppercase", "test123!", true},
		{"Invalid password - no lowercase", "TEST123!", true},
		{"Invalid password - no number", "TestTest!", true},
		{"Invalid password - no special char", "Test1234", true},
		{"Invalid password - too short", "Te1!", true},
		{"Invalid password - empty", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			test := TestStruct{Password: tt.password}
			err := validate.Struct(test)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
