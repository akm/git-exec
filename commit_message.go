package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var commitPrefix = func() string {
	if prefix := os.Getenv("GIT_EXEC_COMMIT_PREFIX"); prefix != "" {
		return prefix
	}
	return "🤖 @%s $"
}()

func buildCommitMessage(command *Command) string {
	firstLine := buildCommitMessageFirstLine(command.Envs, command.Args)
	return fmt.Sprintf("%s\n\n%s\n",
		firstLine,
		command.Output.String(),
	)
}

func buildCommitMessageFirstLine(envs []string, commandArgs []string) string {
	commandParts := make([]string, len(commandArgs))
	for i, arg := range commandArgs {
		if strings.Contains(arg, " ") && !(strings.HasPrefix(arg, "'") && strings.HasSuffix(arg, "'")) {
			commandParts[i] = fmt.Sprintf("'%s'", arg)
		} else {
			commandParts[i] = arg
		}
	}

	head, err := buildCommitMessageHead()
	if err != nil {
		fmt.Printf("Failed to build commit prefix: %+v\n", err)
		panic(err)
	}

	parts := []string{head}
	if len(envs) > 0 {
		parts = append(parts, strings.Join(envs, " "))
	}
	parts = append(parts, strings.Join(commandParts, " "))
	return strings.Join(parts, " ")
}

func buildCommitMessageHead() (string, error) {
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

	var location string
	if relPath == "." {
		location = rootDirName
	} else if strings.HasPrefix(relPath, "./") {
		location = rootDirName + relPath[1:]
	} else if strings.HasPrefix(relPath, "/") {
		location = relPath
	} else {
		location = rootDirName + "/" + relPath
	}

	return fmt.Sprintf(commitPrefix, location), nil
}

func gitRootDir() (string, error) {
	out, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}
