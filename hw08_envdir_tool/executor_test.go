package main

import (
	"testing"
)

func TestRunCmdExitCode(t *testing.T) {
	cmd := []string{"sh", "-c", "exit 42"}
	rc := RunCmd(cmd, Environment{})
	if rc != 42 {
		t.Errorf("Expected exit code 42, but got %d", rc)
	}
}
