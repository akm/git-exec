package main

import (
	"bytes"
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
	return "🤖 %s $"
}()

func main() {
	// 1. このプログラムに渡された引数をコマンドとして実行する。
	//    その際には、コマンドの照準出力と標準エラー出力をバッファに格納する。
	// 2. "git add ." を実行し、コマンドによって作成・変更されたカレントディレクトリ以下のファイルを staging area に追加する。
	// 3. "git commit" を以下のオプションと標準力を指定して実行する。
	//    オプション : --file -
	//    標準入力: "🤖 (実行したコマンド)\n\n(バッファ)"

	if len(os.Args) < 2 {
		fmt.Println("Usage: git-exec <command>")
		return
	}

	envs, commandArgs := splitArgsToEnvsAndCommand(os.Args[1:])

	// 1. このプログラムに渡された引数をコマンドとして実行する。
	cmd := exec.Command(commandArgs[0], commandArgs[1:]...)
	cmd.Env = append(os.Environ(), envs...)
	var outBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &outBuf

	if err := cmd.Run(); err != nil {
		fmt.Printf("Command execution failed: %+v\n", err)
		return
	}

	// 2. "git add ." を実行し、コマンドによって作成・変更されたカレントディレクトリ以下のファイルを staging area に追加する。
	if err := exec.Command("git", "add", ".").Run(); err != nil {
		fmt.Printf("git add failed: %+v\n", err)
		return
	}

	// 3. "git commit" を以下のオプションと標準力を指定して実行する。
	commitMessage := buildCommitMessage(envs, commandArgs, &outBuf)

	// See https://tracpath.com/docs/git-commit/
	commitCmd := exec.Command("git", "commit", "--file", "-")
	commitCmd.Stdin = bytes.NewBufferString(commitMessage)

	if err := commitCmd.Run(); err != nil {
		fmt.Printf("git commit failed: %+v\n", err)
		return
	}
}

func splitArgsToEnvsAndCommand(args []string) ([]string, []string) {
	var envs []string
	var command []string
	equalNotFound := false
	for _, arg := range args {
		if !equalNotFound && strings.Contains(arg, "=") {
			envs = append(envs, arg)
		} else {
			equalNotFound = true
			command = append(command, arg)
		}
	}
	return envs, command
}

func buildCommitMessage(envs []string, commandArgs []string, outBuf *bytes.Buffer) string {
	firstLine := buildCommitMessageFirstLine(envs, commandArgs)
	return fmt.Sprintf("%s\n\n%s\n",
		firstLine,
		outBuf.String(),
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
	if strings.HasPrefix(relPath, "./") {
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
