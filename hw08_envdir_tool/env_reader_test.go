package main

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	t.Run("success read variables", func(t *testing.T) {
		env, err := ReadDir("testdata/env")

		expected := Environment{
			"BAR":   EnvValue{Value: "bar", NeedRemove: false},
			"EMPTY": EnvValue{Value: "", NeedRemove: false},
			"FOO":   EnvValue{Value: "   foo\nwith new line", NeedRemove: false},
			"HELLO": EnvValue{Value: "\"hello\"", NeedRemove: false},
			"UNSET": EnvValue{Value: "", NeedRemove: true},
		}

		require.Equal(t, expected, env)
		require.NoError(t, err)
	})

	t.Run("directory not specified", func(t *testing.T) {
		env, err := ReadDir("")

		require.Len(t, env, 0)
		require.NotNil(t, err)
	})

	t.Run("invalid variables name", func(t *testing.T) {
		dir := "testdata/tmp"
		err := os.Mkdir(dir, os.FileMode(0o755))
		require.NoError(t, err)
		defer func() {
			err = os.RemoveAll(dir)
			require.NoError(t, err)
		}()

		f, err := os.Create(path.Join(dir, "X=Y"))
		require.NoError(t, err)
		defer func() {
			err = f.Close()
			require.NoError(t, err)
		}()

		env, err := ReadDir(dir)
		require.Len(t, env, 0)
	})
}
