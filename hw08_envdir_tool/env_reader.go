package main

import (
	"bufio"
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
			env[f.Name()] = EnvValue{Value: "", NeedRemove: true}
		} else {
			value, err := getValueFromFile(filepath.Join(dir, f.Name()))
			if err != nil {
				return nil, err
			}
			env[f.Name()] = EnvValue{Value: value, NeedRemove: false}
		}
	}

	return env, nil
}

func getValueFromFile(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	scanner.Scan()
	v := scanner.Bytes()
	v = bytes.ReplaceAll(v, []byte("\x00"), []byte("\n"))
	value := strings.TrimRight(string(v), " \t\n")

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return value, nil
}
