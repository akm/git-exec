# git-exec

## install

```
go install github.com/akm/git-exec@latest
```

You might have to run `asdf reshim` .

## Usage

```
git exec (your command with argument)
```

git-exec executes given command and `git add .` and `git commit` automatically.

### Help

```
Usage: git-exec [options ...] [key=value ...] <command> [args ...]

Options:
  -C, --directory                      Specify the directory where the command is executed. (default: "")
  -e, --emoji                          Specify the emoji used in commit message. (default: "🤖")
  -p, --prompt                         Specify the prompt used in commit message. (default: "$")
  -t, --template                       Specify the template to build commit message. (default: "{{.Emoji}} [{{.Location}}] {{.Prompt}} {{.Command}}")
      --skip-guard                     Skip the guard check for uncommitted changes and untracked files before executing command. (default: false)
      --skip-guard-uncommitted-changes Skip the guard check for uncommitted changes before executing command. (default: false)
      --skip-guard-untracked-files     Skip the guard check for untracked files before executing command. (default: false)
  -D, --debug-log                      Output debug log.
  -i, --interactive                    Interactive mode for command which requires input. tmux is required to use.
  -h, --help                           Show this message.
  -v, --version                        Show version.

Environment variable mapping:
--emoji                          GIT_EXEC_EMOJI
--prompt                         GIT_EXEC_PROMPT
--template                       GIT_EXEC_TEMPLATE
--skip-guard                     GIT_EXEC_SKIP_GUARD
--skip-guard-uncommitted-changes GIT_EXEC_SKIP_GUARD_UNCOMMITTED_CHANGES
--skip-guard-untracked-files     GIT_EXEC_SKIP_GUARD_UNTRACKED_FILES
--debug-log                      GIT_EXEC_DEBUG_LOG
--interactive                    GIT_EXEC_INTERACTIVE

Examples:
* Specify environment variables.
	git exec FOO=fooooo make args1 args2

* Use shell to work with redirect operator.
	git exec /bin/bash -c 'echo "foo" >> README.md'

* Use interactive mode for command which requires input such as "npx sv create" for SvelteKit.
	git exec -i npx sv create my-app


```
