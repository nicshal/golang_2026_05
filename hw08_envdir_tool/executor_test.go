package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	env := make(Environment)

	t.Run("empty command", func(t *testing.T) {
		code := RunCmd([]string{}, env)
		require.Equal(t, 0, code)
	})

	t.Run("invalid command", func(t *testing.T) {
		code := RunCmd([]string{"invalid_command"}, env)
		require.NotEqual(t, 0, code)
	})

	t.Run("success", func(t *testing.T) {
		code := RunCmd([]string{"echo", "777"}, env)
		require.Equal(t, 0, code)
	})

	env["TST777"] = EnvValue{"777", false}

	t.Run("success with env", func(t *testing.T) {
		code := RunCmd([]string{"ls", "-la"}, env)
		require.Equal(t, 0, code)
		require.Contains(t, os.Environ(), "TST777=777")
	})
}
