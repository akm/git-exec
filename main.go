package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

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
		fmt.Printf("Command execution failed: %+v\n%s", err, outBuf.String())
		return
	}

	if err := exec.Command("git", "diff", "--exit-code").Run(); err == nil {
		fmt.Printf("No changes to commit\n%s", outBuf.String())
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
