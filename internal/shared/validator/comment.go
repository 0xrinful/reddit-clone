package validator

import "fmt"

const (
	commentContentMax = 40000
)

func ValidateCommentContent(v *Validator, content string) {
	v.Check(NotBlank(content), "content", "must not be empty")
	v.Check(MaxLength(content, commentContentMax), "content",
		fmt.Sprintf("must not exceed %d characters", commentContentMax))
}
