package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	var offset int64
	var limit int64
	testDir := "testdata"
	tmpDir := t.TempDir()
	entries, _ := os.ReadDir(testDir)

	t.Run("test -> not regular file", func(t *testing.T) {
		offset = 10
		limit = 10
		tmpFile, _ := os.CreateTemp(tmpDir, "*.txt")
		err := Copy("/dev/urandom", tmpFile.Name(), offset, limit)
		require.Error(t, err, "expected return ErrUnsupportedFile")
	})

	t.Run("test -> offset more than file length", func(t *testing.T) {
		offset = 1000000
		limit = 1
		for _, e := range entries {
			tmpFile, _ := os.CreateTemp(tmpDir, "*.txt")
			err := Copy(filepath.Join(testDir, e.Name()), tmpFile.Name(), offset, limit)
			require.Error(t, err, "expected return ErrOffsetExceedsFileSize")

			tmpFile.Close()
		}
	})

	t.Run("test -> copy all test files", func(t *testing.T) {
		offset = 0
		limit = 0
		for _, e := range entries {
			srcFileStat, _ := e.Info()
			tmpFile, _ := os.CreateTemp(tmpDir, "*.txt")
			err := Copy(filepath.Join(testDir, e.Name()), tmpFile.Name(), offset, limit)

			dstFileStat, _ := tmpFile.Stat()
			require.NoError(t, err, "expected return nil")
			require.Equal(t, srcFileStat.Size(), dstFileStat.Size(), "expected files must be the same size")

			tmpFile.Close()
		}
	})

	t.Run("test -> check limit", func(t *testing.T) {
		offset = 0
		limit = 100
		for _, e := range entries {
			tmpFile, _ := os.CreateTemp(tmpDir, "*.txt")
			err := Copy(filepath.Join(testDir, e.Name()), tmpFile.Name(), offset, limit)

			dstFileStat, _ := tmpFile.Stat()
			require.NoError(t, err, "expected return nil")
			require.LessOrEqual(t, dstFileStat.Size(), limit, "expected file size must be less or equal to limit")

			tmpFile.Close()
		}
	})

	t.Run("test -> check offset", func(t *testing.T) {
		offset = 100
		limit = 0
		for _, e := range entries {
			srcFileStat, _ := e.Info()
			if srcFileStat.Size() <= offset {
				continue
			}
			tmpFile, _ := os.CreateTemp(tmpDir, "*.txt")
			err := Copy(filepath.Join(testDir, e.Name()), tmpFile.Name(), offset, limit)

			dstFileStat, _ := tmpFile.Stat()
			require.NoError(t, err, "expected return nil")
			require.Equal(t, srcFileStat.Size()-offset, dstFileStat.Size(), "expected copied file must be lesser by offset")

			tmpFile.Close()
		}
	})
}
