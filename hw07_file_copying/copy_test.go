package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	// Place your code here.
	t.Run("copy full file without params",
		func(t *testing.T) {
			tmp, err := os.CreateTemp("", "output")
			require.NoError(t, err)
			defer os.Remove(tmp.Name())

			mustBe, err := os.ReadFile("testdata/out_offset0_limit0.txt")
			require.NoError(t, err)

			err = Copy("testdata/input.txt", tmp.Name(), 0, 0)
			require.NoError(t, err)

			actual, err := os.ReadFile(tmp.Name())
			require.NoError(t, err)

			require.Equal(t, mustBe, actual)
		})

	t.Run("offset > file length",
		func(t *testing.T) {
			tmp, err := os.CreateTemp("", "output")
			require.NoError(t, err)
			defer os.Remove(tmp.Name())

			err = Copy("testdata/input.txt", tmp.Name(), 100_000_000, 0)
			require.EqualError(t, err, ErrOffsetExceedsFileSize.Error())
		})
}
