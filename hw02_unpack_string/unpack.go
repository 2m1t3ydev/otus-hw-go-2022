package hw02unpackstring

import (
	"errors"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(source string) (string, error) {
	var stringBuilder strings.Builder
	sourceRuneSlice := []rune(source)
	sourceLength := len(sourceRuneSlice)
	if sourceLength == 0 {
		return "", nil
	}
	var previous rune
	lastRuneIndex := sourceLength - 1
	for i, val := range sourceRuneSlice {
		isDigit := unicode.IsDigit(val)
		isPreviousDigit := unicode.IsDigit(previous)

		if i == 0 && isDigit {
			return "", ErrInvalidString
		}

		if isDigit && isPreviousDigit {
			return "", ErrInvalidString
		}

		if isDigit {
			times := int(val - '0')
			if times != 0 {
				stringBuilder.WriteString(strings.Repeat(string(previous), times))
			}
		} else {
			if previous != 0 && !isPreviousDigit {
				stringBuilder.WriteRune(previous)
			}

			if i == lastRuneIndex {
				stringBuilder.WriteRune(val)
			}
		}

		previous = val
	}
	return stringBuilder.String(), nil
}
