package main

import (
	"errors"
	"os"
	"os/exec"
)

const (
	CantSetEnvVar      = 501
	UndefinedErrorCode = 505
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	// Place your code here.
	for envName, envVal := range env {
		if envVal.NeedRemove {
			_ = os.Unsetenv(envName)
			continue
		}

		err := os.Setenv(envName, envVal.Value)
		if err != nil {
			return CantSetEnvVar
		}
	}

	command := exec.Command(cmd[0], cmd[1:]...)
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	command.Env = os.Environ()

	err := command.Run()
	if err != nil {
		var exitError *exec.ExitError
		if errors.As(err, &exitError) {
			returnCode = exitError.ExitCode()
		} else {
			return UndefinedErrorCode
		}
	}

	return
}
