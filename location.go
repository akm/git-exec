package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func getLocation() (string, error) {
	curDir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	rootDir, err := gitRootDir()
	if err != nil {
		return "", err
	}
	rootDirName := filepath.Base(rootDir)

	relPath, err := filepath.Rel(rootDir, curDir)
	if err != nil {
		return "", err
	}

	if relPath == "." {
		return rootDirName, nil
	} else if strings.HasPrefix(relPath, "./") {
		return rootDirName + relPath[1:], nil
	} else if strings.HasPrefix(relPath, "/") {
		return relPath, nil
	} else {
		return rootDirName + "/" + relPath, nil
	}
}

func gitRootDir() (string, error) {
	out, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}
