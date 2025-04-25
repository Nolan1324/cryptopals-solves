package util

import "strings"

func IsAlpha(c byte) bool {
	return ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z')
}

func IsNumeric(c byte) bool {
	return '0' <= c && c <= '9'
}

func IsAlphanumeric(c byte) bool {
	return IsAlpha(c) || IsNumeric(c)
}

func IsPunctuation(c byte) bool {
	return strings.Contains(" !@#$%^&*()-=_+[]{}\\|;':\",./<>?", string(rune(c)))
}

func IsWordChar(c byte) bool {
	return IsAlphanumeric(c) || IsPunctuation(c)
}
