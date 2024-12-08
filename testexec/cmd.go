package testexec

import (
	"os"
	"os/exec"
	"testing"
)

func Run(t testing.TB, command string, args ...string) {
	t.Helper()

	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to run %s %v: %s", command, args, err)
	}
}

func Stdout(t testing.TB, command string, args ...string) string {
	t.Helper()

	cmd := exec.Command(command, args...)
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("failed to run %s %v: %s", command, args, err)
	}
	return string(out)
}
