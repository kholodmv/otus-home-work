package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var (
	ErrInvalidString         = errors.New("invalid string")
	ErrInvalidFirstElement   = errors.New("first element should not be a number")
	ErrInvalidNumberElements = errors.New("should not be 2 consecutive numeric elements")
)

func Unpack(s string) (string, error) {
	var unpackStr strings.Builder
	for i := range s {
		if unicode.IsDigit(rune(s[0])) {
			return "", ErrInvalidFirstElement
		}
		if unicode.IsDigit(rune(s[i])) && unicode.IsDigit(rune(s[i+1])) {
			return "", ErrInvalidNumberElements
		}

		if !unicode.IsDigit(rune(s[i])) {
			if i+1 < len(s) && unicode.IsDigit(rune(s[i+1])) {
				digit, _ := strconv.Atoi(string(rune(s[i+1])))
				unpackStr.WriteString(strings.Repeat(string(s[i]), digit))
			} else {
				unpackStr.WriteRune(rune(s[i]))
			}
		}
	}

	return unpackStr.String(), nil
}
