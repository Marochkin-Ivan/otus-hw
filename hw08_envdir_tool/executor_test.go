package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	// Place your code here
	t.Run("can't set", func(t *testing.T) {
		code := RunCmd([]string{"cd", "../"}, Environment{
			"CANT_SET=": {
				Value:      "some value",
				NeedRemove: false,
			},
		})

		require.Equal(t, CantSetEnvVar, code)
	})
}
