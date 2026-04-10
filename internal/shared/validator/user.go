package validator

import "fmt"

const (
	userNameMin = 3
	userNameMax = 32
	passwordMin = 8
	passwordMax = 72
)

func ValidateEmail(v *Validator, email string) {
	v.Check(NotBlank(email), "email", "must not be blank")
	v.Check(Matches(email, EmailRX), "email", "must be a valid email address")
}

func ValidateUsername(v *Validator, username string) {
	v.Check(NotBlank(username), "username", "must not be blank")
	v.Check(
		MinLength(username, userNameMin),
		"username",
		fmt.Sprintf("must be at least %d characters", userNameMin),
	)
	v.Check(
		MaxLength(username, userNameMax),
		"username",
		fmt.Sprintf("must not exceed %d characters", userNameMax),
	)
}

func ValidatePassword(v *Validator, password string) {
	v.Check(NotBlank(password), "password", "must not be blank")
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
