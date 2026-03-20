package validator

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

var EmailRX = regexp.MustCompile(
	"^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$",
)

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
