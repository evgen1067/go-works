package main

import (
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	t.Run("test run cmd", func(t *testing.T) {
		text := "HELLO is (\"hello\")\nBAR is ()\nFOO is ()\nUNSET is (TEST)\nEMPTY is ()\narguments are arg1=1 arg2=2\n"
		cmd := []string{"/bin/bash", "testdata/echo.sh", "arg1=1", "arg2=2"}
		envs := make(Environment)
		envs["HELLO"] = EnvValue{
			Value:      `"hello"`,
			NeedRemove: false,
		}
		r, w, err := os.Pipe()
		require.NoError(t, err)

		err = os.Setenv("UNSET", "TEST")
		require.NoError(t, err)
		err = os.Setenv("HELLO", "TEST")
		require.NoError(t, err)

		os.Stdout = w

		code := RunCmd(cmd, envs)

		err = w.Close()
		require.NoError(t, err)

		require.Equal(t, 0, code)

		lines, err := io.ReadAll(r)
		require.NoError(t, err)
		require.Equal(t, text, string(lines))
	})

	t.Run("test non-existent directory", func(t *testing.T) {
		cmd := []string{"/bin/bash", "env/echo.sh", "arg1=1", "arg2=2"}
		envs := make(Environment)
		envs["HELLO"] = EnvValue{
			Value:      `"hello"`,
			NeedRemove: false,
		}
		code := RunCmd(cmd, envs)
		require.Equal(t, 1, code)
	})

	t.Run("test run cmd", func(t *testing.T) {
		text := "HELLO is (TEST)\nBAR is ()\nFOO is ()\nUNSET is ()\nEMPTY is ()\narguments are arg1=1 arg2=2\n"
		cmd := []string{"/bin/bash", "testdata/echo.sh", "arg1=1", "arg2=2"}

		envs := make(Environment)
		envs["UNSET"] = EnvValue{
			Value:      "UNSET",
			NeedRemove: true,
		}
		r, w, err := os.Pipe()
		require.NoError(t, err)

		err = os.Setenv("UNSET", "TEST")

		require.NoError(t, err)

		os.Stdout = w

		code := RunCmd(cmd, envs)

		err = w.Close()
		require.NoError(t, err)

		require.Equal(t, 0, code)

		lines, err := io.ReadAll(r)
		require.NoError(t, err)
		require.Equal(t, text, string(lines))
	})
}
