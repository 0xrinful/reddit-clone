package validator

func ValidateToken(v *Validator, tokenPlaintext string) {
	v.Check(tokenPlaintext != "", "token", "must be provided")
	v.Check(len(tokenPlaintext) == 22, "token", "must be 26 bytes long")
}
