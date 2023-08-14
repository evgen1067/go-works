package main

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	t.Run("test offset is larger than the file size", func(t *testing.T) {
		dirName, err := os.MkdirTemp("", "dir")
		defer os.RemoveAll(dirName)
		require.NoError(t, err)
		err = Copy("testdata/empty.txt", dirName+"/out.txt", 1000, 0)
		require.Error(t, err)
		require.ErrorIs(t, ErrOffsetExceedsFileSize, err)
	})

	t.Run("test unknown length", func(t *testing.T) {
		dirName, err := os.MkdirTemp("", "dir")
		defer os.RemoveAll(dirName)
		require.NoError(t, err)
		err = Copy("/dev/urandom", dirName+"/out.txt", 1000, 0)
		require.Error(t, err)
		require.ErrorIs(t, ErrUnsupportedFile, err)
	})

	t.Run("test directory", func(t *testing.T) {
		dirName, err := os.MkdirTemp("", "dir")
		defer os.RemoveAll(dirName)
		require.NoError(t, err)
		err = Copy("testdata", dirName+"/out.txt", 1000, 0)
		require.Error(t, err)
		require.ErrorIs(t, ErrUnsupportedFile, err)
	})

	t.Run("test file with EOF", func(t *testing.T) {
		dirName, err := os.MkdirTemp("", "dir")
		defer os.RemoveAll(dirName)
		from := "testdata/out_offset0_limit1000.txt"
		to := dirName + "/out.txt"
		require.NoError(t, err)
		err = Copy(from, to, 0, 1000)
		require.NoError(t, err)

		f1, err := os.Open(from)
		require.NoError(t, err)
		defer f1.Close()
		f2, err := os.Open(to)
		require.NoError(t, err)
		defer f2.Close()

		stat, err := f1.Stat()
		require.NoError(t, err)

		for {
			b1 := make([]byte, stat.Size())
			_, err1 := f1.Read(b1)
			b2 := make([]byte, stat.Size())
			_, err2 := f2.Read(b2)
			if err1 != nil || err2 != nil {
				require.Equal(t, true, err1 == io.EOF && err2 == io.EOF)
				break
			}
			require.Equal(t, true, bytes.Equal(b1, b2))
		}
	})
}
