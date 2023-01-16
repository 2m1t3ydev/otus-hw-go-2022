package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	t.Run("not exist directory", func(t *testing.T) {
		dirPath := "./testdata/notexist"

		env, err := ReadDir(dirPath)

		require.Error(t, err)
		require.Equal(t, len(env), 0)
	})

	t.Run("empty directory", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "")
		require.NoError(t, err)
		defer os.Remove(tmpDir)

		env, err := ReadDir(tmpDir)

		require.NoError(t, err)
		require.Equal(t, len(env), 0)
	})

	t.Run("file name with '='", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "")
		require.NoError(t, err)
		defer os.RemoveAll(tmpDir)

		tmpFile, err := os.CreateTemp(tmpDir, "=")
		require.NoError(t, err)
		tmpFile.Close()

		env, err := ReadDir(tmpDir)

		require.NoError(t, err)
		require.Equal(t, len(env), 0)
	})

	t.Run("directory with file and subdirectory", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "")
		require.NoError(t, err)
		defer os.RemoveAll(tmpDir)

		tmpSubDir, err := os.MkdirTemp(tmpDir, "")
		_ = tmpSubDir
		require.NoError(t, err)

		tmpFile, err := os.CreateTemp(tmpDir, "")
		require.NoError(t, err)
		tmpFile.Close()

		env, err := ReadDir(tmpDir)

		require.NoError(t, err)
		require.Equal(t, len(env), 1)
	})

	t.Run("check file content", func(t *testing.T) {
		dirPath := "./testdata/env"
		expected := Environment{
			"BAR":   EnvValue{Value: "bar", NeedRemove: false},
			"EMPTY": EnvValue{Value: "", NeedRemove: false},
			"FOO":   EnvValue{Value: "   foo\nwith new line", NeedRemove: false},
			"HELLO": EnvValue{Value: "\"hello\"", NeedRemove: false},
			"UNSET": EnvValue{Value: "", NeedRemove: true},
		}

		env, err := ReadDir(dirPath)

		require.NoError(t, err)
		require.Equal(t, expected, env)
	})
}
