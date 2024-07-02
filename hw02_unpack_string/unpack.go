package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

const backSlash = '\u005C'

var (
	ErrDigitStart      = errors.New("invalid string: can't start with a digit")
	ErrHasNumber       = errors.New("invalid string: can't have a number")
	ErrIncorrectEscape = errors.New("invalid string: incorrect usage of escape characters")
)

func Unpack(input string) (string, error) {
	res := strings.Builder{}
	runes := []rune(input)

	if len(runes) == 0 {
		return "", nil
	}

	for i := 0; i < len(runes); i++ {
		if i == 0 && unicode.IsDigit(runes[i]) {
			return "", ErrDigitStart
		}

		if runes[i] == backSlash { // экранирование
			if i+1 < len(runes) {
				if unicode.IsDigit(runes[i+1]) || runes[i+1] == backSlash {
					if i+2 < len(runes) && unicode.IsDigit(runes[i+2]) {
						repeatCount, _ := strconv.Atoi(string(runes[i+2]))
						res.WriteString(strings.Repeat(string(runes[i+1]), repeatCount))
						i += 2 // Пропускаем следующий символ и цифру после экранирования
					} else {
						res.WriteRune(runes[i+1])
						i++ // Пропускаем следующий символ после экранирования
					}
				} else {
					return "", ErrIncorrectEscape
				}
			} else {
				return "", ErrIncorrectEscape
			}
			continue
		}

		if unicode.IsLetter(runes[i]) || unicode.IsSpace(runes[i]) {
			if i < len(runes)-1 && unicode.IsDigit(runes[i+1]) {
				repeatCount, _ := strconv.Atoi(string(runes[i+1]))
				res.WriteString(strings.Repeat(string(runes[i]), repeatCount))
				i++ // Пропускаем цифру
			} else {
				res.WriteRune(runes[i])
			}
			continue
		}

		if unicode.IsDigit(runes[i]) {
			if i == 0 || unicode.IsDigit(runes[i-1]) {
				return "", ErrHasNumber
			}
		}
	}

	return res.String(), nil
}
