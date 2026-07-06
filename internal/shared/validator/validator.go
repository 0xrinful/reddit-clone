package validator

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

var (
	EmailRX = regexp.MustCompile(
		"^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$",
	)
	UsernameRX      = regexp.MustCompile("^[a-zA-Z0-9][a-zA-Z0-9_]{2,31}$")
	CommunityNameRX = regexp.MustCompile("^[a-zA-Z0-9][a-zA-Z0-9_]{2,31}$")
)

type Validator struct {
	Errors map[string]string
}

func New() *Validator {
	return &Validator{Errors: make(map[string]string)}
}

func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

func (v *Validator) AddError(key, message string) {
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = message
	}
}

func (v *Validator) Check(ok bool, key string, message string) {
	if !ok {
		v.AddError(key, message)
	}
}

func NotBlank(s string) bool {
	return strings.TrimSpace(s) != ""
}

func MinLength(s string, n int) bool {
	return utf8.RuneCountInString(s) >= n
}

func MaxLength(s string, n int) bool {
	return utf8.RuneCountInString(s) <= n
}

func Matches(s string, rx *regexp.Regexp) bool {
	return rx.MatchString(s)
}
