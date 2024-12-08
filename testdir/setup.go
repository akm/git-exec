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
	t.Logf("srcDir: %s\n", srcDir)
	t.Logf("destDir: %s\n", destDir)

	absSrcDir, err := filepath.Abs(srcDir)
	if err != nil {
		t.Fatalf("Failed to get absolute path: %v", err)
	}
	t.Logf("absSrcDir: %s\n", absSrcDir)

	groundDir := filepath.Join(destDir, filepath.Base(absSrcDir)) // srcDir に . が指定された場合でもそのディレクトリ名を取得する
	t.Logf("groundDir: %s\n", groundDir)
	if err := os.MkdirAll(destDir, os.ModePerm); err != nil {
		t.Fatalf("Failed to create directory: %v", err)
	}

	// Copy srcDir to destDir
	copyDir(t, srcDir, groundDir)

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
func copyDir(t testing.TB, src string, dest string) {
	t.Logf("copyDir: %s -> %s\n", src, dest)

	entries, err := os.ReadDir(src)
	if err != nil {
		t.Fatalf("Failed to read directory: %v", err)
	}

	err = os.MkdirAll(dest, os.ModePerm)
	if err != nil {
		t.Fatalf("Failed to create directory: %v", err)
	}

	for _, entry := range entries {
		copyEntry(t, src, dest, entry)
	}
}

func copyEntry(t testing.TB, src string, dest string, entry fs.DirEntry) {
	if entry.Name() == ".git" {
		t.Logf("skip .git directory\n")
		return
	}

	srcPath := filepath.Join(src, entry.Name())
	destPath := filepath.Join(dest, entry.Name())

	if entry.IsDir() {
		copyDir(t, srcPath, destPath)
		return
	}

	t.Logf("copyEntry: %s -> %s\n", srcPath, destPath)

	reader, err := os.Open(srcPath)
	if err != nil {
		t.Fatalf("Failed to open source file: %v", err)
	}
	defer reader.Close()

	writer, err := os.Create(destPath)
	if err != nil {
		t.Fatalf("Failed to create destination file: %v", err)
	}
	defer writer.Close()

	if _, err := io.Copy(writer, reader); err != nil {
		t.Fatalf("Failed to copy file: %v", err)
	}
}
