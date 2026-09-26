package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	if len(entries) == 0 {
		return nil, nil
	}

	env := make(Environment)

	for _, entry := range entries {
		stat, err := entry.Info()
		if err != nil {
			return nil, err
		}

		if regular := stat.Mode().IsRegular(); !regular {
			continue
		}

		if strings.ContainsRune(entry.Name(), '=') {
			continue
		}

		if stat.Size() == 0 {
			env[entry.Name()] = EnvValue{NeedRemove: true}
			continue
		}

		fullBytes, err := os.ReadFile(filepath.Join(dir, entry.Name()))
		if err != nil {
			return nil, err
		}

		if len(fullBytes) == 0 {
			env[entry.Name()] = EnvValue{NeedRemove: true}
			continue
		}

		firstLineBytes := bytes.ReplaceAll(bytes.Split(fullBytes, []byte("\n"))[0], []byte("\x00"), []byte("\n"))
		firstLine := strings.TrimRight(string(firstLineBytes), " \t")

		env[entry.Name()] = EnvValue{firstLine, false}
	}

	return env, nil
}
