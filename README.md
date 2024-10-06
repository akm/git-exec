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

### Option

You can change message prefix by environment variable `GIT_EXEC_COMMIT_PREFIX` .

### Tips

1. Use shell for redirection

```
git exec /bin/bash -c 'echo "foo" >> README.md'
```

2. Pass environment variables before the executable command

```
git exec FOO=fooooo command1 args1 args2
```
