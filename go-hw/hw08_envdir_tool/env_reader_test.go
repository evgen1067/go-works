package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	t.Run("Test reading dir", func(t *testing.T) {
		envs, err := ReadDir("testdata/env")
		require.NoError(t, err)

		require.Equal(t, `"hello"`, envs["HELLO"].Value)
		require.Equal(t, false, envs["HELLO"].NeedRemove)

		require.Equal(t, "bar", envs["BAR"].Value)
		require.Equal(t, false, envs["BAR"].NeedRemove)

		require.Equal(t, "   foo\nwith new line", envs["FOO"].Value)
		require.Equal(t, false, envs["FOO"].NeedRemove)

		require.Equal(t, "", envs["UNSET"].Value)
		require.Equal(t, true, envs["UNSET"].NeedRemove)

		require.Equal(t, "", envs["EMPTY"].Value)
		require.Equal(t, false, envs["EMPTY"].NeedRemove)
	})

	t.Run("test non-existent directory", func(t *testing.T) {
		_, err := ReadDir("fail/env")
		require.Error(t, err)
	})

	t.Run("test non-directory reading", func(t *testing.T) {
		_, err := ReadDir("env_reader.go")
		require.Error(t, err)
	})
}
