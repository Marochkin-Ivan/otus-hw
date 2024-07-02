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
	runes := []rune(input)

	if len(runes) == 0 {
		return "", nil
	}

	if unicode.IsDigit(runes[0]) {
		return "", ErrDigitStart
	}

	var res strings.Builder
	for i := 0; i < len(runes); i++ {
		if runes[i] == backSlash { // экранирование
			escaped, newIndex, err := handleEscape(runes, i)
			if err != nil {
				return "", err
			}

			res.WriteString(escaped)
			i = newIndex
			continue
		}

		if unicode.IsLetter(runes[i]) || unicode.IsSpace(runes[i]) {
			if i < len(runes)-1 && unicode.IsDigit(runes[i+1]) {
				repeatCount, _ := strconv.Atoi(string(runes[i+1]))
				res.WriteString(strings.Repeat(string(runes[i]), repeatCount))
				i++ // Пропускаем кол-во повторений
			} else {
				res.WriteRune(runes[i])
			}
			continue
		}

		if unicode.IsDigit(runes[i]) {
			return "", ErrHasNumber
		}
	}

	return res.String(), nil
}

func handleEscape(runes []rune, index int) (string, int, error) {
	if index+1 >= len(runes) {
		return "", index, ErrIncorrectEscape
	}

	nextChar := runes[index+1]
	if unicode.IsDigit(nextChar) || nextChar == backSlash { // экранируются только цифры и слеши
		if index+2 < len(runes) && unicode.IsDigit(runes[index+2]) {
			repeatCount, _ := strconv.Atoi(string(runes[index+2]))
			return strings.Repeat(string(nextChar), repeatCount),
				index + 2, // пропускаем экранируемый символ и кол-во повторений
				nil
		}
		return string(nextChar),
			index + 1, // пропускаем экранируемый символ
			nil
	}

	return "", index, ErrIncorrectEscape
}
