package main

import (
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestCopy(t *testing.T) {
	// Place your code here.
	t.Run("copy full file without params",
		func(t *testing.T) {
			tmp, _ := os.CreateTemp("", "output")
			defer os.Remove(tmp.Name())

			mustBe, _ := os.ReadFile("testdata/out_offset0_limit0.txt")

			err := Copy("testdata/input.txt", tmp.Name(), 0, 0)
			actual, _ := os.ReadFile(tmp.Name())

			require.Nil(t, err)
			require.Equal(t, mustBe, actual)
		})

	t.Run("offset > file length",
		func(t *testing.T) {
			err := Copy("testdata/input.txt", "", 100_000_000, 0)

			require.EqualError(t, err, ErrOffsetExceedsFileSize.Error())
		})
}
