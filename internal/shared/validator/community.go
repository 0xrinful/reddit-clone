package validator

import "fmt"

const (
	communityNameMin        = 3
	communityNameMax        = 32
	communityDescriptionMax = 500
)

func ValidateCommunityName(v *Validator, name string) {
	v.Check(
		MinLength(name, communityNameMin),
		"name",
		fmt.Sprintf("must be at least %d characters", communityNameMin),
	)
	v.Check(
		MaxLength(name, communityNameMax),
		"name",
		fmt.Sprintf("must not exceed %d characters", communityNameMax),
	)
	v.Check(
		Matches(name, CommunityNameRX),
		"name",
		"must contain only letters, numbers, and underscores, and must start with a letter or number",
	)
}

func ValidateCommunityDescription(v *Validator, description string) {
	v.Check(NotBlank(description), "description", "must not be empty")
	v.Check(
		MaxLength(description, communityDescriptionMax),
		"description",
		fmt.Sprintf("must not exceed %d characters", communityDescriptionMax),
	)
}
