package main

import (
	"errors"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	if len(cmd) == 0 {
		return 1
	}
	command := cmd[0]

	var params []string
	if len(cmd) > 1 {
		params = cmd[1:]
	}

	setEnvVar(env)
	return execCmd(command, params)
}

func setEnvVar(env Environment) {
	if len(env) == 0 {
		return
	}

	for envName, envValue := range env {
		if _, ok := os.LookupEnv(envName); ok {
			os.Unsetenv(envName)
		}

		if !envValue.NeedRemove {
			os.Setenv(envName, envValue.Value)
		}
	}
}

func execCmd(command string, args []string) int {
	exeCmd := exec.Command(command, args...)
	exeCmd.Stdin = os.Stdin
	exeCmd.Stdout = os.Stdout
	exeCmd.Stderr = os.Stderr

	if err := exeCmd.Run(); err != nil {
		if e := (&exec.ExitError{}); errors.As(err, &e) {
			return e.ExitCode()
		}
		return -1
	}

	return 0
}
