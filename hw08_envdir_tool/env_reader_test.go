package main

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "test")
	if err != nil {
		require.NoError(t, err, "Failed to create temp directory")
	}
	defer os.RemoveAll(tempDir)

	testCases := []struct {
		name     string
		files    map[string]string
		expected Environment
	}{
		{
			name: "Normal files",
			files: map[string]string{
				"FOO": "bar",
				"BAZ": "qux",
			},
			expected: Environment{
				"FOO": EnvValue{Value: "bar", NeedRemove: false},
				"BAZ": EnvValue{Value: "qux", NeedRemove: false},
			},
		},
		{
			name: "Empty file",
			files: map[string]string{
				"EMPTY": "",
			},
			expected: Environment{
				"EMPTY": EnvValue{NeedRemove: true},
			},
		},
		{
			name: "File with spaces and tabs",
			files: map[string]string{
				"SPACES": "value with spaces  \t  ",
			},
			expected: Environment{
				"SPACES": EnvValue{Value: "value with spaces", NeedRemove: false},
			},
		},
		{
			name: "File with null bytes",
			files: map[string]string{
				"NULL": "before\x00after",
			},
			expected: Environment{
				"NULL": EnvValue{Value: "before\nafter", NeedRemove: false},
			},
		},
		{
			name: "File with '=' in name",
			files: map[string]string{
				"INVALID=NAME": "value",
			},
			expected: Environment{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			for name, content := range tc.files {
				err := os.WriteFile(filepath.Join(tempDir, name), []byte(content), 0o644)
				if err != nil {
					require.NoError(t, err, "Failed to create test file")
				}
			}

			result, err := ReadDir(tempDir)
			if err != nil {
				require.NoError(t, err, "ReadDir failed")
			}

			if !reflect.DeepEqual(result, tc.expected) {
				require.Equal(t, tc.expected, result, "Unexpected result from ReadDir")
			}

			dir, _ := os.Open(tempDir)
			names, _ := dir.Readdirnames(-1)
			for _, name := range names {
				os.Remove(filepath.Join(tempDir, name))
			}
		})
	}
}

func TestReadDirNonExistentError(t *testing.T) {
	_, err := ReadDir("/non/test")
	if err == nil {
		require.Error(t, err, "Expected error for non-existent directory")
	}
}
