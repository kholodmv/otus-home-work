package main

import (
	"fmt"
	"os"
)

func main() {
	envDir := os.Args[1]
	command := os.Args[2:]

	env, err := ReadDir(envDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading environment: %v\n", err)
		os.Exit(1)
	}

	code := RunCmd(command, env)
	os.Exit(code)
}
