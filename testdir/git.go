package testdir

import (
	"os"
	"path/filepath"
	"testing"
)

// GitRoot searches for the directory containing the .git directory starting from the current directory and moving upwards.
// If the .git directory is found, it returns the path to that directory.
// If the .git directory is not found, it fails the test.
func GitRoot(t *testing.T) string {
	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, ".git")); err == nil {
			return dir
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatalf(".git directory not found")
		}
		dir = parent
	}
}

// FromGitRoot returns the absolute path by joining the Git root directory with the provided relative path.
// If the .git directory is not found, it fails the test.
func FromGitRoot(t *testing.T, relpath string) string {
	return filepath.Join(GitRoot(t), relpath)
}
