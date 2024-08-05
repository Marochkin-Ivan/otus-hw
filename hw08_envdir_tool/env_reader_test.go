package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	// Place your code here
	t.Run("simple", func(t *testing.T) {
		env, err := ReadDir("./testdata/env")
		expected := Environment{
			"BAR": {
				Value:      "bar",
				NeedRemove: false,
			},
			"EMPTY": {
				Value:      "",
				NeedRemove: false,
			},
			"FOO": {
				Value:      "   foo\nwith new line",
				NeedRemove: false,
			},
			"HELLO": {
				Value:      "\"hello\"",
				NeedRemove: false,
			},
			"UNSET": {
				Value:      "",
				NeedRemove: true,
			},
		}

		require.NoError(t, err)
		require.Equal(t, expected, env)
	})

	t.Run("wrong directory path", func(t *testing.T) {
		_, err := ReadDir("./testdataaaaa/env")

		require.ErrorIs(t, err, ErrCantReadDir)
	})

	t.Run("filename with '='", func(t *testing.T) {
		f, err := os.CreateTemp("./testdata/env", "with_=_*")
		defer os.Remove(f.Name())
		require.NoError(t, err)

		_, err = ReadDir("./testdata/env")

		require.ErrorIs(t, err, ErrHasForbiddenSymbols)
	})

	t.Run("incorrect file", func(t *testing.T) {
		f, err := os.CreateTemp("./testdata/env", "tmp_*")
		defer os.Remove(f.Name())
		require.NoError(t, err)

		_, err = f.Write([]byte("some data"))
		require.NoError(t, err)

		err = f.Chmod(0o222) // only write permissions
		require.NoError(t, err)

		_, err = ReadDir("./testdata/env")
		require.ErrorIs(t, err, ErrCantOpenFile)
	})
}
