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

func IsSentenceAscii(c byte) bool {
	return IsAlphanumeric(c) || strings.Contains(" \"',.\n\r", string(rune(c)))
}

func IsVisibleAscii(c byte) bool {
	return (32 <= c && c <= 126) || c == '\r' || c == '\n'
}
