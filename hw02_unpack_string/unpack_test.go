package hw02unpackstring

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnpack(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{input: "a4bc2d5e", expected: "aaaabccddddde"},
		{input: "abccd", expected: "abccd"},
		{input: "", expected: ""},
		{input: "aaa0b", expected: "aab"},
		{input: "a2", expected: "aa"},
		{input: "\n2", expected: "\n\n"},
		{input: "a", expected: "a"},
		// uncomment if task with asterisk completed
		{input: `qwe\4\5`, expected: `qwe45`},
		{input: `qwe\45`, expected: `qwe44444`},
		{input: `qwe\\5`, expected: `qwe\\\\\`},
		{input: `qwe\\\3`, expected: `qwe\3`},
		{input: `\\\\`, expected: `\\`},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			result, err := unpack(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestUnpackInvalidString(t *testing.T) {
	invalidStrings := []string{"3abc", "45", `2\3\4`}
	for _, tc := range invalidStrings {
		tc := tc
		t.Run(tc, func(t *testing.T) {
			_, err := unpack(tc)
			require.Truef(t, errors.Is(err, ErrDigitStart), "actual error %q", err)
		})
	}

	invalidStrings = []string{"abc34", "rrr45sdf", `\2\3\445`}
	for _, tc := range invalidStrings {
		tc := tc
		t.Run(tc, func(t *testing.T) {
			_, err := unpack(tc)
			require.Truef(t, errors.Is(err, ErrHasNumber), "actual error %q", err)
		})
	}

	invalidStrings = []string{`\\\`, `qq\q2`, `\\2\&`}
	for _, tc := range invalidStrings {
		tc := tc
		t.Run(tc, func(t *testing.T) {
			_, err := unpack(tc)
			require.Truef(t, errors.Is(err, ErrIncorrectEscape), "actual error %q", err)
		})
	}
}
