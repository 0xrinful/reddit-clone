package validator

import "fmt"

const (
	usernameMin = 3
	usernameMax = 32
	passwordMin = 8
	passwordMax = 72
)

func ValidateEmail(v *Validator, email string) {
	v.Check(NotBlank(email), "email", "must not be empty")
	v.Check(Matches(email, EmailRX), "email", "must be a valid email address")
}

func ValidateUsername(v *Validator, username string) {
	v.Check(
		MinLength(username, usernameMin),
		"username",
		fmt.Sprintf("must be at least %d characters", usernameMin),
	)
	v.Check(
		MaxLength(username, usernameMax),
		"username",
		fmt.Sprintf("must not exceed %d characters", usernameMax),
	)
	v.Check(
		Matches(username, UsernameRX),
		"username",
		"must contain only letters, numbers, and underscores, and must start with a letter or number",
	)
}

func ValidatePassword(v *Validator, password string) {
	v.Check(NotBlank(password), "password", "must not be empty")
	v.Check(
		MinLength(password, passwordMin),
		"password",
		fmt.Sprintf("must be at least %d characters", passwordMin),
	)
	v.Check(
		MaxLength(password, passwordMax),
		"password",
		fmt.Sprintf("must not exceed %d characters", passwordMax),
	)
}
