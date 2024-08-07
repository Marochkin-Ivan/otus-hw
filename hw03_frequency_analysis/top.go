package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

const (
	punctuation = `,.!?;:-"'¡¿`
	dash        = "-"
)

var wordSurroundedPunctuation = regexp.MustCompile(`^[[:punct:]]*[^[:punct:]]+[[:punct:]]*$`)

type rating map[string]int

func (r rating) keys() []string {
	keys := make([]string, 0, len(r))

	for k := range r {
		keys = append(keys, k)
	}

	return keys
}

func (r rating) sortByKeysAndValue() []string {
	keys := r.keys()

	sort.SliceStable(keys, func(i, j int) bool {
		if r[keys[i]] == r[keys[j]] {
			return keys[i] < keys[j]
		}
		return r[keys[i]] > r[keys[j]]
	})

	return keys
}

func Top10(input string) []string {
	words := strings.Fields(input)

	rate := make(rating)
	for _, word := range words {
		if wordSurroundedPunctuation.MatchString(word) {
			word = correctWord(word)
		}
		if word == dash {
			continue
		}
		rate[word]++
	}

	res := rate.sortByKeysAndValue()

	if len(res) > 10 {
		res = res[:10]
	}

	return res
}

func correctWord(word string) string {
	return strings.Trim(strings.ToLower(word), punctuation)
}
