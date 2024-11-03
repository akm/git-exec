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
	}

	if err := process(options, commandArgs); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func process(options *Options, commandArgs []string) error {
	if options.Directory != "" {
		if err := os.Chdir(options.Directory); err != nil {
			return fmt.Errorf("Failed to change directory: %s", err.Error())
		}
	}

	if err := guard(options); err != nil {
		if isGuardError(err) {
			return err
		} else {
			return fmt.Errorf("Guard failed: %+v", err)
		}
	}

	command := newCommand(commandArgs)

	if err := command.Run(); err != nil {
		return fmt.Errorf("Command execution failed: %+v\n%s", err, command.Output)
	}

	if err := add(); err != nil {
		return err
	}

	// 3. "git commit" を以下のオプションと標準力を指定して実行する。
	commitMessage, err := newCommitMessage(command, options).Build()
	if err != nil {
		return fmt.Errorf("Failed to build commit message: %+v", err)
	}

	// See https://tracpath.com/docs/git-commit/
	commitCmd := exec.Command("git", "commit", "--file", "-")
	commitCmd.Stdin = bytes.NewBufferString(commitMessage)

	if err := commitCmd.Run(); err != nil {
		return fmt.Errorf("git commit failed: %+v", err)
	}

	return nil
}
