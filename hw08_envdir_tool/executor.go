package main

import (
	"errors"
	"log/slog"
	"os"
	"os/exec"
)

const (
	Success = iota
	Failure
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	if len(cmd) == 0 {
		slog.Info("no command to execute found")
		return Success
	}

	name, args := cmd[0], cmd[1:]

	command := exec.Command(name, args...)
	command.Stderr = os.Stderr
	command.Stdout = os.Stdout
	command.Stdin = os.Stdin

	for k, v := range env {
		err := os.Unsetenv(k)
		if err != nil {
			slog.Error("error unsetting variable from env", "error", err)
			return Failure
		}

		if !v.NeedRemove {
			err = os.Setenv(k, v.Value)
			if err != nil {
				slog.Error("error setting variable into env", "error", err)
				return Failure
			}
		}
	}

	if err := command.Run(); err != nil {
		var exitError *exec.ExitError
		if errors.As(err, &exitError) {
			return exitError.ExitCode()
		}

		slog.Error("error executing command", "error", err)
		return Failure
	}

	return command.ProcessState.ExitCode()
}
