package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	t.Run("no command", func(t *testing.T) {
		cmd := []string{}
		env := Environment{}

		errorCode := RunCmd(cmd, env)

		require.Equal(t, 1, errorCode)
	})

	t.Run("with command simple", func(t *testing.T) {
		cmd := []string{"pwd"}
		env := Environment{}

		errorCode := RunCmd(cmd, env)

		require.Equal(t, 0, errorCode)
	})

	t.Run("setting env var when there is no existing one", func(t *testing.T) {
		cmd := []string{"pwd"}
		env := Environment{"NOT_EXIST_ENV": EnvValue{Value: "NOT_EXIST_ENV_VALUE"}}

		errorCode := RunCmd(cmd, env)

		require.Equal(t, 0, errorCode)
		require.Contains(t, os.Environ(), "NOT_EXIST_ENV=NOT_EXIST_ENV_VALUE")
	})

	t.Run("setting environment variable when there is already", func(t *testing.T) {
		cmd := []string{"pwd"}
		env := Environment{"EXIST_ENV": EnvValue{Value: "MODIFY_EXIST_ENV_VALUE"}}

		err := os.Setenv("EXIST_ENV", "EXIST_ENV_VALUE")
		require.NoError(t, err)

		errorCode := RunCmd(cmd, env)

		require.Equal(t, 0, errorCode)
		require.Contains(t, os.Environ(), "EXIST_ENV=MODIFY_EXIST_ENV_VALUE")
	})

	t.Run("delete when there is no environment variable", func(t *testing.T) {
		cmd := []string{"pwd"}
		env := Environment{"NOT_EXIST_ENV": EnvValue{Value: "", NeedRemove: true}}

		errorCode := RunCmd(cmd, env)

		require.Equal(t, 0, errorCode)
		require.NotContains(t, os.Environ(), "NOT_EXIST_ENV")
	})

	t.Run("deleting environment variable when there is already", func(t *testing.T) {
		cmd := []string{"pwd"}
		env := Environment{"EXIST_ENV": EnvValue{Value: "", NeedRemove: true}}

		err := os.Setenv("EXIST_ENV", "EXIST_ENV_VALUE")
		require.NoError(t, err)

		errorCode := RunCmd(cmd, env)

		require.Equal(t, 0, errorCode)
		require.NotContains(t, os.Environ(), "EXIST_ENV")
	})
}
