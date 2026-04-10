package validator

const (
	postTitleMin = 3
	postTitleMax = 120
	postBodyMin  = 10
	postBodyMax  = 40000
)

func ValidatePostTitle(v *Validator, title string) {
	v.Check(NotBlank(title), "title", "must not be blank")
	v.Check(MinLength(title, postTitleMin), "title", "must be at least 3 characters")
	v.Check(MaxLength(title, postTitleMax), "title", "must not exceed 120 characters")
}

func ValidatePostBody(v *Validator, body string) {
	v.Check(NotBlank(body), "body", "must not be blank")
	v.Check(MinLength(body, postBodyMin), "body", "must be at least 10 characters")
	v.Check(MaxLength(body, postBodyMax), "body", "must not exceed 40000 characters")
}
