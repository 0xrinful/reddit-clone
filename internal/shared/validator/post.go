package validator

const (
	postTitleMax = 120
	postBodyMax  = 40000
)

func ValidatePostTitle(v *Validator, title string) {
	v.Check(NotBlank(title), "title", "must not be blank")
	v.Check(MaxLength(title, postTitleMax), "title", "must not exceed 120 characters")
}

func ValidatePostBody(v *Validator, body string) {
	v.Check(NotBlank(body), "body", "must not be blank")
	v.Check(MaxLength(body, postBodyMax), "body", "must not exceed 40000 characters")
}
