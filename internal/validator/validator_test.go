package validator

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateEmail(t *testing.T) {
	assertions := assert.New(t)

	// Test cases
	testCases := []struct {
		email  string
		expect error
	}{
		{"test@example.com", nil},
		{"invalid-email", errors.New("invalid email format")},
		{"user@domain.com", nil},
		{"user@domain", errors.New("invalid email format")},
	}

	for _, testCase := range testCases {
		err := Email(testCase.email)
		if testCase.expect == nil {
			assertions.NoError(err, "Expected no error for email: %s", testCase.email)
		} else {
			assertions.EqualError(err, testCase.expect.Error(), "Expected error for email: %s", testCase.email)
		}
	}
}

func TestValidatePassword(t *testing.T) {
	assertions := assert.New(t)

	// Test cases
	testCases := []struct {
		password string
		expect   error
	}{
		{"Password123!", nil},
		{"short", errors.New("password must be at least 8 characters long")},
		{"lowercaseonly123", errors.New("password must contain at least one uppercase letter")},
		{"UPPERCASEONLY123", errors.New("password must contain at least one lowercase letter")},
		{"NoDigits!", errors.New("password must contain at least one digit")},
		{"NoSpecial123", errors.New("password must contain at least one special character")},
	}

	for _, testCase := range testCases {
		err := Password(testCase.password)
		if testCase.expect == nil {
			assertions.NoError(err, "Expected no error for password: %s", testCase.password)
		} else {
			assertions.EqualError(err, testCase.expect.Error(), "Expected error for password: %s", testCase.password)
		}
	}
}
