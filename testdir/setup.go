package testdir

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
)

func Setup(t testing.TB, srcDir, destDir string) func() {
	// srcDir を destDir にコピーして、コピーされたディレクトリにカレントディレクトリを移動する
	// 戻り値は カレントディレクトリを元のディレクトリに戻し、コピーされたディレクトリを削除する関数

	groundDir := filepath.Join(destDir, filepath.Base(srcDir))
	if err := os.MkdirAll(destDir, os.ModePerm); err != nil {
		t.Fatalf("Failed to create directory: %v", err)
	}

	// Copy srcDir to destDir
	if err := copyDir(srcDir, groundDir); err != nil {
		t.Fatalf("Failed to copy directory: %v", err)
	}

	// Save the current working directory
	origWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}

	// Change to the new directory
	if err := os.Chdir(groundDir); err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}

	// Return a function to restore the original directory and clean up
	return func() {
		// Change back to the original directory
		if err := os.Chdir(origWd); err != nil {
			t.Fatalf("Failed to change back to original directory: %v", err)
		}

		// Remove the copied directory
		if err := os.RemoveAll(groundDir); err != nil {
			t.Fatalf("Failed to remove directory %s: %v", groundDir, err)
		}
	}
}

// Helper function to copy a directory
func copyDir(src string, dest string) error {
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	err = os.MkdirAll(dest, os.ModePerm)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if err := copyEntry(src, dest, entry); err != nil {
			return err
		}
	}

	return nil
}

func copyEntry(src string, dest string, entry fs.DirEntry) error {
	srcPath := filepath.Join(src, entry.Name())
	destPath := filepath.Join(dest, entry.Name())

	if entry.IsDir() {
		if err := copyDir(srcPath, destPath); err != nil {
			return err
		}
		return nil
	}

	reader, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer reader.Close()

	writer, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer writer.Close()

	if _, err := io.Copy(writer, reader); err != nil {
		return err
	}
	return nil
}
