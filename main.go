package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

func main() {
	// 1. このプログラムに渡された引数をコマンドとして実行する。
	//    その際には、コマンドの照準出力と標準エラー出力をバッファに格納する。
	// 2. "git add ." を実行し、コマンドによって作成・変更されたカレントディレクトリ以下のファイルを staging area に追加する。
	// 3. "git commit" を以下のオプションと標準力を指定して実行する。
	//    オプション : --file -
	//    標準入力: "🤖 (実行したコマンド)\n\n(バッファ)"

	if len(os.Args) < 2 {
		help()
		os.Exit(1)
	}

	options, commandArgs := splitToOptionsAndCommandArgs(os.Args[1:])
	for _, option := range options {
		switch option {
		case "-h", "--help":
			help()
			os.Exit(0)
		default:
			fmt.Fprintf(os.Stderr, "Unknown option: %s\n", option)
		}
	}

	if err := guard(); err != nil {
		if isGuardError(err) {
			fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		} else {
			fmt.Printf("Guard failed: %+v\n", err)
		}
		os.Exit(1)
	}

	command := newCommand(commandArgs)

	if err := command.Run(); err != nil {
		fmt.Printf("Command execution failed: %+v\n%s", err, command.Output.String())
		return
	}

	if !hasDiff() {
		fmt.Printf("No changes to commit\n%s", command.Output.String())
		return
	}

	// 2. "git add ." を実行し、コマンドによって作成・変更されたカレントディレクトリ以下のファイルを staging area に追加する。
	if err := exec.Command("git", "add", ".").Run(); err != nil {
		fmt.Printf("git add failed: %+v\n", err)
		return
	}

	// 3. "git commit" を以下のオプションと標準力を指定して実行する。
	commitMessage := buildCommitMessage(command)

	// See https://tracpath.com/docs/git-commit/
	commitCmd := exec.Command("git", "commit", "--file", "-")
	commitCmd.Stdin = bytes.NewBufferString(commitMessage)

	if err := commitCmd.Run(); err != nil {
		fmt.Printf("git commit failed: %+v\n", err)
		return
	}
}
