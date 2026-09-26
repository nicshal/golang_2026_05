package main

import (
	"log/slog"
	"os"
)

func main() {
	args := os.Args
	if len(args) < 3 {
		slog.Error("error command call. Usage: go-evndir /path/to/evndir command arg1 arg2...")
		os.Exit(1)
	}
	dir, cmd := args[1], args[2:]

	env, err := ReadDir(dir)
	if err != nil {
		slog.Error("error parsing environment variables", "error", err)
		os.Exit(1)
	}

	returnCode := RunCmd(cmd, env)
	os.Exit(returnCode)
}
