package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrDigitStart = errors.New("invalid string: can't start with a digit")
var ErrHasNumber = errors.New("invalid string: can't have a number")
var ErrIncorrectEscape = errors.New("invalid string: incorrect usage of escape characters")

func unpack(input string) (string, error) {
	res := strings.Builder{}
	runes := []rune(input)

	err := validate(runes)
	if err != nil {
		return "", err
	}

	for i := 0; i < len(runes); i++ {
		if unicode.IsLetter(runes[i]) || unicode.IsSpace(runes[i]) {
			if i < len(runes)-1 && unicode.IsDigit(runes[i+1]) {
				repeatCount, _ := strconv.Atoi(string(runes[i+1]))
				res.WriteString(strings.Repeat(string(runes[i]), repeatCount))
				i++ // перескакиваем, чтобы не проверять цифру
				continue
			} else {
				res.WriteString(string(runes[i]))
				continue
			}
		}
		if runes[i] == '\u005C' {
			if i < len(runes)-2 && unicode.IsDigit(runes[i+2]) {
				repeatCount, _ := strconv.Atoi(string(runes[i+2]))
				res.WriteString(strings.Repeat(string(runes[i+1]), repeatCount))
				i += 2 // перескакиваем, чтобы не проверять \ и цифру
				continue
			} else {
				res.WriteString(string(runes[i+1]))
				i++ // перескакиваем, чтобы не проверять \
				continue
			}
		}
	}

	return res.String(), nil
}

func validateEscapes(input []rune) ([]rune, error) {
	var res []rune
	for i := 0; i < len(input); i++ {
		if input[i] == '\u005C' {
			// за \ обязательно должна быть цифра или другой \
			if i+1 < len(input) && (unicode.IsDigit(input[i+1]) || input[i+1] == '\u005C') {
				res = append(res, 'a')
				i++
			} else { // если после \ ничего нет - это ошибка
				return nil, ErrIncorrectEscape
			}
		} else {
			res = append(res, input[i])
		}
	}

	return res, nil
}

func validate(input []rune) error {
	// сначала проверяем экранирование
	input, err := validateEscapes(input)
	if err != nil {
		return err
	}

	for i := 0; i < len(input)-1; i++ {
		// не может начинаться с цифры
		if i == 0 && unicode.IsDigit(input[i]) {
			return ErrDigitStart
		}

		// не может иметь числа
		if unicode.IsDigit(input[i]) && unicode.IsDigit(input[i+1]) {
			return ErrHasNumber
		}
	}

	return nil
}
