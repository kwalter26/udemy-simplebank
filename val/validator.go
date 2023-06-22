package val

import (
	"fmt"
	"net/mail"
	"regexp"
)

var (
	isValidUsername = regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString
	isValidFullName = regexp.MustCompile(`^[a-zA-Z\s]+$`).MatchString
)

// ValidateString function to Validate a string and make sure it falls between a  min and max length. Returns an error
// if the string is invalid.
func ValidateString(str string, min int, max int) error {
	if len(str) < min || len(str) > max {
		return fmt.Errorf("invalid string length: must be between %d and %d characters", min, max)
	}
	return nil
}

// ValidateUsername function to Validate a username and make sure it falls between a  min and max length. Also,
// has to have a specific character set Returns an error if the username is invalid.
func ValidateUsername(username string) error {
	if err := ValidateString(username, 3, 100); err != nil {
		return err
	}
	if !isValidUsername(username) {
		return fmt.Errorf("username can only contain letters, numbers, and underscores")
	}
	return nil
}

// ValidateFullName function to Validate a full name and make sure it falls between a  min and max length.
// Returns an error if the full name is invalid.
func ValidateFullName(fullName string) error {
	if err := ValidateString(fullName, 3, 100); err != nil {
		return err
	}
	if !isValidFullName(fullName) {
		return fmt.Errorf("fullname can only contain letters and spaces")
	}
	return nil
}

// ValidatePassword function to Validate a password and make sure it falls between a  min and max length.
// Returns an error if the password is invalid.
func ValidatePassword(password string) error {
	if err := ValidateString(password, 6, 100); err != nil {
		return err
	}
	return nil
}

// ValidateEmail function to Validate an email and make sure it falls between a  min and max length. Must also be a valid email.
// Returns an error if the email is invalid.
func ValidateEmail(email string) error {
	if err := ValidateString(email, 6, 100); err != nil {
		return err
	}
	if _, err := mail.ParseAddress(email); err != nil {
		return fmt.Errorf("invalid email address")
	}
	return nil
}

// ValidateEmailId function to Validate an email id and make sure it is greater than 0.
func ValidateEmailId(value int64) error {
	if value < 0 {
		return fmt.Errorf("invalid email id")
	}
	return nil
}

// ValidateSecretCode function to Validate a secret code and make sure it falls between a  min and max length.
func ValidateSecretCode(value string) error {
	if err := ValidateString(value, 32, 128); err != nil {
		return err
	}
	return nil
}
