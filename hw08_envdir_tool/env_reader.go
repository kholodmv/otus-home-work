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
	env := make(Environment)

	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		if f.IsDir() || strings.Contains(f.Name(), "=") {
			continue
		}

		path := filepath.Join(dir, f.Name())
		content, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}

		if len(content) == 0 {
			env[f.Name()] = EnvValue{NeedRemove: true}
		} else {
			value := string(content)
			value = strings.TrimRight(value, " \t")
			value = string(bytes.ReplaceAll([]byte(value), []byte{0}, []byte{'\n'}))
			env[f.Name()] = EnvValue{Value: value, NeedRemove: false}
		}
	}

	return env, nil
}
