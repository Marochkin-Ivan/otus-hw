package hw03frequencyanalysis

import (
	"github.com/wasilibs/go-re2"
	"strings"
)

const puncts = `,.!?;:-"'¡¿`

var wordSurroundedPunctuation = re2.MustCompile(`^[[:punct:]]+[[:graph:]]+[[:punct:]]+$`)

func Top10(input string) []string {
	words := strings.Fields(input)

	rating := make(map[string]int)
	for _, word := range words {
		if wordSurroundedPunctuation.MatchString(word) {
			word = correctWord(word)
		}
		rating[word]++
	}

	return nil
}

func correctWord(word string) string {
	return strings.Trim(strings.ToLower(word), puncts)
}
