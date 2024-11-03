package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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

	options, commandArgs, err := parseOptions(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse arguments: %s\n", err.Error())
	}
	if options.Help {
		help()
		os.Exit(0)
	} else if options.Version {
		if len(commandArgs) == 0 {
			showVersion()
			os.Exit(0)
		} else {
			showVersionWithExecName(filepath.Base(os.Args[0]))
		}
	} else if options.Directory != "" {
		if err := os.Chdir(options.Directory); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to change directory: %s\n", err.Error())
			os.Exit(1)
		}
	}

	if err := guard(options); err != nil {
		if isGuardError(err) {
			fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		} else {
			fmt.Printf("Guard failed: %+v\n", err)
		}
		os.Exit(1)
	}

	command := newCommand(commandArgs)

	if err := command.Run(); err != nil {
		fmt.Printf("Command execution failed: %+v\n%s", err, command.Output)
		return
	}

	uncommittedChanges, err := hasUncommittedChanges()
	if err != nil {
		fmt.Printf("git diff failed: %+v\n", err)
		return
	}
	untrackedFiles, err := hasUntrackedFiles()
	if err != nil {
		fmt.Printf("git ls-files failed: %+v\n", err)
		return
	}

	if !uncommittedChanges && !untrackedFiles {
		fmt.Printf("No changes to commit and No untracked files\n")
		return
	}

	// 2. "git add ." を実行し、コマンドによって作成・変更されたカレントディレクトリ以下のファイルを staging area に追加する。
	if err := exec.Command("git", "add", ".").Run(); err != nil {
		fmt.Printf("git add failed: %+v\n", err)
		return
	}

	// 3. "git commit" を以下のオプションと標準力を指定して実行する。
	commitMessage, err := newCommitMessage(command, options).Build()
	if err != nil {
		fmt.Printf("Failed to build commit message: %+v\n", err)
		return
	}

	// See https://tracpath.com/docs/git-commit/
	commitCmd := exec.Command("git", "commit", "--file", "-")
	commitCmd.Stdin = bytes.NewBufferString(commitMessage)

	if err := commitCmd.Run(); err != nil {
		fmt.Printf("git commit failed: %+v\n", err)
		return
	}
}
