package validator

import "fmt"

const (
	postTitleMax = 120
	postBodyMax  = 40000
)

func ValidatePostTitle(v *Validator, title string) {
	v.Check(NotBlank(title), "title", "must not be empty")
	v.Check(MaxLength(title, postTitleMax), "title",
		fmt.Sprintf("must not exceed %d characters", postTitleMax))
}

func ValidatePostBody(v *Validator, body string) {
	v.Check(NotBlank(body), "body", "must not be empty")
	v.Check(MaxLength(body, postBodyMax), "body",
		fmt.Sprintf("must not exceed %d characters", postBodyMax))
}
