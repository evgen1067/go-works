package main

import (
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmdArgs []string, env Environment) (returnCode int) {
	name, arg := cmdArgs[0], cmdArgs[1:]
	cmd := exec.Command(name, arg...)
	// стандартные потоки ввода/вывода/ошибок пробрасывались в вызываемую программу;
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Env = os.Environ()
	for fName, val := range env {
		if !val.NeedRemove {
			cmd.Env = append(cmd.Env, fName+"="+val.Value)
		} else {
			cmd.Env = append(cmd.Env, fName+"=")
		}
	}
	err := cmd.Run()
	if err != nil {
		return 1
	}

	return 0
}
